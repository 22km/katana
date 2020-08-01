package log

import (
	"fmt"
	"io"
	"os"

	"code.cloudfoundry.org/go-diodes"
	"github.com/arthurkiller/rollingwriter"
	"github.com/gin-gonic/gin"
)

// 目前支持的级别，线上只支持INFO及以上级别
const (
	DEBUG   = "DEBUG"
	INFO    = "INFO"
	NOTICE  = "NOTICE"
	WARNING = "WARNING"
	ERROR   = "ERROR"
	FATAL   = "FATAL"
)

const (
	ktnLogTimeFormat = "2006-01-02T15:04:05.999Z0700"
)

var (
	l  *logger // info&warn 日志
	pl *logger // public 日志
)

// Config 日志配置
type Config struct {
	LogDir        string `yaml:"logDir"`
	LogName       string `yaml:"logName"`
	PublicLogDir  string `yaml:"publicLogDir"`
	PublicLogName string `yaml:"publicLogName"`
}

type logger struct {
	running bool
	writer  io.Writer
	i       *diodes.ManyToOne
	o       *diodes.Poller
}

// Init 初始化日志
func Init(c *Config) error {
	var err error

	if l, err = newLogger(c.LogDir, c.LogName); err != nil {
		return err
	}
	gin.DefaultErrorWriter = l.writer
	gin.DefaultWriter = l.writer

	if pl, err = newLogger(c.PublicLogDir, c.PublicLogName); err != nil {
		return err
	}

	return nil
}

// Close TODO: 清空缓存里的日志
func Close() {}

func newLogger(path, fileName string) (*logger, error) {
	if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
		os.Mkdir(path, 0777)
	}

	d := diodes.NewManyToOne(4096, nil)
	logger := &logger{
		running: true,
		i:       d,
		o:       diodes.NewPoller(d),
	}

	// writer 实现了 io.Writer 的全部接口
	// 使用配置方式生成一个 writer 或者 Option 都可以
	c := rollingwriter.Config{
		LogPath:       path,         // 日志路径
		TimeTagFormat: "2006010215", // 时间格式串
		FileName:      fileName,     // 日志文件名
		MaxRemain:     0,            // 配置日志最大存留数, 0代表不限制 (使用odin来限制日志保留时间)

		// 目前有2种滚动策略: 按照时间滚动按照大小滚动
		// - 时间滚动: 配置策略如同 crontable, 例如,每天0:0切分, 则配置 0 0 0 * * *
		// - 大小滚动: 配置单个日志文件(未压缩)的滚动大小门限, 如1G, 500M
		RollingPolicy:      rollingwriter.TimeRolling, // 配置滚动策略 norolling timerolling volumerolling
		RollingTimePattern: "0 0 * * * *",             // 配置时间滚动策略, 每小时滚动
		RollingVolumeSize:  "",                        // 配置截断文件下限大小
		Compress:           false,                     // 配置是否压缩存储

		// writer 支持3种方式:
		// - 无保护的 writer: 不提供并发安全保障
		// - lock 保护的 writer: 提供由 mutex 保护的并发安全保障
		// - 异步 writer: 异步 write, 并发安全. 异步开启后忽略 Lock 选项
		Asynchronous: false, // 配置是否异步写
		Lock:         true,  // 配置是否同步加锁写
	}

	// 创建一个 writer
	writer, err := rollingwriter.NewWriterFromConfig(&c)
	if err != nil {
		return nil, err
	}

	logger.writer = writer
	// gin.DefaultErrorWriter = writer
	// gin.DefaultWriter = writer

	go func() {
		for logger.running {
			i := logger.o.Next()
			fmt.Fprintln(logger.writer, *(*string)(i))
		}
	}()

	return logger, nil
}

func writeLog(logger *logger, format string, a ...interface{}) {
	if logger == nil {
		return
	}

	msg := ""
	if a == nil {
		// 防止format中有%字符影响打印结果
		msg = format
	} else {
		msg = fmt.Sprintf(format, a...)
	}
	logger.i.Set(diodes.GenericDataType(&msg))
}

// WriteLog ...
func WriteLog(format string, a ...interface{}) {
	writeLog(l, format, a...)
}

// WritePublicLog ...
func WritePublicLog(format string, a ...interface{}) {
	writeLog(pl, format, a...)
}
