package logger

import (
	"encoding/json"
	"fmt"
	"github.com/dougEfresh/logzio-go"
	"github.com/fatih/color"
	"log"
	"srv-con-app-golang/app"
	"strings"
	"time"
)

type H map[string]interface{}

func LogZ(message string, h H) {

	h["message"] = message

	msg, err := json.Marshal(h)
	if err != nil {
		log.Println(err)
	}

	l, err := logzio.New(app.CFG.LogZKey)
	if err != nil {
		panic(err)
	}
	defer l.Stop()

	err = l.Send(msg)
	if err != nil {
		panic(err)
	}
}

func Info(data ...string) {
	color.Cyan(format(data...))
}

func Log(data ...string) {
	color.White(format(data...))
}

func Warning(data ...string) {
	color.Yellow(format(data...))
}

func Error(data ...string) {
	color.Red(format(data...))
}

func format(data ...string) string {
	return fmt.Sprintf("[%s] %s", currentTime(), strings.Join(data, " "))
}

func currentTime() string {
	return time.Now().Format("2006.01.02 15:04:05")
}
