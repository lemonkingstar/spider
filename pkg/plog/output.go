package plog

import (
	"io"

	"gopkg.in/natefinch/lumberjack.v2"
)

func SetOutput(output io.Writer) { logger.SetOutput(output) }

func SetRotateFile(logFile string) {
	rotateWriter := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    50,    // 达到50MB切割
		MaxBackups: 30,    // 最多保留30个备份文件
		MaxAge:     21,    // 21天前的日志自动删除
		LocalTime:  true,  // 本地时间作后缀
		Compress:   false, // 备份文件是否压缩
	}
	SetOutput(rotateWriter)
}
