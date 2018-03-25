package controller

import (
	"log"
	"net/http"

	"github.com/zachvanuum/FoodHelperBot/bot"
)

func GreetingHandler(foodBot bot.FoodHelperBot) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[GreetingHandler] got request")
		w.Write([]byte(foodBot.Greeting()))
	}
}
