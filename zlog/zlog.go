package zlog

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/rs/zerolog"
)

/*
time 时间 -
hostname 主机名 -
pid 程序pid -
connId 连接id
uid 用户id
token 用户token
service 服务名 -
func 函数名 -
caller 位置 -
err 报错内容 -
msg 描述信息 -
level 日至级别 -
*/

var (
	logDir   = "/dev/temp/"
	pid      = 0
	service  = ""
	hostname = ""
	stdout   = stdoutToTerminal
	logsMap  sync.Map

	//
	dunno     = "???"
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

const (
	// 日志输出类型
	stdoutToTerminal = 0 // 终端
	stdoutToFile     = 1 // 文件
	stdoutToES       = 2 // elasticSearch

	// 日志级别

	// 日志前缀
	LogInfoSuffix = ".info.log"
	LogErrSuffix  = ".error.log"
	LogWarnSuffix = ".warn.log"
	LogGinAccess  = ".ginAccess.log"
)

// LogT :
type LogT struct {
	lastTime string   // 最近创建日志文件的时间,用于对比是否需要切割日志文件,当前每天切割一次
	filepath string   // 日志文件的路径
	file     *os.File // 日志文件的句柄
}

func InitLog(srcName, dir string, stdType, level int) (err error) {
	if len(dir) != 0 {
		logDir = dir
	}
	if len(srcName) != 0 {
		service = srcName
	}
	hostname, err = os.Hostname()
	if err != nil {
		return
	}
	stdout = stdType
	pid = os.Getpid()
	zerolog.SetGlobalLevel(zerolog.Level(level))
	zerolog.TimestampFunc = time.Now().UTC
	zerolog.ErrorFieldName = "err"
	zerolog.MessageFieldName = "msg"
	fmt.Printf("服务:%s的日志文件路径为:%s\n", service, dir)
	return
}

// 日志级别分别创建日志句柄 fileTail 日志类型的后缀名
func newLog(fileTail string) *zerolog.Logger {
	switch stdout {
	case stdoutToTerminal:
		logger := zerolog.New(os.Stdout).With().Time("time", time.Now().UTC()).Logger()
		return &logger
	case stdoutToFile:
		file := newFile(fileTail)
		logger := zerolog.New(file).With().Time("time", time.Now().UTC()).Caller().Logger().Output(file)
		return &logger
	case stdoutToES:
	}

	logger := zerolog.New(os.Stdout).With().Time("time", time.Now().UTC()).Logger()
	return &logger
}

// 创建文件句柄,用于日志写入文件
func newFile(fileTail string) *os.File {
	// 每天创建一个日志文件
	timeStr := time.Unix(time.Now().Unix(), 0).Format("20060102")
	logT, ok := logsMap.Load(fileTail)
	if ok && logT.(*LogT) != nil && timeStr == logT.(*LogT).lastTime {
		return logT.(*LogT).file
	}

	if ok && logT.(*LogT) != nil {
		err := logT.(*LogT).file.Close()
		if err != nil {
			fmt.Printf("日志文件:%s关闭失败:%v\n", fileTail, err)
		}
	}

	filepath := logDir + "/" + service + "." + hostname + timeStr + fileTail
	ff, err := os.OpenFile(filepath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0775)
	if err != nil {
		fmt.Printf("日志文件:%s创建失败:%v\n", fileTail, err)
	}

	logNode := &LogT{lastTime: timeStr, filepath: filepath, file: ff}
	logsMap.Store(fileTail, logNode)

	return ff
}

// 打印错误级别日志
func Err(err error, msg string) {
	pc, filePath, line, _ := runtime.Caller(1)
	caller := filePath + ":" + strconv.Itoa(line)
	newLog(LogErrSuffix).Error().Err(err).
		Str("service", service).
		Str("pid", strconv.Itoa(pid)).
		Str("func", getFuncName(pc)).
		Str("caller", caller).
		Msg(msg)
}

// 可自定义添加字段和描述
func ErrWithStr(err error) *zerolog.Event {
	pc, filePath, line, _ := runtime.Caller(1)
	caller := filePath + ":" + strconv.Itoa(line)
	return newLog(LogErrSuffix).Error().Err(err).
		Str("service", service).
		Str("pid", strconv.Itoa(pid)).
		Str("func", getFuncName(pc)).
		Str("caller", caller)
}

func Warn(msg string) {
	pc, filePath, line, _ := runtime.Caller(1)
	caller := filePath + ":" + strconv.Itoa(line)
	newLog(LogWarnSuffix).Warn().
		Str("service", service).
		Str("pid", strconv.Itoa(pid)).
		Str("func", getFuncName(pc)).
		Str("caller", caller).Msg(msg)
}

func WarnWithStr() *zerolog.Event {
	pc, filePath, line, _ := runtime.Caller(1)
	caller := filePath + ":" + strconv.Itoa(line)
	return newLog(LogWarnSuffix).Warn().
		Str("service", service).
		Str("pid", strconv.Itoa(pid)).
		Str("func", getFuncName(pc)).
		Str("caller", caller)
}

func Info(msg string) {
	pc, filePath, line, _ := runtime.Caller(1)
	caller := filePath + ":" + strconv.Itoa(line)
	newLog(LogInfoSuffix).Info().
		Str("service", service).
		Str("pid", strconv.Itoa(pid)).
		Str("func", getFuncName(pc)).
		Str("caller", caller).Msg(msg)
}

func InfoWithStr() *zerolog.Event {
	pc, filePath, line, _ := runtime.Caller(1)
	caller := filePath + ":" + strconv.Itoa(line)
	return newLog(LogInfoSuffix).Info().
		Str("service", service).
		Str("pid", strconv.Itoa(pid)).
		Str("func", getFuncName(pc)).
		Str("caller", caller)
}

func ginRequest(clientIp, remoteAddr, uri, method, state, latencyTime string, err error) {
	newLog(LogGinAccess).Info().
		Str("client_ip", clientIp).
		Str("remote_addr", remoteAddr).
		Str("uri", uri).
		Str("method", method).
		Str("status", state).
		Str("latency_time", latencyTime).
		Err(err).Send()
}

func GinLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		// 处理请求
		c.Next()
		// 执行时间
		latencyTime := time.Now().Sub(startTime).String()
		ginRequest(c.ClientIP(), c.Request.RemoteAddr, c.Request.RequestURI, c.Request.Method,
			strconv.Itoa(c.Writer.Status()), latencyTime, c.Err())
	}
}

// 还原函数名称
func getFuncName(pc uintptr) string {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	if lastSlash := bytes.LastIndex(name, slash); lastSlash >= 0 {
		name = name[lastSlash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return string(name)
}
