package model

type BotInfoResponseWrapper struct {
	OK     bool    `json:"ok"`
	Result BotInfo `json:"result"`
}

type BotInfo struct {
	ID       int64  `json:"id"`
	IsBot    bool   `json:"is_bot"`
	Name     string `json:"first_name"`
	Username string `json:"username"`
}

type SetWebhookResponse struct {
	OK          bool   `json:"ok"`
	Result      bool   `json:"result"`
	Description string `json:"description"`
}

type ReceivedMessage struct {
	UpdateID int64       `json:"update_id"`
	Message  MessageInfo `json:"message"`
}

type MessageInfo struct {
	Date        int64         `json:"date"`
	Chat        ChatInfo      `json:"chat"`
	MessageID   int64         `json:"message_id"`
	From        UserInfo      `json:"from"`
	Text        string        `json:"text"`
	ForwardFrom ForwarderInfo `json:"forward_from,omitempty"`
	ForwardDate int64         `json:"forward_date,omitempty"`
}

type ChatInfo struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Type      string `json:"type"`
}

type UserInfo struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	ID        int64  `json:"id"`
	Username  string `json:"username"`
}

type ForwarderInfo struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	ID        int64  `json:"id"`
}

type SendMessageResponse struct {
	OK     bool              `json:"ok"`
	Result SendMessageResult `json:"result"`
}

type SendMessageResult struct {
	MessageID int      `json:"message_id"`
	From      UserInfo `json:"from"`
	Chat      ChatInfo `json:"chat"`
	Date      int64    `json:"date"`
	Text      string   `json:"text"`
}
