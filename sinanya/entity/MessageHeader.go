package entity

import (
	"fmt"
	"github.com/GehirnInc/crypt/md5_crypt"
	"strconv"
)

type MessageHeader struct {
	BotId        int    `json:"botId"`
	PlatformType string `json:"platformType"`
	Type         string `json:"type"`
	Token        string `json:"token"`
}

func GetToken(botId int, platformType string) string {
	passwd := []byte(fmt.Sprintf("%s%s", strconv.Itoa(botId), platformType))
	salt := []byte("$1$445158612931")
	ret, _ := md5_crypt.New().Generate(passwd, salt)
	return ret
}
