package conf

import (
	"github.com/google/wire"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"kraots-xa/configs"
	"log"
	"os"
	"time"
)

var ZapProvider = wire.NewSet(NewZapConfig)

func NewZapConfig(allConfig *configs.AllConfig) *zap.Logger {

	//日志级别
	level := zap.NewAtomicLevel()
	switch allConfig.Log.Level {
	case "debug":
		level.SetLevel(zap.DebugLevel)
	case "info":
		level.SetLevel(zap.InfoLevel)
	case "warning":
		level.SetLevel(zap.WarnLevel)
	case "error":
		level.SetLevel(zap.ErrorLevel)
	case "fatal":
		level.SetLevel(zap.FatalLevel)
	default:
		log.Fatal("日志级别指定错误")
		return nil
	}

	var logger *zap.Logger

	//设置我们指定的级别
	newCoreError := zapcore.NewCore(getEncoding(allConfig), getLogWriter(allConfig), level.Level())
	logger = zap.New(newCoreError, zap.AddCaller())

	//这里使用了wire，严格准守di原则
	//但是有些地方可能不太方便传递logger对象，比如中间件的地方使用全局的也可以
	zap.ReplaceGlobals(logger)
	return logger
}

func getEncoding(allConfig *configs.AllConfig) zapcore.Encoder {
	var newEncoder zapcore.Encoder

	encodeTime := func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(t.Format(time.DateTime))
	}

	log.Println("日志级别为：", allConfig.Log.Level) //设置编码方式和自定义的时间，开发环境就是json，生产环境是Console
	if isDev := allConfig.IsDev(); isDev {
		config := zap.NewDevelopmentEncoderConfig()
		config.EncodeTime = encodeTime
		newEncoder = zapcore.NewConsoleEncoder(config)
	} else {
		config := zap.NewProductionEncoderConfig()
		config.EncodeTime = encodeTime
		newEncoder = zapcore.NewJSONEncoder(config)
	}
	return newEncoder
}

func getLogWriter(allConfig *configs.AllConfig) zapcore.WriteSyncer {
	//如果是开发环境，向控制台输出，生产环境应该向文件输出
	//文件则向lumber输出，由lumber切割
	var writer io.Writer
	if isDev := allConfig.IsDev(); isDev {
		writer = os.Stdout
	} else {
		writer = lumberJackConfig()
	}
	return zapcore.AddSync(writer)
}

// 日志切割
func lumberJackConfig() *lumberjack.Logger {
	//获取项目目录，如果本目录下logs目录不存在
	//就在当前项目运行目录下创建logs目录
	dir, _ := os.Getwd()
	dir = dir + "/logs"

	//判断有没有logs目录
	_, err := os.ReadDir(dir)
	if err != nil {
		//目录不存在，则创建
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			log.Fatal("logs目录创建失败，请手动创建")
			return nil
		}
	}

	//普通日志输出到info中，错误日志输出到error中
	//外面会对zapcore进行合并，输出2个文件，但是你里面的文件名称要正确
	//fileName := ""
	//switch level {
	//case "info":
	//	fileName = dir + "/app.info.log"
	//case "error":
	//	fileName = dir + "/app.error.log"
	//}

	//统一输出到app.log里，暂时不区分error和info
	//到时候再统一处理日志
	fileName := dir + "/app.log"

	return &lumberjack.Logger{
		Filename: fileName,

		//日志文件的最大尺寸,单位MB
		//切割出来的每个文件都是xMB,但是最开始的主文件可能会小一点
		MaxSize: 10,

		//保留的旧的最大个数.此时我们输出了10MB的内容.
		//但是只有5个切割文件+1个主文件.其余5个都删掉了.按照切割出来的日期.早期的会优先进行删除
		//如果旧日志一直没有删除(没有满5个).但是已经过去30天了.这时候会自动删除
		MaxBackups: 5,

		//保留旧文件的最大天数
		MaxAge: 30,

		//是否压缩旧文件
		Compress: false,
	}
}
