package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// 设置颜色刷
type Brush func(string) string

func NewBrush(color string) Brush {
	pre := "\033["
	reset := "\033[0m"
	return func(text string) string {
		return pre + color + "m" + text + reset
	}
}

var colors = map[Level]Brush{
	TRACE:   NewBrush("1;32"), // Trace      cyan
	INFO:    NewBrush("1;34"), // Info		blue
	WARNING: NewBrush("1;33"), // Warning    yellow
	DEBUG:   NewBrush("1;36"), // Debug
	ERROR:   NewBrush("1;31"), // Error      red
	FATAL:   NewBrush("1;37"), // Fatal		white

}

type ConsoleLog struct {
	log   *log.Logger
	level Level
}

// 初始化控制台输出引擎
func NewConsole() LogEngine {
	return &ConsoleLog{
		log:   log.New(os.Stdout, "", log.Ldate|log.Ltime),
		level: DEBUG,
	}
}

func (c *ConsoleLog) Init(conf string) error {
	if len(conf) == 0 {
		return nil
	}

	return json.Unmarshal([]byte(conf), c)
}

func (c *ConsoleLog) Write(msg string, level Level) error {
	if level < c.level {
		return nil
	}
	c.log.Println(colors[level](msg))
	return nil
}

func (c *ConsoleLog) Destroy() {}
func (c *ConsoleLog) Flush()   {}

func init() {
	Register("console", NewConsole)
}

// DEBUG

func Debug(v ...interface{}) {
	msg := fmt.Sprint("[DEBUG] " + fmt.Sprintln(v...))
	log.Println(colors[DEBUG](msg))
}

func Debugf(format string, v ...interface{}) {
	msg := fmt.Sprintf("[DEBUG] "+format, v...)
	log.Println(colors[DEBUG](msg))
}

// Trace
func Trac(v ...interface{}) {
	msg := fmt.Sprint("[TRAC] " + fmt.Sprintln(v...))
	log.Println(colors[TRACE](msg))
}

func Tracf(format string, v ...interface{}) {
	msg := fmt.Sprintf("[TRAC] "+format, v...)
	log.Println(colors[TRACE](msg))
}

// INFO
func Info(v ...interface{}) {
	msg := fmt.Sprint("[INFO] " + fmt.Sprintln(v...))
	log.Println(colors[INFO](msg))
}

func Infof(format string, v ...interface{}) {
	msg := fmt.Sprintf("[INFO] "+format, v...)
	log.Println(colors[INFO](msg))
}

//WARNING
func Warn(v ...interface{}) {
	msg := fmt.Sprint("[WARN] " + fmt.Sprintln(v...))
	log.Println(colors[WARNING](msg))
}

func Warnf(format string, v ...interface{}) {
	msg := fmt.Sprintf("[WARN] "+format, v...)
	log.Println(colors[WARNING](msg))
}

// ERROR
func Error(v ...interface{}) {
	msg := fmt.Sprint("[ERRO] " + fmt.Sprintln(v...))
	log.Println(colors[ERROR](msg))
}

func Errorf(format string, v ...interface{}) {
	msg := fmt.Sprintf("[ERRO] "+format, v...)
	log.Println(colors[ERROR](msg))
}

// FATAL
func Fatal(v ...interface{}) {
	msg := fmt.Sprintf("[FATA] " + fmt.Sprintln(v...))
	log.Println(colors[FATAL](msg))
}

func Fatalf(format string, v ...interface{}) {
	msg := fmt.Sprintf("[FATA] "+format, v...)
	log.Println(colors[FATAL](msg))
}
