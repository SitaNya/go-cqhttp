package windows

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/Mrs4s/go-cqhttp/cmd/gocq"
	"github.com/Mrs4s/go-cqhttp/global/terminal"
	"github.com/Mrs4s/go-cqhttp/sinanya/entity"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
)

func getInput(label string) string {
	fmt.Print(label)
	reader := bufio.NewReader(os.Stdin)
	cmdString, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(os.Stderr, err)
	}
	return strings.TrimSuffix(cmdString, "\n")
}

func LoginByLinux() {
	_, err := os.Stat(entity.LOGIN_CONFIG_FILE)
	if err == nil {
		log.Infof("使用LoginConfig.json中的数据进行登录，如登录失败请确认LoginConfig.json中的账户密码是否正确")
		terminal.SetTitle()
		gocq.InitBase()
		gocq.PrepareData()
		gocq.LoginInteract()
		_ = terminal.DisableQuickEdit()
		_ = terminal.EnableVT100()
		gocq.WaitSignal()
		_ = terminal.RestoreInputMode()
		log.Exit(0)
	}
	userName := getInput("请输入QQ:\t")
	passWd := getInput("请输入密码:\t")
	serverIp := getInput("请输入目标服务器IP，以逗号分割，默认值为127.0.0.1:\n服务器后以|分割可以设置一个别名，如\"127.0.0.1|本机")
	userNameInt, err := strconv.Atoi(userName)
	if err != nil {
		panic(err)
	}
	loginConfig := entity.LoginConfig{UserName: int64(userNameInt), Passwd: passWd, ServerIp: entity.MakeServerName(strings.Split(serverIp, ",")), ServerRetryTime: 60}
	filePtr, err := os.Create(entity.LOGIN_CONFIG_FILE)
	if err != nil {
		log.Fatalf("文件创建失败%s", err.Error())
	}
	defer func(filePtr *os.File) {
		err := filePtr.Close()
		if err != nil {
			panic(err)
		}
	}(filePtr)
	// 创建Json编码器
	encoder := json.NewEncoder(filePtr)
	err = encoder.Encode(loginConfig)
	terminal.SetTitle()
	gocq.InitBase()
	gocq.PrepareData()
	gocq.LoginInteract()
	_ = terminal.DisableQuickEdit()
	_ = terminal.EnableVT100()
	gocq.WaitSignal()
	_ = terminal.RestoreInputMode()
}
