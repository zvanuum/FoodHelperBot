package controller

import (
	"log"
	"net/http"

	"github.com/zachvanuum/FoodHelperBot/bot"
)

func ReceiveMessageHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[ReceiveMessageHandler] Got request: %s", r.Body)
		w.Write([]byte("test"))
	}
}

func GreetingHandler(foodBot bot.FoodHelperBot) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[GreetingHandler] Got request")
		w.Write([]byte(foodBot.Greeting()))
	}
}
