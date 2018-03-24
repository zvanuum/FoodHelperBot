package telegram

type BotInfoResponseWrapper struct {
	OK     bool    `json:"ok"`
	Result BotInfo `json:"result"`
}

type BotInfo struct {
	ID        int64  `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
}
