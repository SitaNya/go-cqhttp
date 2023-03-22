package service

import (
	"encoding/json"
	"fmt"
	"github.com/Mrs4s/go-cqhttp/sinanya/entity"
	"github.com/go-netty/go-netty"
	"github.com/go-netty/go-netty/codec/format"
	"github.com/go-netty/go-netty/codec/frame"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)

type Client struct {
	BotId int
	Host  string
	Port  int
}

func (l Client) HandleActive(ctx netty.ActiveContext) {
	log.Infof("与%s服务器建立连接", l.Host)
	messageHeader := entity.MessageHeader{BotId: l.BotId, PlatformType: "QQ", Type: "MessageHeader", Token: entity.GetToken(l.BotId, "QQ")}
	str, _ := json.Marshal(messageHeader)
	ctx.Write(str)
	ctx.HandleActive()
	entity.ChannelList[l.Host] = ctx.Channel()
}

func (l Client) HandleRead(ctx netty.InboundContext, message netty.Message) {
	sitaContext, err := entity.ParseSitaContext(message.(string))
	if err != nil {
		panic(err)
	}
	entity.ChannelMessage <- sitaContext
	for _, messageEvery := range sitaContext.MessagesList.Messages {
		log.Infof("[%d]发来信息->\t%s", sitaContext.UserId, messageEvery)
	}
	ctx.HandleRead(message)
}

func (l Client) HandleInactive(ctx netty.InactiveContext, ex netty.Exception) {
	delete(entity.ChannelList, l.Host)
	ctx.Close(ex)
	loginConfig := entity.LoginConfig{}
	content, err := os.ReadFile(entity.LOGIN_CONFIG_FILE)
	err = json.NewDecoder(strings.NewReader(string(content))).Decode(&loginConfig)
	if err != nil {
		log.Fatal("配置文件不合法!", err)
	}
	if _, ok := entity.RETRY_SERVER_LIST[l.Host]; ok {
		entity.RETRY_SERVER_LIST[l.Host] = entity.RETRY_SERVER_LIST[l.Host] * 2
	} else {
		entity.RETRY_SERVER_LIST[l.Host] = time.Duration(loginConfig.ServerRetryTime*1000) * time.Millisecond
	}
	log.Warnf("从%s服务器断开连接，%d秒后重连", l.Host, entity.RETRY_SERVER_LIST[l.Host]/time.Millisecond/1000)
	time.Sleep(entity.RETRY_SERVER_LIST[l.Host])
	if _, ok := entity.ChannelList[l.Host]; ok {
		CreateClient(l.BotId, l.Host, l.Port)
	}
}

func CreateClient(botId int, host string, port int) {
	clientInitializer := func(channel netty.Channel) {
		channel.Pipeline().
			AddLast(frame.DelimiterCodec(1000000000, "\n", true)).
			AddLast(format.TextCodec()).
			AddLast(Client{Host: host, Port: port, BotId: botId})
	}
	log.Info(fmt.Sprintf("尝试以botId:%d为id注册服务器%s", botId, host))
	// new bootstrap
	var bootstrap = netty.NewBootstrap(netty.WithClientInitializer(clientInitializer))

	go func() { _, _ = bootstrap.Connect(fmt.Sprintf("%s:%d", host, port), nil) }()
}
