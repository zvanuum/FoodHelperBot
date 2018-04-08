package handler

import (
	"log"
	"net/http"

	"github.com/zachvanuum/FoodHelperBot/model"
	"github.com/zachvanuum/FoodHelperBot/service"
	"github.com/zachvanuum/FoodHelperBot/util"
)

func ReceiveMessageHandler(svc service.TelegramService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var message model.ReceivedMessage

		if err := util.UnmarshalBody(r.Body, &message); err != nil {
			log.Printf("[ReceiveMessageHandler] failed marshall sendMessage response to struct: %s", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		log.Printf(
			"[ReceiveMessageHandler] Got message - chat ID: %d, message ID: %d, user ID: %d, text: \"%s\"",
			message.Message.Chat.ID,
			message.Message.MessageID,
			message.Message.From.ID,
			message.Message.Text,
		)

		if err := svc.RespondToMessage(message); err != nil {
			log.Printf("[ReceiveMessageHandler] Error responding to message: %s", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}
