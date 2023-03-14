package entity

type SitaContext struct {
	BotId        int          `json:"botId"`
	UserId       int          `json:"userId"`
	GroupId      int          `json:"groupId"`
	Type         string       `json:"type"`
	MessagesList MessagesList `json:"messagesList"`
	ActionTypes  []string     `json:"actionTypes"`
	Platform     string       `json:"platform"`
}

type MessagesList struct {
	Messages     []Message `json:"messages"`
	MessageTypes string    `json:"messageTypes"`
}

type Message struct {
	Text string `json:"text,omitempty"`
	Type string `json:"type"`
	Id   int    `json:"id,omitempty"`
}
