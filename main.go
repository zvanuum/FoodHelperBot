package main

import (
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
	telegramToken := os.Getenv("TELEGRAM_TOKEN")
	port := os.Getenv("PORT")
	services := createServices(telegramToken)

	botInfo, err := services.TelegramService.GetMe()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	foodBot := bot.NewTelegramBot(botInfo)
	routes := createRoutes(foodBot)
	server := createServer(port, routes)

	botServer := createBotServer(services, server, foodBot)

	log.Printf("[main] Starting server\n")
	log.Fatal(botServer.Server.ListenAndServe())
}

func createBotServer(services *Services, server *http.Server, foodBot bot.FoodHelperBot) *BotServer {
	return &BotServer{
		Services: services,
		Bot:      foodBot,
		Server:   server,
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

func createServer(port string, routes *mux.Router) *http.Server {
	return &http.Server{
		Handler:      routes,
		Addr:         ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}
