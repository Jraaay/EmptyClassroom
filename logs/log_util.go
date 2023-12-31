package logs

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

const (
	LogIDKey = "K_LOGID"
)

func Init(isMain bool) {
	if isMain {
		stdoutWriter := os.Stdout
		os.Mkdir("run_log", 0755)
		fileWriter, err := os.OpenFile("run_log/ec.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			log.Fatalf("create file log.txt failed: %v", err)
		}
		defaultLogger = log.New(io.MultiWriter(stdoutWriter, fileWriter), "", log.Ldate|log.Lmicroseconds)
	} else {
		stdoutWriter := os.Stdout
		defaultLogger = log.New(io.MultiWriter(stdoutWriter), "", log.Ldate|log.Lmicroseconds)
	}
}

func SetNewContextForGinContext(c *gin.Context) {
	newCtx := GenNewContext()
	c.Set("ctx", newCtx)
	c.Writer.Header().Set("LogID", GetLogIDFromContext(newCtx))
}

func FillZeroForInt(i int, w int) string {
	rawStr := fmt.Sprintf("%d", i)
	for len(rawStr) < w {
		rawStr = "0" + rawStr
	}
	return rawStr
}

func RandomHex(n int) string {
	bytes := make([]byte, n)
	_, _ = rand.Read(bytes)
	return hex.EncodeToString(bytes)

}

func GenLogID() string {
	str := ""
	t := time.Now()
	year, month, day := t.Date()
	hour, minute, second := t.Clock()
	str += FillZeroForInt(year, 4)
	str += FillZeroForInt(int(month), 2)
	str += FillZeroForInt(day, 2)
	str += FillZeroForInt(hour, 2)
	str += FillZeroForInt(minute, 2)
	str += FillZeroForInt(second, 2)
	str += strings.ToUpper(RandomHex(9))
	return str
}

func GenNewContext() context.Context {
	return context.WithValue(context.Background(), LogIDKey, GenLogID())
}

func GetLogIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	v := ctx.Value(LogIDKey)
	if v == nil {
		return ""
	}

	return v.(string)
}

func GetContextFromGinContext(c *gin.Context) context.Context {
	return c.MustGet("ctx").(context.Context)
}
