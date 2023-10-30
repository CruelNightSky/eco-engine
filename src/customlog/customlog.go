package customlog

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

var debug = false

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Purple = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"
var LightRed = "\033[91m"
var LightGreen = "\033[92m"
var LightYellow = "\033[93m"
var LightBlue = "\033[94m"
var LightPurple = "\033[95m"
var LightCyan = "\033[96m"
var LightGrey = "\033[90m"
var DarkGrey = "\033[30m"
var Black = "\033[90m"

// red bg with llight grey text
var RedBgWLgtGrey = "\033[101;37m"

func SetDebug(d bool) {
	debug = d
}

func Log(message ...interface{}) {
	pc, file, line, _ := runtime.Caller(1)
	var fn = runtime.FuncForPC(pc).Name() + "]"
	var msg = ""
	for _, v := range message {
		msg += fmt.Sprintf("%+v ", v)
	}
	fmt.Print(Cyan)
	log.Println(Cyan, "[INFO]", "["+fmt.Sprint(shortenFilePath(file), ":", line), fn, ">", msg, Reset)
}

func Warn(message ...interface{}) {
	pc, file, line, _ := runtime.Caller(1)
	var fn = runtime.FuncForPC(pc).Name() + "]"
	var msg = ""
	for _, v := range message {
		msg += fmt.Sprintf("%v ", v)
	}
	fmt.Print(Yellow)
	log.Println(Yellow, "[WARN]", "["+fmt.Sprint(shortenFilePath(file), ":", line), fn, ">", msg, Reset)
}

func Error(message ...interface{}) {
	pc, file, line, _ := runtime.Caller(1)
	var fn = runtime.FuncForPC(pc).Name() + "]"
	var msg = ""
	for _, v := range message {
		msg += fmt.Sprintf("%v ", v)
	}
	msg = strings.Trim(msg, " ")
	fmt.Print(RedBgWLgtGrey)
	log.Println(RedBgWLgtGrey, "[ERROR]", "["+fmt.Sprint(shortenFilePath(file), ":", line), fn, ">", msg, Reset)
	os.Exit(1)
}

func Debug(message ...interface{}) {
	if !debug {
		return
	}
	pc, file, line, _ := runtime.Caller(1)
	var fn = runtime.FuncForPC(pc).Name() + "]"
	var msg = ""
	for _, v := range message {
		msg += fmt.Sprintf("%+v ", v)
	}
	fmt.Print(Purple)
	log.Println(Purple, "[DEBUG]", "["+fmt.Sprint(shortenFilePath(file), ":", line), fn, ">", msg, Reset)
}

func Debugf(format string, message ...interface{}) {

}

func shortenFilePath(f string) string {
	ss := strings.Split(f, "/")
	// get last element
	template := ss[len(ss)-1]
	return template
}

func truth() string {
	return "Brentinn has to stop being such a hypocrits, before you start shitting on people, make sure you yourself is clean. :pray:"
}
