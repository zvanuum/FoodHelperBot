package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/zachvanuum/FoodHelperBot/model"
	"github.com/zachvanuum/FoodHelperBot/util"
)

type TelegramService interface {
	GetMe() (model.BotInfo, error)
	RegisterWebhook(url string) error
	RespondToMessage(model.ReceivedMessage) error
}

type telegramService struct {
	Token       string
	YelpService YelpService
	BotService  BotService
}

func NewTelegramService(token string, yelpService YelpService) TelegramService {
	service := telegramService{
		Token:       token,
		YelpService: yelpService,
	}

	botService := service.setupBotService()
	service.BotService = botService

	return service
}

func (svc telegramService) setupBotService() BotService {
	botInfo, err := svc.GetMe()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	webhookURL := "https://zvanuum.com/message"
	if err := svc.RegisterWebhook(webhookURL); err != nil {
		log.Fatalf("[setupBot] Failed to register webhook for bot using url %s", webhookURL)
	}

	return NewTelegramBot(botInfo, svc.YelpService)
}

func (svc telegramService) GetMe() (model.BotInfo, error) {
	var botInfo model.BotInfo
	var err error

	res, err := http.Get(fmt.Sprintf("https://api.telegram.org/bot%s/getMe", svc.Token))
	if err != nil {
		return botInfo, fmt.Errorf("failed to get bot, %s", err.Error())
	}

	defer res.Body.Close()

	log.Printf("[GetMe] Response status: %s", res.Status)
	if res.StatusCode >= 300 {
		return botInfo, fmt.Errorf("bad response status when getting bot info: %s", res.Status)
	}

	var botInfoWrapper model.BotInfoResponseWrapper
	err = util.UnmarshalBody(res.Body, &botInfoWrapper)
	if err != nil {
		return botInfo, fmt.Errorf("failed to marshall getMe response to struct: %s", err.Error())
	}

	log.Printf("[GetMe] /getMe response succeeded: %t.\n", botInfoWrapper.OK)
	if botInfoWrapper.OK {
		log.Printf("[GetMe] Bot info -  ID: %d, Name: %s, Username: %s\n", botInfoWrapper.Result.ID, botInfoWrapper.Result.Name, botInfoWrapper.Result.Username)
	}

	return botInfoWrapper.Result, nil
}

func (svc telegramService) RegisterWebhook(url string) error {
	telegramURL := fmt.Sprintf("https://api.telegram.org/bot%s/setWebhook?url=%s", svc.Token, url)

	res, err := http.Post(telegramURL, "", nil)
	if err != nil {
		return fmt.Errorf("failed to do POST request to %s: %s", telegramURL, err.Error())
	}

	log.Printf("[RegisterWebhook] Response status: %s", res.Status)
	if res.StatusCode >= 300 {
		return fmt.Errorf("bad response status when setting webhook: %s", res.Status)
	}

	defer res.Body.Close()

	var setWebhookResponse model.SetWebhookResponse
	if err := util.UnmarshalBody(res.Body, &setWebhookResponse); err != nil {
		return fmt.Errorf("failed to marshall setWebhook response to struct: %s", err.Error())
	}

	log.Printf("[RegisterWebhook] /setWebhook response succeeded: %t.\n", setWebhookResponse.OK)
	if setWebhookResponse.OK {
		log.Printf("[RegisterWebhook] Result: %t, Description: %s\n", setWebhookResponse.Result, setWebhookResponse.Description)
	}

	return nil
}

func (svc telegramService) RespondToMessage(message model.ReceivedMessage) error {
	responseMessage := svc.BotService.CreateResponseMessage(message)

	log.Printf(
		"[RespondToMessage] Response message - chat ID: %d, message ID: %d, user ID: %d, text: \"%s\"",
		message.Message.Chat.ID,
		message.Message.MessageID,
		message.Message.From.ID,
		responseMessage.Text,
	)

	var sendMessageURL string
	var req *http.Request
	var err error
	if len(responseMessage.ReplyMarkup.Keyboard) > 0 {
		postBody, err := json.Marshal(responseMessage)
		if err != nil {
			return fmt.Errorf("failed to marshal struct %v to json: %s", responseMessage, err.Error())
		}

		sendMessageURL = fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", svc.Token)
		req, err = http.NewRequest("POST", sendMessageURL, bytes.NewBuffer(postBody))
		if err != nil {
			return fmt.Errorf("failed to make POST request to %s: %s", sendMessageURL, err.Error())
		}

		req.Header.Set("Content-Type", "application/json")
	} else {
		sendMessageURL = formatSendMessageURL(svc.Token, responseMessage.ChatID, responseMessage.Text)
		req, err = http.NewRequest("GET", sendMessageURL, nil)
		if err != nil {
			return fmt.Errorf("failed to make GET request to %s: %s", sendMessageURL, err.Error())
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to do request to %s: %s", sendMessageURL, err.Error())
	}

	defer res.Body.Close()

	log.Printf("[RespondToMessage] Response status: %s", res.Status)

	var sendMessageResponse model.SendMessageResponse
	if err := util.UnmarshalBody(res.Body, &sendMessageResponse); err != nil {
		return fmt.Errorf("failed to marshall sendMessage response to struct: %s", err.Error())
	}

	log.Printf("[RespondToMessage] /sendMessage response succeeded: %t.\n", sendMessageResponse.OK)
	if !sendMessageResponse.OK {
		return fmt.Errorf("failed to send message")
	}

	return nil
}

func formatSendMessageURL(token string, chatId int64, responseText string) string {
	return fmt.Sprintf(
		"https://api.telegram.org/bot%s/sendMessage?chat_id=%d&text=%s",
		token,
		chatId,
		url.QueryEscape(responseText),
	)
}
