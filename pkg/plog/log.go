package plog

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	// 参考 log4j
	// %L 输出代码中的行号
	// %l 输出日志事件的发生位置，包括类目名、发生的线程，以及在代码中的行数 如：Testlog.main(TestLog.java:10)
	// %m 输出代码中指定的消息
	// %p 输出优先级，即DEBUG,INFO,WARN,ERROR,FATAL
	// %c 输出所属的类目,通常就是所在类的全名
	// %t 输出产生该日志事件的线程名
	// %n 输出一个回车换行符，Windows平台为“\r\n”，Unix平台为“\n”
	// %d 输出日志时间点的日期或时间，默认格式为ISO8601，也可以在其后指定格式 如：%d{yyyy年MM月dd日 HH:mm:ss,SSS}，输出类似：2012年01月05日 22:10:28,921
	defaultLogFormat       = "[ %p ][ %d ] %m"
	fileLineLogFormat      = "[ %p ][ %d ][ %l ] %m"
	defaultTimestampFormat = time.RFC3339
	// 混淆码
	obfuscatedCode = "*#06#"
)

var (
	logger = logrus.New()

	DebugLevel = logrus.DebugLevel
	InfoLevel  = logrus.InfoLevel
	WarnLevel  = logrus.WarnLevel
	ErrorLevel = logrus.ErrorLevel

	Print  = logger.Print
	Printf = logger.Printf
	Debug  = logger.Debug
	Debugf = logger.Debugf
	Info   = logger.Info
	Infof  = logger.Infof
	Warn   = logger.Warn
	Warnf  = logger.Warnf
	Error  = logger.Error
	Errorf = logger.Errorf
	Fatal  = logger.Fatal
	Fatalf = logger.Fatalf
	Panic  = logger.Panic
	Panicf = logger.Panicf
)

var (
	defaultFormatter = &stdFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		LevelTruncation: true,
	}
)

type Logger = logrus.Logger
type Entry = logrus.Entry

func init() {
	defaultFormatter.build()
	logger.SetFormatter(defaultFormatter)
	logger.SetLevel(InfoLevel)
	logger.SetOutput(os.Stdout)
}

func GetLogger() *Logger                              { return logger }
func WithField(key string, value interface{}) *Entry  { return logger.WithField(key, value) }
func WithFields(fields map[string]interface{}) *Entry { return logger.WithFields(fields) }
func SetLevel(l logrus.Level)                         { logger.SetLevel(l) }
func SetOutput(output io.Writer)                      { logger.SetOutput(output) }

func SetReportCaller(reportCaller bool) {
	logger.SetReportCaller(reportCaller)
	defaultFormatter.LogFormat = fileLineLogFormat
	defaultFormatter.build()
}

// SetRotateFile enable circular log files. e.g. logs/*.log.
func SetRotateFile(file string) {
	rotate := &rotateFileHook{
		level:     logger.GetLevel(),
		formatter: logger.Formatter,
		logWriter: &lumberjack.Logger{
			Filename:   file,
			MaxSize:    50,
			MaxBackups: 30,
			MaxAge:     7,
			LocalTime:  true,
		},
	}
	logger.AddHook(rotate)
}

type stdFormatter struct {
	TimestampFormat string
	LogFormat       string
	LevelTruncation bool

	formatContent string
	fLevelCode    string
	fDateCode     string
	fCallerCode   string
	fMessageCode  string
}

func (f *stdFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	output := f.formatContent
	level := strings.ToUpper(entry.Level.String())
	if f.LevelTruncation {
		level = level[:4]
	}
	// append log level
	output = strings.ReplaceAll(output, f.fLevelCode, level)
	// append log time
	output = strings.ReplaceAll(output, f.fDateCode, entry.Time.Format(f.TimestampFormat))
	if entry.HasCaller() {
		// append file/line info
		callerMessage := fmt.Sprintf("%s/%s:%d", entry.Caller.Function,
			filepath.Base(entry.Caller.File), entry.Caller.Line)
		output = strings.Replace(output, f.fCallerCode, callerMessage, 1)
	}
	if len(entry.Data) > 0 {
		fieldMessage := make([]string, 0, len(entry.Data))
		for k, v := range entry.Data {
			fieldMessage = append(fieldMessage, fmt.Sprintf("%s: %v", k, v))
		}
		// append fields info before message
		output = strings.Replace(output, f.fMessageCode,
			fmt.Sprintf("[ %s ] %s", strings.Join(fieldMessage, ","), f.fMessageCode), 1)
	}

	// append message
	output = strings.Replace(output, f.fMessageCode, entry.Message, 1)
	output += "\n"
	return []byte(output), nil
}

func (f *stdFormatter) build() {
	if f.LogFormat == "" {
		f.LogFormat = defaultLogFormat
	}
	if f.TimestampFormat == "" {
		f.TimestampFormat = defaultTimestampFormat
	}
	f.fLevelCode = "%p" + obfuscatedCode + "%"
	f.fDateCode = "%d" + obfuscatedCode + "%"
	f.fCallerCode = "%l" + obfuscatedCode + "%"
	f.fMessageCode = "%m" + obfuscatedCode + "%"
	f.formatContent = strings.ReplaceAll(f.LogFormat, "%p", f.fLevelCode)
	f.formatContent = strings.ReplaceAll(f.formatContent, "%d", f.fDateCode)
	f.formatContent = strings.ReplaceAll(f.formatContent, "%l", f.fCallerCode)
	f.formatContent = strings.ReplaceAll(f.formatContent, "%m", f.fMessageCode)
}

type rotateFileHook struct {
	logWriter io.Writer
	formatter logrus.Formatter
	level     logrus.Level
}

func (h *rotateFileHook) Levels() []logrus.Level {
	return logrus.AllLevels[:h.level+1]
}

func (h *rotateFileHook) Fire(entry *logrus.Entry) (err error) {
	b, err := h.formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = h.logWriter.Write(b)
	return
}
