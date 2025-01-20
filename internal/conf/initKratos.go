package conf

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/transport/http"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"go.uber.org/zap"
	"kraots-xa/configs"
	controller "kraots-xa/internal/control"
	"kraots-xa/internal/utils/kratosMiddle"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type KratosServer struct {
	Logger           *zap.Logger
	GrpcServer       *grpc.Server
	HttpServer       *http.Server
	KratosEtcdClient *etcd.Registry
	AllConfig        *configs.AllConfig

	//注入各种controller，我们已经创建好了grpc+中间件的grpc对象
	//然后创建GrpcServer的时候才会注册各种grpc方法
	UserAccountControl *controller.UserAccountControl
}

func NewGrpcServer(logger *zap.Logger, allConfig *configs.AllConfig) *grpc.Server {
	grpcServer := grpc.NewServer(
		grpc.Address(":"+fmt.Sprintf("%d", allConfig.Server.GrpcPort)),
		grpc.Middleware(
			//请求log中间件
			kratosMiddle.KratosZapMiddle(logger),

			//recovery中间件
			//会自动recovery，然后把error传递进来，这个和gin的很像
			recovery.Recovery(
				recovery.WithHandler(func(ctx context.Context, req, err interface{}) error {
					//出现系统的panic，那么就是打印出来
					logger.Sugar().Errorf("系统出现panic: %+v", err)
					return nil
				}),
			),
		),
	)
	return grpcServer
}

func NewHttpServer(logger *zap.Logger, allConfig *configs.AllConfig) *http.Server {
	httpServer := http.NewServer(
		http.Address(":"+fmt.Sprintf("%d", allConfig.Server.HttpPort)),
		http.Middleware(
			//请求log中间件
			kratosMiddle.KratosZapMiddle(logger),

			//recovery中间件
			//会自动recovery，然后把error传递进来，这个和gin的很像
			recovery.Recovery(
				recovery.WithHandler(func(ctx context.Context, req, err interface{}) error {
					//出现系统的panic，那么就是打印出来
					logger.Sugar().Errorf("系统出现panic: %+v", err)
					return nil
				}),
			),
		),
	)
	return httpServer
}

func (receiver *KratosServer) RunServer() {
	//创建kratos app
	app := kratos.New(
		kratos.Name(receiver.AllConfig.Server.ServerName),
		kratos.Registrar(receiver.KratosEtcdClient),
		kratos.Server(receiver.GrpcServer, receiver.HttpServer))

	go func() {
		err := app.Run()
		if err != nil {
			receiver.Logger.Error("kratos服务启动失败：", zap.Error(err))
			return
		}
	}()

	//平滑重启
	//创建信号，返回一个channel
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	//程序关闭了，则协程有值，执行到这里
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	//无需重复关闭，在app.run里有对应的监控信号会自动关闭server
	//if err := app.Stop(); err != nil {
	//  log.Error("关闭服务错误", zap.Error(err))
	//}

	receiver.Logger.Info(receiver.AllConfig.Server.ServerName + "  释放资源中....")
	receiver.Logger.Info(receiver.AllConfig.Server.ServerName + "  退出了....")
	<-ctx.Done()
}
