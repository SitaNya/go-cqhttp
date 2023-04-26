package entity

import (
	"fmt"
	"fyne.io/fyne/v2/app"
	"github.com/go-netty/go-netty"
	"runtime"
	"strings"
	"time"
)

var ChannelList = make(map[string]netty.Channel)

var ChannelMessage = make(chan SitaContext)

var ChannelAction = make(chan SitaContext)

var SERVER_REQUEST = []string{"GROUP_LEAVE_REQUEST", "GROUP_GET_NAME_REQUEST", "GROUP_GET_LIST_REQUEST", "GROUP_GET_LAST_SENDER_TIME_REQUEST", "GROUP_CHANGE_USER_NAME_REQUEST", "USER_GET_NAME_REQUEST"}

var LOGIN_CONFIG_FILE = "./LoginConfig.json"

type LoginConfig struct {
	UserName        int64
	Passwd          string
	ServerIp        map[string]string
	ServerRetryTime int
}

type RetryServerList struct {
	ServerIp  string
	RetryTime int64
}

var RETRY_SERVER_LIST = make(map[string]time.Duration)

var LOG_CHANNEL = make(chan string)

var OS_TYPE = runtime.GOOS

var WINDOW = app.New()

func GetMapKeys(m map[string]string) []string {
	// 数组默认长度为map长度,后面append时,不需要重新申请内存和拷贝,效率很高
	j := 0
	keys := make([]string, len(m))
	for k := range m {
		keys[j] = k
		j++
	}
	return keys
}

func MakeServerName(m []string) map[string]string {
	mapResult := make(map[string]string, len(m))
	for _, k := range m {
		if strings.Contains(k, "|") {
			mapResult[strings.Split(k, "|")[0]] = strings.Split(k, "|")[1]
		} else {
			mapResult[k] = k
		}
	}
	return mapResult
}

func MakeMapToServersIp(m map[string]string) []string {
	// 数组默认长度为map长度,后面append时,不需要重新申请内存和拷贝,效率很高
	j := 0
	keys := make([]string, len(m))
	for k, v := range m {
		keys[j] = fmt.Sprintf("%s|%s", k, v)
		j++
	}
	return keys
}
