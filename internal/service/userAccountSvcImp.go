package service

import (
	"database/sql"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"kraots-xa/internal/dao"
	"kraots-xa/proto/pb/userAccountPb"
)

type UserAccountSvcImp struct {
	UserAccountDao *dao.UserAccountDaoImp
	Logger         *zap.Logger
}

/*
这里的db需要用dtm传递过来的db才行
*/
func (receiver UserAccountSvcImp) UpdateAccount(db *gorm.DB, request *userAccountPb.UserAccountUpdateRequest) error {
	return receiver.UserAccountDao.UpdateAccount(db, request)
}
