package entity

import (
	"encoding/json"
)

type SitaContext struct {
	BotId        int          `json:"botId"`
	UserId       int          `json:"userId"`
	GroupId      int          `json:"groupId"`
	Type         string       `json:"type"`
	MessagesList MessagesList `json:"messagesList"`
	ActionTypes  ActionTypes  `json:"actionTypes"`
	Platform     string       `json:"platform"`
}

type ActionTypes struct {
	Id      int64  `json:"id"`
	Context string `json:"context"`
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

func parseMessagesList(messagesList map[string]interface{}) MessagesList {
	messages := messagesList["messages"].([]interface{})
	if messages == nil {
		messages = make([]interface{}, 0)
	}
	messageTypes := messagesList["messageTypes"].(string)

	var resultMessages []IMessage
	for _, message := range messages {
		messageObj := message.(map[string]interface{})
		messageType := messageObj["type"].(string)
		switch messageType {
		case "Text":
			textMessage := MessageText{}
			jsonBytes, _ := json.Marshal(messageObj)
			_ = json.Unmarshal(jsonBytes, &textMessage)
			resultMessages = append(resultMessages, textMessage)
		case "At":
			atMessage := MessageAt{}
			jsonBytes, _ := json.Marshal(messageObj)
			_ = json.Unmarshal(jsonBytes, &atMessage)
			resultMessages = append(resultMessages, atMessage)
		case "Image":
			imageMessage := MessageImage{}
			jsonBytes, _ := json.Marshal(messageObj)
			_ = json.Unmarshal(jsonBytes, &imageMessage)
			resultMessages = append(resultMessages, imageMessage)
		}
	}

	return MessagesList{Messages: resultMessages, MessageTypes: messageTypes}
}

func ParseSitaContext(jsonStr string) (SitaContext, error) {
	sitaContext := SitaContext{}
	_ = json.Unmarshal([]byte(jsonStr), &sitaContext)
	data := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		panic(err)
	}

	sitaContext.MessagesList = parseMessagesList(data["messagesList"].(map[string]interface{}))

	return sitaContext, nil
}
