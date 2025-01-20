package test

import (
	"context"
	"github.com/go-kratos/kratos/v2/selector"
	"github.com/go-kratos/kratos/v2/selector/filter"
	"github.com/go-kratos/kratos/v2/selector/random"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	rawGrpc "google.golang.org/grpc"
	"kraots-xa/configs"
	"kraots-xa/internal/conf"
	"log"
)

// 解析出所有配置

func GetGrpcClient() *rawGrpc.ClientConn {
	flagPath := "../configs/config.dev.yaml"
	//初始化配置
	AllConfig := configs.ReadConfig(flagPath)

	logger := conf.NewZapConfig(AllConfig)
	KratosEtcdClient := conf.NewKratosEtcdClient(logger, AllConfig)

	//创建全局的负载均衡算法为random
	//还有p2c,wrr，具体看官方和资料
	//由于 gRPC 框架的限制，只能使用全局 balancer name 的方式来注入 selector
	selector.SetGlobalSelector(random.NewBuilder())

	//创建路由 Filter：筛选版本号为"xxx"的实例
	//这里我注册的时候就填的空,所以也为空
	filterVersion := filter.Version("")

	grpcClient, err := grpc.DialInsecure(
		context.Background(),

		//不打印服务发现的日志，不然invoke之后会输出很多etcd相关的东西
		grpc.WithPrintDiscoveryDebugLog(false),

		//服务发现语法
		//<schema>://[authority]/<service-name>
		grpc.WithEndpoint("discovery:///"+AllConfig.Server.ServerName),
		grpc.WithDiscovery(KratosEtcdClient),
		grpc.WithNodeFilter(filterVersion),
	)
	if err != nil {
		log.Fatal("服务发现错误", err)
		return nil
	}

	return grpcClient
}
