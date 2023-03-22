package tools

import (
	"github.com/Mrs4s/go-cqhttp/sinanya/entity"
	"strconv"
	"strings"
)

const (
	image = "[CQ:image,file=h"
	at    = "[CQ:at,qq="
)

func MakeMessageList(input string) (result []entity.IMessage) {
	var tmpText string
	var tmpType string
	var tmpTypeBool bool
	for _, charItem := range input {
		if charItem == '[' {
			tmpTypeBool = true
			result = append(result, entity.MessageText{Type: "Text", Text: tmpText})
			tmpText = ""
		} else if charItem == ']' {
			tmpTypeBool = false
			if strings.HasPrefix(tmpType, image) {
				result = append(result, entity.MessageImage{Type: "Image", Url: strings.TrimSuffix(strings.TrimPrefix(tmpType, image), "]")})
			} else if strings.HasPrefix(tmpType, at) {
				id, _ := strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(tmpType, image), "]"))
				result = append(result, entity.MessageAt{Type: "At", Id: id})
			}
		} else if tmpTypeBool {
			tmpType = tmpType + string(charItem)
		} else {
			tmpText = tmpText + string(charItem)
		}
	}
	if tmpText != "" {
		result = append(result, entity.MessageText{Type: "Text", Text: tmpText})
		tmpText = ""
	}
	return
}
