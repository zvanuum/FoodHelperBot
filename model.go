package main

type TelegramResponse struct {
	OK     bool        `json:"ok"`
	Result interface{} `json:"result"`
}
