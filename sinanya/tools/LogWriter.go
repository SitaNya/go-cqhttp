package tools

import "github.com/Mrs4s/go-cqhttp/sinanya/entity"

type LogWriter struct {
}

func (receiver LogWriter) Write(p []byte) (n int, err error) {
	entity.LOG_CHANNEL <- string(p)
	return len(p), nil
}
