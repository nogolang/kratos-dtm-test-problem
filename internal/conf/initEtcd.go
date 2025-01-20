package conf

import (
	"context"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/transport/grpc/resolver/discovery"
	"github.com/google/wire"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc/resolver"
	"kraots-xa/configs"
	// 导入 kratos 的 dtm 驱动
	_ "github.com/dtm-labs/driver-kratos"
	"time"
)

var EtcdProvider = wire.NewSet(NewKratosEtcdClient)

func NewKratosEtcdClient(logger *zap.Logger, allConfig *configs.AllConfig) *etcd.Registry {

	//指定所有的endpoints
	etcdConfig := clientv3.Config{
		Endpoints: allConfig.Etcd.Url,
	}

	//3.3x版本以后，超时不会直接通过error返回，必须要使用Status方法判断
	client, _ := clientv3.New(etcdConfig)
	timeout, _ := context.WithTimeout(context.Background(), 3*time.Second)
	_, err := client.Status(timeout, etcdConfig.Endpoints[0])
	if err != nil {
		logger.Error("连接etcd失败", zap.Error(err))
		return nil
	}

	//返回register对象，而非原始client
	r := etcd.New(client)

	//注册全局的resolver，现在我们的业务可以使用discovery:///dtmservice 来访问dtm和服务名称
	//  记得引入driver-kratos的驱动
	//我们在kratos里调用其它微服务，也是是通过discovery:///来的
	//  这一点kratos给我们实现了，然后dtm相当于是接入kratos的这种方式
	resolver.Register(
		discovery.NewBuilder(r, discovery.WithInsecure(true)))

	logger.Info("连接etcd成功")
	return r
}
