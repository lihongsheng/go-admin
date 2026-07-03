package cron

import (
	log2 "github.com/lihongsheng/pay-gateway/log"
)

type log struct{}

func (l *log) Info(msg string, keysAndValues ...interface{}) {
	log2.Info(msg, keysAndValues...)
}

func (l *log) Error(err error, msg string, keysAndValues ...interface{}) {
	keysAndValues = append(keysAndValues, "error", err)
	log2.Error(msg, keysAndValues...)
}
