package log

import (
	"fmt"
	"runtime"
	"time"
)

// Bean ...
type Bean struct {
	Level    string
	Title    string
	FuncInfo string
	Body     string
}

// NewBean ...
func NewBean(title string) *Bean {
	b := &Bean{
		Title: title,
	}
	return b
}

// SetTitle ...
func (b *Bean) SetTitle(title string) {
	b.Title = title
}

// SetFuncInfo ...
func (b *Bean) SetFuncInfo(info string) {
	b.FuncInfo = info
}

// Add ...
func (b *Bean) Add(k string, v interface{}) {
	value := fmt.Sprint(v)
	b.Body += "||" + k + "=" + value
}

// WritePublicLog ...
func (b *Bean) WritePublicLog() {
	writeLog(pl, b.Title+b.Body)
}

func (b *Bean) write() {

	if b.Level == "" {
		b.Level = INFO
	}

	if b.Title == "" {
		b.Title = "_undef"
	}

	logtime := time.Now().Format(ktnLogTimeFormat)
	log := fmt.Sprintf("[%s][%s][%s] %s%s", b.Level, logtime, b.FuncInfo, b.Title, b.Body)

	writeLog(l, log)
}

// Debug ...
func (b *Bean) Debug() {
	b.Level = DEBUG
	b.write()
}

// Info ...
func (b *Bean) Info() {
	b.Level = INFO
	b.write()
}

// Notice ...
func (b *Bean) Notice() {
	b.Level = NOTICE
	b.write()
}

// Warning ...
func (b *Bean) Warning() {
	b.Level = WARNING
	b.write()
}

// Error ...
func (b *Bean) Error() {
	b.Level = ERROR
	b.write()
}

// Fatal ...
func (b *Bean) Fatal() {
	b.Level = FATAL
	b.write()
}

// GetFuncInfo ...
func GetFuncInfo(pc uintptr) string {
	file, line := runtime.FuncForPC(pc).FileLine(pc)
	return fmt.Sprintf("%s:%d", suffix(file, 2), line)
}

func suffix(s string, n int) string {
	if n <= 0 {
		return s
	}

	i := len(s) - 1
	for ; i > 0; i-- {
		if s[i] == '/' {
			n--
			if n == 0 {
				break
			}
		}
	}
	return s[i+1:]
}
