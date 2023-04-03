// Package config 包含go-cqhttp操作配置文件的相关函数
package config

import (
	_ "embed" // embed the default config file
	"encoding/json"
	"github.com/Mrs4s/go-cqhttp/sinanya/entity"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
	"regexp"
	"strings"
)

// defaultConfig 默认配置文件
//
//go:embed default_config.yml
var defaultConfig string

// Reconnect 重连配置
type Reconnect struct {
	Disabled bool `yaml:"disabled"`
	Delay    uint `yaml:"delay"`
	MaxTimes uint `yaml:"max-times"`
	Interval int  `yaml:"interval"`
}

// Account 账号配置
type Account struct {
	Uin              int64      `yaml:"uin"`
	Password         string     `yaml:"password"`
	Encrypt          bool       `yaml:"encrypt"`
	Status           int        `yaml:"status"`
	ReLogin          *Reconnect `yaml:"relogin"`
	UseSSOAddress    bool       `yaml:"use-sso-address"`
	AllowTempSession bool       `yaml:"allow-temp-session"`
}

// Config 总配置文件
type Config struct {
	Account   *Account `yaml:"account"`
	Heartbeat struct {
		Disabled bool `yaml:"disabled"`
		Interval int  `yaml:"interval"`
	} `yaml:"heartbeat"`

	Message struct {
		PostFormat          string `yaml:"post-format"`
		ProxyRewrite        string `yaml:"proxy-rewrite"`
		IgnoreInvalidCQCode bool   `yaml:"ignore-invalid-cqcode"`
		ForceFragment       bool   `yaml:"force-fragment"`
		FixURL              bool   `yaml:"fix-url"`
		ReportSelfMessage   bool   `yaml:"report-self-message"`
		RemoveReplyAt       bool   `yaml:"remove-reply-at"`
		ExtraReplyData      bool   `yaml:"extra-reply-data"`
		SkipMimeScan        bool   `yaml:"skip-mime-scan"`
		ConvertWebpImage    bool   `yaml:"convert-webp-image"`
	} `yaml:"message"`

	Output struct {
		LogLevel    string `yaml:"log-level"`
		LogAging    int    `yaml:"log-aging"`
		LogForceNew bool   `yaml:"log-force-new"`
		LogColorful *bool  `yaml:"log-colorful"`
		Debug       bool   `yaml:"debug"`
	} `yaml:"output"`

	Servers  []map[string]yaml.Node `yaml:"servers"`
	Database map[string]yaml.Node   `yaml:"database"`
}

// Server 的简介和初始配置
type Server struct {
	Brief   string
	Default string
}

// Parse 从默认配置文件路径中获取
func Parse(path string) *Config {
	loginConfig := entity.LoginConfig{}
	content, err := os.ReadFile(entity.LOGIN_CONFIG_FILE)
	err = json.NewDecoder(strings.NewReader(string(content))).Decode(&loginConfig)
	if err != nil {
		log.Fatal("配置文件不合法!", err)
	}
	config := &Config{}
	configStr := generateConfig()
	err = yaml.NewDecoder(strings.NewReader(configStr)).Decode(config)
	if err != nil {
		log.Fatal("配置文件不合法!", err)
	}
	config.Account.Uin = loginConfig.UserName
	config.Account.Password = loginConfig.Passwd
	config.Output.LogLevel = "info"
	config.Servers = make([]map[string]yaml.Node, 0)
	openFile, _ := os.Create(path)
	err = yaml.NewEncoder(openFile).Encode(config)
	if err != nil {
		panic(err)
	}
	return config
}

var serverconfs []*Server

// AddServer 添加该服务的简介和默认配置
func AddServer(s *Server) {
	serverconfs = append(serverconfs, s)
}

// generateConfig 生成配置文件
func generateConfig() string {
	sb := strings.Builder{}
	sb.WriteString(defaultConfig)
	readString := "0"
	readMax := len(serverconfs)
	if readMax > 10 {
		readMax = 10
	}
	for _, r := range readString {
		r -= '0'
		if r >= 0 && r < rune(readMax) {
			sb.WriteString(serverconfs[r].Default)
		}
	}
	return sb.String()
}

// expand 使用正则进行环境变量展开
// os.ExpandEnv 字符 $ 无法逃逸
// https://github.com/golang/go/issues/43482
func expand(s string, mapping func(string) string) string {
	r := regexp.MustCompile(`\${([a-zA-Z_]+[a-zA-Z0-9_:/.]*)}`)
	return r.ReplaceAllStringFunc(s, func(s string) string {
		s = strings.Trim(s, "${}")
		before, after, ok := strings.Cut(s, ":")
		m := mapping(before)
		if ok && m == "" {
			return after
		}
		return m
	})
}
