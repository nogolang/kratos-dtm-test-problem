package gormUtils

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
	"time"
)
import (
	gormlogger "gorm.io/gorm/logger"
)

type MyGormZap struct {
	ZapLogger *zap.Logger
	LogLevel  gormlogger.LogLevel
}

func NewMyGormZap(zapLog *zap.Logger, logLevel gormlogger.LogLevel) MyGormZap {
	//这里用值传递创建一个新的logger对象
	//然后设置一下logger
	logger := zapLog

	//跳出3层，这样才能知道是哪一个gorm语句调用的
	newLogger := logger.WithOptions(zap.AddCallerSkip(3))

	return MyGormZap{
		ZapLogger: newLogger,
		LogLevel:  logLevel,
	}
}

// 和原来的gormlooger一样使用，设置log的级别后返回新的log对象
func (receiver MyGormZap) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return MyGormZap{
		ZapLogger: receiver.ZapLogger,
		LogLevel:  level,
	}
}

func (receiver MyGormZap) Info(ctx context.Context, str string, args ...interface{}) {
	//判定现在指定的级别，当我们设置的级别是warn(3)，小于4，那么肯定打印不出来info信息
	if receiver.LogLevel < gormlogger.Info {
		return
	}
	//这里可以用info也可以用debug，无所谓的
	receiver.ZapLogger.Sugar().Info(str, args)
}
func (receiver MyGormZap) Warn(ctx context.Context, str string, args ...interface{}) {
	//当我们设置的级别是error(2)，小于3(warn)，那么肯定打印不出来warn的信息
	if receiver.LogLevel < gormlogger.Warn {
		return
	}
	receiver.ZapLogger.Sugar().Warn(str, args)
}

func (receiver MyGormZap) Error(ctx context.Context, str string, args ...interface{}) {
	if receiver.LogLevel < gormlogger.Error {
		return
	}
	receiver.ZapLogger.Sugar().Error(str, args)
}

// Trace Trace是gorm自动调用的，在执行完一条语句后会调用Trace
func (receiver MyGormZap) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	//花费时间
	latency := time.Since(begin)
	switch {
	//如果发现错误，则打印error
	case err != nil:
		sql, rowsAffected := fc()
		fields := []zapcore.Field{
			zap.String("sql", sql),
			zap.Int64("rowsAffected", rowsAffected),
			zap.Duration("latency", latency),
			zap.Error(err),
		}
		receiver.ZapLogger.Error("", fields...)
	case receiver.LogLevel >= gormlogger.Warn:
		//其他情况一律打印info即可，目前是这样，后续自己调整
		sql, rowsAffected := fc()

		//我们的sql，如果里面有换行符\r\n\t，那么存放json里的时候会被zap转换到对应的符号
		//所以这里替换掉
		sql = strings.Replace(sql, "\t", "", -1)
		sql = strings.Replace(sql, "\n", "", -1)
		sql = strings.Replace(sql, "\r", "", -1)

		fields := []zapcore.Field{
			zap.String("sql", sql),
			zap.Int64("rowsAffected", rowsAffected),
			zap.Duration("latency", latency),
		}
		receiver.ZapLogger.Info("", fields...)
	case receiver.LogLevel >= gormlogger.Info:
		sql, rowsAffected := fc()
		fields := []zapcore.Field{
			zap.String("sql", sql),
			zap.Int64("rowsAffected", rowsAffected),
			zap.Duration("latency", latency),
		}
		receiver.ZapLogger.Info("", fields...)
	}
}
