package service

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/zachvanuum/FoodHelperBot/model"
)

const (
	// Recognized user commands
	HelpCommand   = "/help"
	RandomCommand = "/random"
	SearchCommand = "/search"
	StartCommand  = "/start"

	// Response messages
	BadCommandResponse   = "Valid queries start with \"/\", for example \"/search <term>\" will search for businesses near you."
	DefaultResponse      = "Sorry, but I don't know how to answer that query."
	FailedResponse       = "Sorry, I was unable to perform that search."
	greetingStringFormat = `Hello, my name is %s. You can contact me by messaging @%s. 
	Accepted requests are:
		"/search <cuisine/business> in <location>",
		"/search <cuisine/business> nearby/near me", and
		"/random"
	To see these again send "/start" or "/help".`
	LocationResponse = "Please provide your location so that I can search for businesses near you."
	ThanksResponse   = "Thank you!"

	LocationKeyboardText = "Provide Location"

	searchLocationDelimiter     = " in "
	userLocationDelimiterNearMe = " near me"
	userLocationDelimiterNearby = " nearby"
)

type BotService interface {
	CreateResponseMessage(message model.ReceivedMessage) *model.Message
	Greeting() string
}

type botService struct {
	ID          int64
	Name        string
	Username    string
	YelpService YelpService
	UsersCache  map[int64]*model.UserLocationInfo
}

func NewTelegramBot(info model.BotInfo, yelp YelpService) BotService {
	return &botService{
		ID:          info.ID,
		Name:        info.Name,
		Username:    info.Username,
		YelpService: yelp,
		UsersCache:  make(map[int64]*model.UserLocationInfo),
	}
}

func (svc botService) CreateResponseMessage(message model.ReceivedMessage) *model.Message {
	command, remaining := splitUserMessageToQuery(message.Message.Text)

	log.Printf("[CreateResponseMessage] User query: %s, remaining message: \"%s\"", command, remaining)

	response := model.NewMessage(message.Message.Chat.ID, "")

	switch command {
	case StartCommand, HelpCommand:
		response.Text = svc.Greeting()
	case SearchCommand:
		term := getUserSearchTerm(remaining)
		svc.updateUserLastSearchTerm(message.Message.Chat.ID, term)

		if isUserLocationSearchQuery(remaining) {
			log.Printf("[createResponseMessage] User search term: %s", term)
			addLocationKeyboardMarkup(response)
			response.Text = LocationResponse
			break
		}

		location := getUserSepcifiedSearchLocation(remaining)
		log.Printf("[createResponseMessage] User search term: %s, user search location: %s", term, location)

		searchResults, err := svc.YelpService.SearchByLocation(term, location)
		if err != nil {
			log.Printf("[createResponseMessage] %s", err.Error())

			response.Text = FailedResponse
		} else {
			svc.createSearchResponse(response, searchResults)
		}
	case RandomCommand:
		addLocationKeyboardMarkup(response)
		response.Text = LocationResponse
		svc.updateUserLastSearchTerm(message.Message.Chat.ID, getRandomCuisine())
	default:
		if isProvidingLocation(message) &&
			svc.UsersCache[message.Message.Chat.ID].LastCommand == SearchCommand ||
			svc.UsersCache[message.Message.Chat.ID].LastCommand == RandomCommand {
			log.Printf("[createResponseMessage] Got user's location - Chat ID: %d, Message ID: %d, Location: %f, %f",
				message.Message.Chat.ID,
				message.Message.MessageID,
				message.Message.Location.Latitude,
				message.Message.Location.Longitude,
			)

			searchResults, err := svc.YelpService.SearchByCoordinates(svc.UsersCache[message.Message.Chat.ID].LastSearchTerm, message.Message.Location.Latitude, message.Message.Location.Longitude)
			if err != nil {
				log.Printf("[createResponseMessage] %s", err.Error())

				response.Text = FailedResponse
			} else {
				svc.createSearchResponse(response, searchResults)
			}

			break
		}

		if !strings.Contains(command, "/") {
			response.Text = BadCommandResponse
		} else {
			response.Text = DefaultResponse
		}
	}

	svc.updateUserLastCommand(message.Message.Chat.ID, command)

	response.ReplyToMessageID = message.Message.MessageID
	return response
}

func (svc botService) Greeting() string {
	return fmt.Sprintf(greetingStringFormat, svc.Name, svc.Username)
}

func (svc botService) updateUserLocation(chatID int64, latitude float64, longitude float64) {
	if _, ok := svc.UsersCache[chatID]; ok {
		svc.UsersCache[chatID].Location.Latitude = latitude
		svc.UsersCache[chatID].Location.Longitude = longitude
	} else {
		svc.UsersCache[chatID] = &model.UserLocationInfo{
			Location: model.Coordinates{
				Latitude:  latitude,
				Longitude: longitude,
			},
		}
	}
}

func (svc botService) updateUserLastCommand(chatID int64, command string) {
	if _, ok := svc.UsersCache[chatID]; ok {
		svc.UsersCache[chatID].LastCommand = command
	} else {
		svc.UsersCache[chatID] = &model.UserLocationInfo{
			LastCommand: command,
		}
	}
}

func (svc botService) updateUserLastSearchTerm(chatID int64, term string) {
	if _, ok := svc.UsersCache[chatID]; ok {
		svc.UsersCache[chatID].LastSearchTerm = term
	} else {
		svc.UsersCache[chatID] = &model.UserLocationInfo{
			LastSearchTerm: term,
		}
	}
}

func splitUserMessageToQuery(text string) (string, string) {
	splitText := strings.Split(text, " ")
	command := splitText[0]
	remaining := strings.Join(splitText[1:len(splitText)], " ")

	return command, remaining
}

func isUserLocationSearchQuery(query string) bool {
	userLocationSearchDelimitters := []string{userLocationDelimiterNearMe, userLocationDelimiterNearby}

	for _, delimitter := range userLocationSearchDelimitters {
		if strings.Contains(query, delimitter) {
			return true
		}
	}

	return false
}

func getUserSearchTerm(text string) string {
	nearbyIndex := strings.LastIndex(text, userLocationDelimiterNearby)
	nearMeIndex := strings.LastIndex(text, userLocationDelimiterNearMe)
	inIndex := strings.LastIndex(text, searchLocationDelimiter)
	var term string

	if nearbyIndex > -1 {
		term = text[0:nearbyIndex]
	} else if nearMeIndex > -1 {
		term = text[0:nearMeIndex]
	} else if inIndex > -1 {
		term = text[0:inIndex]
	}

	return term
}

func getUserSepcifiedSearchLocation(text string) string {
	inIndex := strings.LastIndex(text, searchLocationDelimiter)
	var location string

	if inIndex > -1 {
		location = text[inIndex+len(searchLocationDelimiter) : len(text)]
	}

	return location
}

func (svc botService) createSearchResponse(response *model.Message, result model.SearchResponse) {
	response.ParseMode = "Markdown"
	response.ReplyMarkup.Keyboard = [][]model.KeyboardButton{}

	showCount := 10
	if result.Total < 10 {
		showCount = result.Total
	}

	responseString := fmt.Sprintf(
		"Got %d results searching for %s, here are the top %d!\n\n",
		result.Total,
		svc.UsersCache[response.ChatID].LastSearchTerm,
		showCount,
	)

	for i, business := range result.Businesses[0:showCount] {
		businessStr := fmt.Sprintf(
			"[%d: %s](%s)\n%s, %s\n%s\n\n",
			i+1,
			business.Name,
			business.URL,
			getStars(business.Rating, business.ReviewCount),
			business.Price,
			business.Location.Address1,
		)

		responseString += businessStr
	}

	response.Text = responseString
}

func getStars(rating float64, reviewCount int) string {
	var stars string

	for i := 0; i < int(math.Round(rating)); i++ {
		stars += "⭐️"
	}

	return fmt.Sprintf("%s (%.2f, %d reviews)", stars, rating, reviewCount)
}

func isProvidingLocation(message model.ReceivedMessage) bool {
	return message.Message.Text == "" &&
		message.Message.Location.Latitude != 0 &&
		message.Message.Location.Longitude != 0
}

func addLocationKeyboardMarkup(message *model.Message) {
	message.ReplyMarkup = model.ReplyMarkup{
		Keyboard: [][]model.KeyboardButton{
			[]model.KeyboardButton{
				model.KeyboardButton{
					Text:            LocationKeyboardText,
					RequestLocation: true,
				},
			},
		},
		ResizeKeyboard: true,
	}
}

func getRandomCuisine() string {
	possibilities := []string{
		"mexican", "indian", "breakfast", "cafe", "seafood",
		"chinese", "japanese", "thai", "vietnamese", "ethiopian",
		"american", "burgers", "gastropub", "sandwiches", "filipino",
		"ramen", "pho", "french", "greek", "german",
		"moroccan", "soul food", "cajun", "carribean",
		"turkish", "spanish", "italian", "korean", "lebanese",
		"hawaiian", "jamaican", "brazillian", "british", "mediterranean",
	}

	rand.Seed(time.Now().Unix())

	return possibilities[rand.Intn(len(possibilities)-1)]
}
