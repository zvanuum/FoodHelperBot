package controller

import (
	"net/http"

	"github.com/zachvanuum/FoodHelperBot/bot"
)

func GreetingHandler(foodBot bot.FoodHelperBot) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(foodBot.Greeting()))
	}
}
