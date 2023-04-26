package windows

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/Mrs4s/go-cqhttp/cmd/gocq"
	"github.com/Mrs4s/go-cqhttp/sinanya/entity"
	"github.com/Mrs4s/go-cqhttp/sinanya/service"
	"github.com/Mrs4s/go-cqhttp/theme"
	log "github.com/sirupsen/logrus"
	"golang.org/x/image/colornames"
	"os"
	"strconv"
	"strings"
	"time"
)

func LoginWindowsOrMac() {
	version := "0.1.1"
	//新建一个app
	a := entity.WINDOW
	//设置窗口栏，任务栏图标
	a.Settings().SetTheme(&theme.MyTheme{})
	//新建一个窗口
	loginWindow := a.NewWindow(fmt.Sprintf("SinaNyaX 登录界面 %s", version))
	loginWindow.SetCloseIntercept(func() {
		a.Quit()
	})

	//主界面框架布局
	MainShow(loginWindow, a)
	//尺寸
	loginWindow.Resize(fyne.Size{Width: 500, Height: 100})
	//w居中显示
	loginWindow.CenterOnScreen()
	//循环运行
	loginWindow.Show()
	a.Run()

	err := os.Unsetenv("FYNE_FONT")
	if err != nil {
		return
	}
}

type LogSave struct {
	logList []string
}

func (receiver *LogSave) add(item string) {
	if len(receiver.logList) > 50 {
		receiver.logList = receiver.logList[:1]
	}
	receiver.logList = append(receiver.logList, strings.Trim(item, "\n"))
}

type serverStatusList struct {
	ip     string
	label  *canvas.Text
	button *widget.Button
}

func RunningWindow(app fyne.App) {
	loginConfig := entity.LoginConfig{}
	content, err := os.ReadFile(entity.LOGIN_CONFIG_FILE)
	if err == nil {
		_ = json.NewDecoder(strings.NewReader(string(content))).Decode(&loginConfig)
	}
	runningWindow := app.NewWindow(fmt.Sprintf("SinaNyaX 日志界面:\t%d", loginConfig.UserName))
	runningWindow.SetCloseIntercept(func() {
		app.Quit()
	})
	runningWindow.Show()
	logSave := LogSave{}
	txtResults := widget.NewLabel("")
	cntScrolling := container.NewScroll(txtResults)
	serverButtonList := container.NewGridWithColumns(3)
	var statusList []serverStatusList
	for ip, name := range loginConfig.ServerIp {
		button := widget.NewButton("重连", func(ip string) func() {
			return func() {
				service.CreateClient(int(loginConfig.UserName), ip, 60083)
			}
		}(ip))
		status := canvas.NewText("离线", colornames.Green)
		statusList = append(statusList, serverStatusList{ip: ip, label: status, button: button})
		serverButtonList.Add(widget.NewLabel(name))
		serverButtonList.Add(status)
		serverButtonList.Add(button)
	}
	//txtResults.ShowLineNumbers = true
	cnt1 := container.NewGridWrap(fyne.NewSize(600, 600), cntScrolling)
	cnt2 := container.NewGridWrap(fyne.NewSize(600, float32(30*len(statusList))), serverButtonList)
	cntContent := container.NewVBox(cnt2, cnt1)

	runningWindow.SetContent(cntContent)
	go func() {
		for true {
			time.Sleep(500 * time.Millisecond)
			strMessage := <-entity.LOG_CHANNEL
			logSave.add(strMessage)
			txtResults.SetText(strings.Join(logSave.logList, "\n"))
		}
	}()
	for true {
		time.Sleep(5 * time.Second)
		for _, status := range statusList {
			enabled := getServerStatus(status.ip)
			if enabled {
				status.label.Text = "在线"
				status.label.Color = colornames.Green
			} else {
				status.label.Text = "离线"
				status.label.Color = colornames.Red
			}
			if getServerStatus(status.ip) {
				status.button.Disable()
			} else {
				status.button.Enable()
			}
		}
		//txtResults.SetText(strMessage)
		cntContent.Refresh()
	}
}

func getServerStatus(ip string) bool {
	if _, ok := entity.ChannelList[ip]; ok {
		return true
	} else {
		return false
	}
}

// MainShow 主界面函数
func MainShow(w fyne.Window, app fyne.App) {
	loginConfig := entity.LoginConfig{}
	content, err := os.ReadFile(entity.LOGIN_CONFIG_FILE)
	if err == nil {
		_ = json.NewDecoder(strings.NewReader(string(content))).Decode(&loginConfig)
	}
	//var ctrl *beep.Ctrl
	title := widget.NewLabel("SinaNyaX 登录界面")
	userNameLabel := widget.NewLabel("QQ:")
	userNameEntry := widget.NewEntry()
	passWdLabel := widget.NewLabel("密码:")
	passWd := widget.NewPasswordEntry()
	serverIpLabel := widget.NewLabel("服务器列表:")
	serverIp := widget.NewMultiLineEntry()
	userNameEntry.Text = strconv.FormatInt(loginConfig.UserName, 10)
	passWd.Text = loginConfig.Passwd
	passWd.Password = true
	serverIp.Text = strings.Join(entity.MakeMapToServersIp(loginConfig.ServerIp), "\n")
	dia1 := widget.NewButton("登录", func() { //回调函数：打开选择文件对话框
		userName, err := strconv.Atoi(userNameEntry.Text)
		if err != nil {
			log.Fatalf("您输入的账号:%s不是有效的QQ号", userNameEntry.Text)
		}
		loginConfig := entity.LoginConfig{UserName: int64(userName), Passwd: passWd.Text, ServerIp: entity.MakeServerName(strings.Split(serverIp.Text, "\n")), ServerRetryTime: 60}
		_, err = os.Stat(entity.LOGIN_CONFIG_FILE)
		if err == nil {
			err := os.Remove(entity.LOGIN_CONFIG_FILE)
			if err != nil {
				panic(err)
			}
		}
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
		if err != nil {
			log.Fatalf("登录失败%s", err.Error())
			w.Close()
		} else {
			log.Infof("登录成功")
			w.Hide()
		}
		go func() {
			//terminal.SetTitle()
			//gocq.InitBase()
			//gocq.PrepareData()
			//gocq.LoginInteract()
			//_ = terminal.DisableQuickEdit()
			//_ = terminal.EnableVT100()
			//gocq.WaitSignal()
			//_ = terminal.RestoreInputMode()
			gocq.InitBase()
			gocq.PrepareData()
			gocq.LoginInteract()
			gocq.WaitSignal()
		}()
		go func() { RunningWindow(app) }()
	})

	head := container.NewCenter(title)
	userNameEntry.CursorColumn = 2
	userNameRow := container.New(layout.NewBorderLayout(nil, nil, userNameLabel, nil), userNameLabel, userNameEntry)
	passWdRow := container.New(layout.NewBorderLayout(nil, nil, passWdLabel, nil), passWdLabel, passWd)
	serverIpRow := container.New(layout.NewBorderLayout(nil, nil, serverIpLabel, nil), serverIpLabel, serverIp)
	dia1.Resize(fyne.NewSize(40, 120))
	button := container.NewCenter(dia1)

	contentBox := container.NewVBox(head, userNameRow, passWdRow, serverIpRow, button) //控制显示位置顺序
	w.SetContent(contentBox)
}
