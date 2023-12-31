package logs

import (
	"fmt"
	"golang.org/x/net/context"
	"log"
	"path"
	"runtime"
)

var (
	defaultLogger *log.Logger
)

func GetCallerInfo() (info string) {
	_, file, lineNo, ok := runtime.Caller(2)
	if !ok {
		info = "runtime.Caller() failed"
		return
	}
	fileName := path.Base(file) // Base函数返回路径的最后一个元素
	return fmt.Sprintf("%s:%d", fileName, lineNo)
}

func CtxInfo(ctx context.Context, format string, v ...interface{}) {
	defaultLogger.SetPrefix("[INFO] ")
	logID := GetLogIDFromContext(ctx)
	defaultLogger.Printf(GetCallerInfo()+" "+logID+" "+format, v...)
}

func CtxWarn(ctx context.Context, format string, v ...interface{}) {
	defaultLogger.SetPrefix("[Warn] ")
	logID := GetLogIDFromContext(ctx)
	defaultLogger.Printf(GetCallerInfo()+" "+logID+" "+format, v...)
}

func CtxError(ctx context.Context, format string, v ...interface{}) {
	defaultLogger.SetPrefix("[Error] ")
	logID := GetLogIDFromContext(ctx)
	defaultLogger.Printf(GetCallerInfo()+" "+logID+" "+format, v...)
}
