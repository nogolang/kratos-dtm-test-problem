package controller

import (
	"github.com/google/wire"
	"kraots-xa/internal/dao"
	"kraots-xa/internal/service"
)

var ControlProvider = wire.NewSet(
	UserAccountControlProvider,
)

var UserAccountControlProvider = wire.NewSet(
	NewUserAccountControl,
	wire.Struct(new(service.UserAccountSvcImp), "*"),
	wire.Struct(new(dao.UserAccountDaoImp), "*"),
)
