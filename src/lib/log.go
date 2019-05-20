package lib

import (
	"log"
	"os"
	"io"
	"time"
	
)

var Error *log.Logger
var debug bool

func Logger() *log.Logger{

	if Error != nil {
		return Error;
	}

	config := NewConfig()
	var logPath string = config.Log.Path
	var name string = "err_log"+time.Now().Format("2006-01-02") + ".log"

	//日志
	file, err := os.OpenFile(logPath + name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open error log file:", err)
		return nil
	}
	return log.New(io.MultiWriter(file, os.Stderr), "Error", log.Ltime|log.Lshortfile)
}