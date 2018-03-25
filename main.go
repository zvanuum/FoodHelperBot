package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"

	"github.com/zachvanuum/FoodHelperBot/bot"
	"github.com/zachvanuum/FoodHelperBot/controller"
	"github.com/zachvanuum/FoodHelperBot/service"
)

type BotServer struct {
	Bot      bot.FoodHelperBot
	Services *Services
	Server   *http.Server
}

type Services struct {
	TelegramService service.TelegramService
}

func main() {
	telegramToken := getToken()
	services := createServices(telegramToken)

	botInfo, err := services.TelegramService.GetMe()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	botServer := createBotServer(services, bot.NewTelegramBot(botInfo))

	fmt.Printf("Starting server\n")
	log.Fatal(botServer.Server.ListenAndServe())
}

func getToken() string {
	var token string
	flag.StringVar(&token, "token", "", "The bot's token")
	flag.Parse()
	return token
}

func createBotServer(services *Services, foodBot bot.FoodHelperBot) *BotServer {
	return &BotServer{
		Services: services,
		Bot:      foodBot,
		Server: &http.Server{
			Handler:      createRoutes(foodBot),
			Addr:         "localhost:8080",
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		},
	}
}

func createServices(telegramToken string) *Services {
	return &Services{
		TelegramService: service.NewTelegramService(telegramToken),
	}
}

func createRoutes(foodBot bot.FoodHelperBot) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", controller.GreetingHandler(foodBot)).Methods("GET")

	return r
}
