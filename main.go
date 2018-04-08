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
	YelpService     service.YelpService
}

type Flags struct {
	TelegramToken string
	YelpKey       string
	Port          string
	Cert          string
	Key           string
}

func main() {
	flags := getFlags()
	services := createServices(flags.TelegramToken, flags.YelpKey)
	routes := createRoutes(services)
	server := createServer(flags.Port, routes)

	log.Printf("[main] Starting server on %s\n", flags.Port)

	if flags.Cert != "" || flags.Key != "" {
		if flags.Cert == "" {
			log.Fatal("[main] Missing certificate, exiting")
		} else if flags.Key == "" {
			log.Fatal("[main] Missing key, exiting")
		}

		log.Fatal(server.ListenAndServeTLS(flags.Cert, flags.Key))
	} else {
		log.Fatal(server.ListenAndServe())
	}
}

func getFlags() Flags {
	var telegramToken string
	flag.StringVar(&telegramToken, "token", "", "The token used to authenticate with Telegram")
	var yelpKey string
	flag.StringVar(&yelpKey, "yelpKey", "", "The API key used to authenticate with Yelp")
	var port string
	flag.StringVar(&port, "port", "8080", "The port to run on")
	var cert, key string
	flag.StringVar(&cert, "cert", "", "SSL Certificate")
	flag.StringVar(&key, "key", "", "Private key")
	flag.Parse()

	return Flags{
		TelegramToken: telegramToken,
		YelpKey:       yelpKey,
		Port:          port,
		Cert:          cert,
		Key:           key,
	}
}

func createServices(telegramToken string, yelpKey string) *Services {
	yelpService := service.NewYelpService(yelpKey)
	telegramService := service.NewTelegramService(telegramToken, yelpService)

	return &Services{
		TelegramService: telegramService,
		YelpService:     yelpService,
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
