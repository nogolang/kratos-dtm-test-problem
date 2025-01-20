package conf

import "github.com/google/wire"

// 除了http的，其他的都要提供
var ProviderSet = wire.NewSet(
	NewGrpcServer,
	NewHttpServer,
	ZapProvider,
	GormProvider,
	EtcdProvider,
	DtmConfProvider,
)
