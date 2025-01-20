package conf

import (
	"github.com/dtm-labs/dtm/client/dtmcli"
	"github.com/google/wire"
	"go.uber.org/zap"
	"kraots-xa/configs"
)

var DtmConfProvider = wire.NewSet(NewDtmConf)

func NewDtmConf(logger *zap.Logger, allConfig *configs.AllConfig) *dtmcli.DBConf {
	return &dtmcli.DBConf{
		Driver:   allConfig.DtmConf.Driver,
		Host:     allConfig.DtmConf.Host,
		Port:     allConfig.DtmConf.Port,
		User:     allConfig.DtmConf.User,
		Password: allConfig.DtmConf.Password,
		Db:       allConfig.DtmConf.Db,
	}
}
