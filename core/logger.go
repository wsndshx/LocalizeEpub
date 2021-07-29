package core

import (
	"io"
	"io/fs"
	"log"
	"os"
)

var (
	LogInfo  *log.Logger
	LogError *log.Logger
)

type Test struct {
	data []byte // 容纳数据的地方
}

func init() {
	// 创建日志文件
	err := LogFS.WriteFile("log", []byte("Nya~\n"), 0755)
	if err != nil {
		log.Fatalln("Faild to open error logger file:", err)
	}
	//自定义日志格式
	Log := &Test{} // 声明实例
	LogInfo = log.New(io.MultiWriter(Log, os.Stderr), "[INFO] ", log.Ldate|log.Ltime)
	LogError = log.New(io.MultiWriter(Log, os.Stderr), "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
}

func (t *Test) Write(p []byte) (n int, err error) {
	// 读取旧文件
	data, err := fs.ReadFile(LogFS, "log")
	oldLen := len(data)
	// 插入数据
	data = append(data, p...)
	// 回写文件
	err = LogFS.WriteFile("log", data, 0755)
	newLen := len(data)
	n = newLen - oldLen
	return
}
