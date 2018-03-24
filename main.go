package main

import (
	"flag"
	"fmt"
	"os"

	. "github.com/zachvanuum/FoodHelperBot/service/telegram"
)

type Services struct {
	TelegramService TelegramService
}

func main() {
	telegramToken := getToken()
	fmt.Println(telegramToken)

	services := createServices(telegramToken)
	botInfo, err := services.TelegramService.GetMe()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("%v\n", botInfo)
}

func getToken() string {
	var token string
	flag.StringVar(&token, "token", "", "The bot's token")
	flag.Parse()
	return token
}

func createServices(telegramToken string) Services {
	return Services{
		TelegramService: NewTelegramService(telegramToken),
	}
}
