package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type SplitType string

const (
	SPLIT_TYPE_BY_SIZE  SplitType = "size"
	SPLIT_TYPE_BY_DAILY SplitType = "daily"
)

const (
	DEFAULT_LEVEL     Level = TRACE
	DEFAULT_FILE_SIZE       = 30
)

type FileLog struct {
	log       *log.Logger
	Level     Level     `json:"level"`
	FileName  string    `json:"filename"`
	MaxSize   int64     `json:"maxsize"` //MB
	SplitType SplitType `json:"split"`   //拆分方式,2种(1-天数,2-文件大小)
	fnum      int

	date time.Time

	lock sync.Mutex

	mw *MuxWriter
}

type MuxWriter struct {
	mu sync.RWMutex
	fd *os.File
}

func (f *MuxWriter) Write(b []byte) (int, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.fd.Write(b)
}

func NewFile() LogEngine {
	t, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	fl := &FileLog{
		Level:     DEFAULT_LEVEL,
		FileName:  "",
		MaxSize:   DEFAULT_FILE_SIZE,
		SplitType: SPLIT_TYPE_BY_SIZE,
		fnum:      0,
		date:      t,
	}

	fl.mw = new(MuxWriter)
	fl.log = log.New(fl.mw, "", log.Ldate|log.Ltime)

	return fl
}

func (l *FileLog) Init(conf string) error {

	if len(conf) != 0 {
		err := json.Unmarshal([]byte(conf), l)
		if err != nil {
			return err
		}
	}

	if len(l.FileName) == 0 {
		name := path.Join("log", "log"+".log")
		l.FileName = name
	}

	if l.MaxSize == 0 {
		l.MaxSize = 30
	}

	return l.initLog()
}

func (l *FileLog) initLog() error {
	fd, err := l.createFile()
	if err != nil {
		return err
	}
	if l.mw.fd != nil {
		l.mw.fd.Close()
	}
	l.mw.fd = fd
	return nil
}

// 检测是否需要重新创建文件
func (l *FileLog) docheck() {
	// 判断log分割类型
	switch l.SplitType {
	case SPLIT_TYPE_BY_SIZE:
		l.splitbysize()
	case SPLIT_TYPE_BY_DAILY:
		l.splitbydaily()
	default:
		l.splitbysize()
	}
}

// 按文件大小分割
func (l *FileLog) splitbysize() {
	l.lock.Lock()
	defer l.lock.Unlock()

	// 检查文件大小
	finfo, err := os.Stat(l.FileName)
	if err != nil {
		fmt.Printf("get %s stat err: %s\n", l.FileName, err)
	}

	if finfo.Size() >= l.MaxSize<<10<<10 {
		l.mw.mu.Lock()
		defer l.mw.mu.Unlock()
		// 关闭之前的fd
		oldfd := l.mw.fd
		oldfd.Close()

		fmt.Printf("文件大小: %d----限制大小: %d\n", finfo.Size(), l.MaxSize<<10<<10)

		l.fnum += 1
		fname := l.FileName + "." + strconv.Itoa(l.fnum)
		os.Rename(l.FileName, fname)
		fd, err := l.createFile()
		if err != nil {
			fmt.Printf("%s\n", err)
		}

		l.mw.fd = fd
	}
}

// 按天数分割
func (l *FileLog) splitbydaily() {
	l.lock.Lock()
	defer l.lock.Unlock()

	t, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))

	if l.date.Before(t) {
		l.mw.mu.Lock()
		defer l.mw.mu.Unlock()
		// 关闭之前的fd
		oldfd := l.mw.fd
		oldfd.Close()

		fname := l.FileName + "." + time.Now().Format(time.RFC3339)
		os.Rename(l.FileName, fname)
		fd, err := l.createFile()
		if err != nil {
			fmt.Printf("%s\n", err)
		}
		l.mw.fd = fd
	}

}

// 重命名文件
func (l *FileLog) rename() {

}

// 创建文件
func (l *FileLog) createFile() (*os.File, error) {

	dir := filepath.Dir(l.FileName)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	return os.OpenFile(l.FileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
}

func (l *FileLog) Write(msg string, level Level) error {
	if level < l.Level {
		return nil
	}

	l.docheck()

	l.log.Println(msg)

	return nil
}

func (l *FileLog) Destroy() {
	l.mw.fd.Close()
}

func (l *FileLog) Flush() {
	l.mw.fd.Sync()
}

func init() {
	Register("file", NewFile)
}
