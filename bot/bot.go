package bot

import (
	"fmt"

	"github.com/zachvanuum/FoodHelperBot/model"
)

type FoodHelperBot interface {
	Greeting() string
}

type foodHelperBot struct {
	ID       int64
	Name     string
	Username string
}

func NewTelegramBot(info model.BotInfo) FoodHelperBot {
	return &foodHelperBot{
		ID:       info.ID,
		Name:     info.Name,
		Username: info.Username,
	}
}

func (bot foodHelperBot) Greeting() string {
	return fmt.Sprintf("Hello, my name is %s. You can contact me by messaging @%s.\n", bot.Name, bot.Username)
}
