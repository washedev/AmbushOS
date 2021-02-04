package log

import (
	"fmt"
	"github.com/fatih/color"
	"sync"
	"time"
)

var (
	lock = sync.Mutex{}
)

func Print(s string) {
	lock.Lock()
	defer lock.Unlock()

	fmt.Fprint(color.Output, s)
}

func format(s, l, t string, c func(string, ...interface{}) string) string {
	f := c("[%v] [%v] %v", l, t, s)
	return fmt.Sprintf("[%v] %v", whiteColor(time.Now().Format("15:04:05.000")), f)
}

func formatln(s, l, t string, c func(string, ...interface{}) string) string {
	f := c("[%v] [%v] %v\n", l, t, s)
	return fmt.Sprintf("[%v] %v", whiteColor(time.Now().Format("15:04:05.000")), f)
}

func Debug(s, t string) {
	s = format(s, "DBG", t, DebugColor)
	Print(s)
}

func Info(s, t string) {
	s = format(s, "INF", t, InfoColor)
	Print(s)
}

func Warning(s, t string) {
	s = format(s, "WRN", t, WarningColor)
	Print(s)
}

func Warn(s, t string) {
	s = format(s, "WRN", t, WarningColor)
	Print(s)
}

func Error(s, t string) {
	s = format(s, "ERR", t, ErrorColor)
	Print(s)
}

func Debugln(s, t string) {
	s = formatln(s, "DBG", t, DebugColor)
	Print(s)
}

func Infoln(s, t string) {
	s = formatln(s, "INF", t, InfoColor)
	Print(s)
}

func Warningln(s, t string) {
	s = formatln(s, "WRN", t, WarningColor)
	Print(s)
}

func Warnln(s, t string) {
	s = formatln(s, "WRN", t, WarningColor)
	Print(s)
}

func Errorln(s, t string) {
	s = formatln(s, "ERR", t, ErrorColor)
	Print(s)
}
