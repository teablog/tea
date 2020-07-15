package config

import (
	ini "gopkg.in/ini.v1"
	"log"
	"os"
	"strings"
)

var Config *ini.File
var logFd = os.Stdout


func Init(conf string) *ini.File {
	var (
		err error
	)
	Config, err = ini.Load(conf)
	if err != nil {
		log.Fatalf("config load filed, %s", err)
		return nil
	}
	return Config
}

func GetKey(key string) *ini.Key {
	parts := strings.Split(key, "::")
	section := parts[0]
	keyStr := parts[1]
	return Config.Section(section).Key(keyStr)
}

func GetLogFD() *os.File {
	return logFd
}
