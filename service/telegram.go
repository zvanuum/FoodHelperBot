package service

import (
	"fmt"
	"log"
	"net/http"

	"github.com/zachvanuum/FoodHelperBot/model"
	"github.com/zachvanuum/FoodHelperBot/util"
)

type TelegramService interface {
	GetMe() (model.BotInfo, error)
	RegisterWebhook(url string) error
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
		return botInfo, fmt.Errorf("failed to get bot, %s", err.Error())
	}

	defer res.Body.Close()

	var resBody model.BotInfoResponseWrapper
	err = util.UnmarshalResponse(res.Body, &resBody)
	if err != nil {
		return botInfo, fmt.Errorf("couldn't marshall getMe Telegram response to struct: %s", err.Error())
	}

	log.Printf("getMe response succeeded: %t.\n", resBody.OK)
	if resBody.OK {
		log.Printf("Bot info -  ID: %d, Name: %s, Username: %s\n", resBody.Result.ID, resBody.Result.Name, resBody.Result.Username)
	}

	return resBody.Result, nil
}

func (svc telegramService) RegisterWebhook(url string) error {
	telegramURL := fmt.Sprintf("https://api.telegram.org/bot%s/setWebhook?url=%s", svc.Token, url)

	res, err := http.Post(telegramURL, "", nil)
	if err != nil {
		return fmt.Errorf("failed to do POST request to %s: %s", telegramURL, err.Error())
	}

	if res.StatusCode >= 300 {
		return fmt.Errorf("bad response status when setting webhook: %s %d", res.Status, res.StatusCode)
	}

	defer res.Body.Close()

	var resBody model.SetWebhookResponse
	if err := util.UnmarshalResponse(res.Body, &resBody); err != nil {
		return fmt.Errorf("couldn't marshall setWebhook Telegram response to struct: %s", err.Error())
	}

	log.Printf("setWebhook response succeeded: %t.\n", resBody.OK)
	if resBody.OK {
		log.Printf("Result: %t, Description: %s\n", resBody.Result, resBody.Description)
	}

	return nil
}
