package logger

import (
	"fmt"
	"github.com/shatylos/ffmpeg-thumbnails/tools/terminal"
	"log"
	"os"
	"runtime"
	"runtime/debug"
	"time"
)

var isInit bool

func logInit() {
	log.SetOutput(os.Stdout)
	isInit = true
}

func PrintError(err error) {
	if !isInit {
		logInit()
	}
	fmt.Printf("%serror:%s %v\n\n", terminal.ColorRed, terminal.ColorReset, err)
}

func Error(msg string) {
	if !isInit {
		logInit()
	}
	pc, file, line, ok := runtime.Caller(1)
	funcName := ""
	if ok {
		funcName = runtime.FuncForPC(pc).Name()
	}
	stackTrace := string(debug.Stack())
	log.Printf("%sError%s %s [%s:%d %s]\nStack Trace:\n%s",
		terminal.ColorRed, terminal.ColorReset, msg, file, line, funcName, stackTrace)
}

func Warning(msg string) {
	if !isInit {
		logInit()
	}
	pc, file, line, ok := runtime.Caller(1)
	funcName := ""
	if ok {
		funcName = runtime.FuncForPC(pc).Name()
	}
	log.Printf("%sWarning%s %s [%s:%d %s]\n",
		terminal.ColorYellow, terminal.ColorReset, msg, file, line, funcName)
}

func Info(msg string) {
	if !isInit {
		logInit()
	}
	log.Println(terminal.ColorGrey, time.Now().Format("2006-01-02 15:04:05"), msg, terminal.ColorReset)
}

func Success(msg string) {
	if !isInit {
		logInit()
	}
	log.Println(terminal.ColorGreen, time.Now().Format("2006-01-02 15:04:05"), msg, terminal.ColorReset)
}
