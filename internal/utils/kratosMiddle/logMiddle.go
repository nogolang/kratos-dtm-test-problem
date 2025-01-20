package kratosMiddle

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

func KratosZapMiddle(outLog *zap.Logger) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			//无需打印caller
			logger := outLog.WithOptions(zap.WithCaller(false))
			var (
				code          int32
				reason        string
				message       string
				metadata      map[string]string
				kind          string //操作种类，grpc还是http
				method        string //操作的方法
				methodAndKind string
			)
			startTime := time.Now()

			//从ctx里可以拿到一些基本信息
			if info, ok := transport.FromServerContext(ctx); ok {
				kind = info.Kind().String()
				method = info.Operation()
				methodAndKind = fmt.Sprintf("%s%s", kind, method)
			}

			//如果返回错误，就打印error，否则打印info
			//并且加入一些自定义的字段
			fields := []zapcore.Field{
				//参数，这里使用sprintf直接打印出request对象了，里面就是参数
				//这个req只包含我们自己的内容，没有proto那些杂七杂八的字段
				zap.String("args", fmt.Sprintf("%+v", req)),
				zap.Duration("latency", time.Since(startTime)),
			}

			//调用下一个中间件，然后拿到最终响应，解析出错误
			//这里的errors是kratos的，解析出来的错误也是kratos的
			reply, err = handler(ctx, req)
			if err != nil {
				//判定取到的是不是kratos的错误，也就是我们自定义的错误
				//那我们就打印的是info，而非error，因为是可预期的
				if kratosErr := errors.FromError(err); kratosErr != nil {
					code = kratosErr.GetCode()
					reason = kratosErr.GetReason()
					message = kratosErr.GetMessage()
					metadata = kratosErr.GetMetadata()

					fields = append(fields,
						zap.Int32("code", code),
						zap.String("message", message),
						zap.String("reason", reason),
						zap.Any("metadata", metadata),
					)

					logger.Info(methodAndKind, fields...)
				} else {
					//如果是其他错误，没有经过我们处理过的
					//那就打印出来错误，并且打印堆栈，所以我们的这种第三方错误最用warp包裹出来
					logger.Error(methodAndKind, fields...)
					logger.Sugar().Errorf("%+v", err)
				}
				return reply, err
			}

			//如果没有error，那么打印 方法，参数，延迟即可
			logger.Info(methodAndKind, fields...)

			//返回值，下一个中间件调用的时候拿到的就是这个值
			return reply, err
		}
	}
}
