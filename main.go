package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"

	"github.com/zachvanuum/FoodHelperBot/handler"
	"github.com/zachvanuum/FoodHelperBot/service"
)

type Services struct {
	TelegramService service.TelegramService
	YelpService     service.YelpService
}

type Flags struct {
	Config string
	Port   string
	Cert   string
	Key    string
}

func main() {
	flags := getFlags()

	viper.SetConfigFile(flags.Config)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("[main] Fatal error config file: %s \n", err.Error())
	}

	services := createServices(viper.GetString("telegram_key"), viper.GetString("yelp_key"))
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
	var config, port, cert, key string
	flag.StringVar(&config, "config", "./config.json", "The server configuration file")
	flag.StringVar(&port, "port", "8080", "The port to run on")
	flag.StringVar(&cert, "cert", "", "SSL Certificate")
	flag.StringVar(&key, "key", "", "Private key")
	flag.Parse()

	return Flags{
		Config: config,
		Port:   port,
		Cert:   cert,
		Key:    key,
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
