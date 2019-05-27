package mggo

import (
    "log"
    "os"
    "path/filepath"
    "runtime"
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
func LogInfo(values ...interface{}) {
    go func() {
        logInfo.Println(values...)
    }()
}

// LogError - send message error in logerror.log
func LogError(values ...interface{}) {
    logError.Println("Unrecovered Error:")
    logError.Println(values...)
    logError.Println("Stack Trace:")
    buf2 := make([]byte, 4096)
    buf2 = buf2[:runtime.Stack(buf2, false)]
    logError.Printf("%s\n", buf2)
}
