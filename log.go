package mggo

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var logInfo *log.Logger
var logError *log.Logger

func init() {
	var err error
	newpath := filepath.Join(".", "logs")
	os.MkdirAll(newpath, os.ModePerm)
	fInfo, err := os.Create(newpath + "/loginfo.log")
	if err != nil {
		log.Fatal(err)
	}
	fError, err := os.Create(newpath + "/logerror.log")
	if err != nil {
		log.Fatal(err)
	}
	logInfo = log.New(fInfo, "", 2)
	logError = log.New(fError, "", 3)
}

// LogInfo - send message in loginfo.log
func LogInfo(ctx *BaseContext, values ...interface{}) {
	go func() {
		var uuid string
		if ctx != nil {
			uuid = ctx.uuid
		}
		logInfo.SetPrefix(time.Now().Format("2006-01-02 15:04:05.000000") + fmt.Sprintf(" [%v] ", uuid))
		logInfo.Println(values...)
	}()
}

// LogError - send message error in logerror.log
func LogError(ctx *BaseContext, values ...interface{}) {
	var uuid string
	if ctx != nil {
		uuid = ctx.uuid
	}
	logError.SetPrefix(time.Now().Format("2006-01-02 15:04:05.000000") + fmt.Sprintf(" [%v] ", uuid))
	logError.Println("Unrecovered Error:")
	logError.Println(values...)
	logError.Println("Stack Trace:")
	buf2 := make([]byte, 4096)
	buf2 = buf2[:runtime.Stack(buf2, false)]
	logError.Printf("%s\n", buf2)
}
