package logger

import (
	"fmt"
	"log"
	"path"
	"runtime"
	"strings"
	"sync"
)

func init() {
	log.SetFlags(log.LstdFlags)
}

type Level byte

const (
	DEBUG Level = iota + 1
	TRACE
	INFO
	WARNING
	ERROR
	FATAL
)

// log输出接口
type LogEngine interface {
	Init(conf string) error              //初始化
	Write(msg string, level Level) error //写入
	Destroy()
	Flush()
}

// log结构体
type Log struct {
	level         Level
	msg           chan *logMsg
	trackFuncCall bool //是否追踪调用函数
	funcCallDepth int
	output        map[string]LogEngine
	lock          sync.Mutex
}

// log内容
type logMsg struct {
	level Level
	msg   string
}

// 定义输出引擎字典
type engineType func() LogEngine

var engines = make(map[string]engineType)

// 注册引擎
func Register(name string, log engineType) {
	if log == nil {
		panic("logs: Register provide is nil")
	}
	if _, dup := engines[name]; dup {
		panic("logs: Register called twice for provider " + name)
	}
	engines[name] = log
}

// 初始化log
// output -- 适配名称 为空(默认)console
// chanlen -- 缓存大小
func NewLog(chanlen uint64) *Log {
	l := &Log{
		level:         DEBUG,
		trackFuncCall: false,
		funcCallDepth: 2,
		msg:           make(chan *logMsg, chanlen),
		output:        make(map[string]LogEngine),
	}

	l.SetEngine("console", "")

	return l
}

// 设置log等级
func (l *Log) SetLevel(lstr string) *Log {
	var level Level

	switch lstr {
	case "D", "Debug", "debug":
		level = DEBUG
	case "T", "Trace", "trace", "TRACE", "trac", "Trac", "TRAC":
		level = TRACE
	case "I", "Info", "info", "INFO":
		level = INFO
	case "W", "Warning", "warning", "WARNING", "Warn", "warn", "WARN":
		level = WARNING
	case "E", "Error", "error", "ERROR":
		level = ERROR
	case "F", "Fatal", "fatal", "FATAL":
		level = FATAL
	case "":
	default:
		level = DEBUG
	}
	l.level = level

	return l
}

// 设置是否输出行号
func (l *Log) SetFuncCall(bool) *Log {

	l.trackFuncCall = true

	return l
}

// 设置是否输出行号
func (l *Log) SetFuncCallDepth(depth int) *Log {
	l.funcCallDepth = depth

	return l
}

// 设置输出引擎
func (l *Log) SetEngine(engname string, conf string) *Log {

	l.lock.Lock()
	defer l.lock.Unlock()

	//获取引擎
	if log, ok := engines[engname]; ok {
		lg := log()
		err := lg.Init(conf)
		if err != nil {
			errmsg := fmt.Errorf("SetEngine error: %s", err)
			fmt.Println(errmsg)
			return nil
		}

		l.output[engname] = lg
	} else {
		fmt.Printf("unknown Enginee %q ", engname)
		return nil
	}

	return l
}

// 删除不希望使用的引擎
func (l *Log) DelEngine(engname string) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	if lg, ok := l.output[engname]; ok {
		lg.Destroy()
		delete(l.output, engname)
		return nil
	} else {
		return fmt.Errorf("unknown engine name %q (forgotten Register?)", engname)
	}
}

// 初始化logMsg
func (l *Log) newMsg(msg string, level Level) {
	l.lock.Lock()
	defer l.lock.Unlock()

	lm := new(logMsg)
	lm.level = level

	if l.trackFuncCall {
		_, file, line, ok := runtime.Caller(l.funcCallDepth)
		if !ok {
			file = "???"
			line = 0
		}
		_, filename := path.Split(file)
		lm.msg = fmt.Sprintf("[%s:%d] %s", filename, line, msg)
	} else {
		lm.msg = msg
	}

	l.msg <- lm

}

// 写入
func (l *Log) write() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("The Logger's write() catch panic: %v\n", err)
		}
	}()

	lm := <-l.msg
	for name, e := range l.output {
		err := e.Write(lm.msg, lm.level)
		if err != nil {
			fmt.Println("ERROR, unable to WriteMsg:", name, err)
		}
	}
}

// 获取调用的位置
func (l *Log) getInvokerLocation() string {
	//runtime.Caller(skip) skip=0 返回当前调用Caller函数的函数名、文件、程序指针PC，1是上一层函数，以此类推
	pc, file, line, ok := runtime.Caller(l.funcCallDepth)
	if !ok {
		return ""
	}
	fname := ""
	if index := strings.LastIndex(file, "/"); index > 0 {
		fname = file[index+1 : len(file)]
	}
	funcPath := ""
	funcPtr := runtime.FuncForPC(pc)
	if funcPtr != nil {
		funcPath = funcPtr.Name()
	}
	return fmt.Sprintf("%s : [%s:%d]", funcPath, fname, line)
}

// DEBUG
func (l *Log) Debug(v ...interface{}) {
	if l.level > DEBUG {
		return
	}
	msg := fmt.Sprint("[DEBUG] " + fmt.Sprintln(v...))
	l.newMsg(msg, DEBUG)
	l.write()
}

func (l *Log) Debugf(format string, v ...interface{}) {
	if l.level > DEBUG {
		return
	}
	msg := fmt.Sprintf("[DEBUG] "+format, v...)
	l.newMsg(msg, DEBUG)
	l.write()
}

// Trace
func (l *Log) Trac(v ...interface{}) {
	// 如果设置的级别比 trace 级别高,不输出
	if l.level > TRACE {
		return
	}
	msg := fmt.Sprint("[TRAC] " + fmt.Sprintln(v...))
	l.newMsg(msg, TRACE)
	l.write()
}

func (l *Log) Tracf(format string, v ...interface{}) {
	// 如果设置的级别比 trace 级别高,不输出
	if l.level > TRACE {
		return
	}
	msg := fmt.Sprintf("[TRAC] "+format, v...)
	l.newMsg(msg, TRACE)
	l.write()
}

// INFO
func (l *Log) Info(v ...interface{}) {
	if l.level > INFO {
		return
	}
	msg := fmt.Sprint("[INFO] " + fmt.Sprintln(v...))
	l.newMsg(msg, INFO)
	l.write()
}

func (l *Log) Infof(format string, v ...interface{}) {
	if l.level > INFO {
		return
	}
	msg := fmt.Sprintf("[INFO] "+format, v...)
	l.newMsg(msg, INFO)
	l.write()
}

//WARNING
func (l *Log) Warn(v ...interface{}) {
	if l.level > WARNING {
		return
	}
	msg := fmt.Sprint("[WARN] " + fmt.Sprintln(v...))
	l.newMsg(msg, WARNING)
	l.write()
}

func (l *Log) Warnf(format string, v ...interface{}) {
	if l.level > WARNING {
		return
	}
	msg := fmt.Sprintf("[WARN] "+format, v...)
	l.newMsg(msg, WARNING)
	l.write()
}

// ERROR
func (l *Log) Error(v ...interface{}) {
	if l.level > ERROR {
		return
	}
	msg := fmt.Sprint("[ERRO] " + fmt.Sprintln(v...))
	l.newMsg(msg, ERROR)
	l.write()
}

func (l *Log) Errorf(format string, v ...interface{}) {
	if l.level > ERROR {
		return
	}
	msg := fmt.Sprintf("[ERRO] "+format, v...)
	l.newMsg(msg, ERROR)
	l.write()
}

// FATAL
func (l *Log) Fatal(v ...interface{}) {
	if l.level > FATAL {
		return
	}
	msg := fmt.Sprintf("[FATA] " + fmt.Sprintln(v...))
	l.newMsg(msg, FATAL)
	l.write()
}

func (l *Log) Fatalf(format string, v ...interface{}) {
	if l.level > FATAL {
		return
	}
	msg := fmt.Sprintf("[FATA] "+format, v...)
	l.newMsg(msg, FATAL)
	l.write()
}

func (l *Log) Close() {
	for {
		if len(l.msg) > 0 {
			bm := <-l.msg
			for _, l := range l.output {
				err := l.Write(bm.msg, bm.level)
				if err != nil {
					fmt.Println("ERROR, unable to WriteMsg (while closing logger):", err)
				}
			}
			continue
		}
		break
	}
	for _, ls := range l.output {
		ls.Flush()
		ls.Destroy()
	}
}
