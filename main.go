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

type Flags struct {
	TelegramToken string
	Port          string
	Cert          string
	Key           string
}

func main() {
	flags := getFlags()
	services := createServices(flags.TelegramToken)
	foodBot := setupBot(services.TelegramService)
	routes := createRoutes(foodBot)
	server := createServer(flags.Port, routes)

	botServer := createBotServer(services, server, foodBot)

	log.Printf("[main] Starting server on %s\n", flags.Port)
	if flags.Port == "443" {
		if flags.Cert == "" || flags.Key == "" {
			log.Fatal("No SSL certificates were provided, exiting")
		}
		log.Fatal(botServer.Server.ListenAndServeTLS(flags.Cert, flags.Key))
	} else {
		log.Fatal(botServer.Server.ListenAndServe())
	}
}

func getFlags() Flags {
	var telegramToken string
	flag.StringVar(&telegramToken, "token", "", "The token used to authenticate with Telegram")
	var port string
	flag.StringVar(&port, "port", "8080", "The port to run on")
	var cert, key string
	flag.StringVar(&cert, "cert", "", "SSL Certificate")
	flag.StringVar(&key, "key", "", "Private key")
	flag.Parse()

	return Flags{
		TelegramToken: telegramToken,
		Port:          port,
		Cert:          cert,
		Key:           key,
	}
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

func setupBot(telegramService service.TelegramService) bot.FoodHelperBot {
	botInfo, err := telegramService.GetMe()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	webhookURL := "https://zvanuum.com/message"
	if err := telegramService.RegisterWebhook(webhookURL); err != nil {
		log.Fatalf("Failed to register webhook for bot using url %s", webhookURL)
	}

	return bot.NewTelegramBot(botInfo)
}

func createRoutes(foodBot bot.FoodHelperBot) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", controller.GreetingHandler(foodBot)).Methods("GET")
	r.HandleFunc("/message", controller.ReceiveMessageHandler()).Methods("POST")

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
