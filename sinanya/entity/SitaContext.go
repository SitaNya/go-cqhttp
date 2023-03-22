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
	Messages     []IMessage `json:"messages"`
	MessageTypes string     `json:"messageTypes"`
}

type IMessage interface {
	GetType() string
}

type MessageText struct {
	Text string `json:"text,omitempty"`
	Type string `json:"type,omitempty"`
}

func (t MessageText) GetType() string {
	return t.Type
}

type MessageAt struct {
	Id   int    `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
}

func (t MessageAt) GetType() string {
	return t.Type
}

type MessageImage struct {
	Type string `json:"type,omitempty"`
	Url  string `json:"url,omitempty"`
}

func (t MessageImage) GetType() string {
	return t.Type
}
