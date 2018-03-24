package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/zachvanuum/FoodHelperBot/service/telegram"
)

type Services struct {
	TelegramService telegram.TelegramService
}

func main() {
	telegramToken := getToken()
	services := createServices(telegramToken)

	botInfo, err := services.TelegramService.GetMe()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	bot := NewTelegramBot(botInfo)

	bot.Greeting()
}

func getToken() string {
	var token string
	flag.StringVar(&token, "token", "", "The bot's token")
	flag.Parse()
	return token
}

func createServices(telegramToken string) Services {
	return Services{
		TelegramService: telegram.NewTelegramService(telegramToken),
	}
}
