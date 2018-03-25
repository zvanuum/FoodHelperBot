package service

import (
	"fmt"
	"net/http"

	"github.com/zachvanuum/FoodHelperBot/model"
	"github.com/zachvanuum/FoodHelperBot/util"
)

type TelegramService interface {
	GetMe() (model.BotInfo, error)
}

type telegramService struct {
	Token string
}

func NewTelegramService(token string) TelegramService {
	return telegramService{
		Token: token,
	}
}

func (svc telegramService) GetMe() (model.BotInfo, error) {
	var botInfo model.BotInfo
	var err error

	res, err := http.Get(fmt.Sprintf("https://api.telegram.org/bot%s/getMe", svc.Token))
	if err != nil {
		return botInfo, fmt.Errorf("failed to get bot, %s", err)
	}

	defer res.Body.Close()

	var resBody model.BotInfoResponseWrapper
	err = util.UnmarshalResponse(res, &resBody)
	if err != nil {
		return botInfo, fmt.Errorf("couldn't marshall Telegram response to struct: %s", err)
	}

	return resBody.Result, nil
}