package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/zachvanuum/FoodHelperBot/handler"
	"github.com/zachvanuum/FoodHelperBot/service"
)

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
	routes := createRoutes(services)
	server := createServer(flags.Port, routes)

	log.Printf("[main] Starting server on %s\n", flags.Port)

	if flags.Port == "443" {
		if flags.Cert == "" || flags.Key == "" {
			log.Fatal("No SSL certificates were provided, exiting")
		}

		log.Fatal(server.ListenAndServeTLS(flags.Cert, flags.Key))
	} else {
		log.Fatal(server.ListenAndServe())
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

func createServices(telegramToken string) *Services {
	return &Services{
		TelegramService: service.NewTelegramService(telegramToken),
	}
}

func createRoutes(services *Services) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/health", handler.HealthHandler()).Methods("GET")
	r.HandleFunc("/message", handler.ReceiveMessageHandler(services.TelegramService)).Methods("POST")

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
