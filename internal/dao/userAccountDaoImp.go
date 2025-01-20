package dao

import (
	"github.com/dtm-labs/dtm/client/dtmcli"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"kraots-xa/proto/pb/userAccountPb"
	"strings"
)

type UserAccountDaoImp struct {
	Logger *zap.Logger
}

func (receiver *UserAccountDaoImp) UpdateAccount(db *gorm.DB,
	request *userAccountPb.UserAccountUpdateRequest) error {
	if strings.Contains(request.TransInResult, dtmcli.ResultFailure) ||
		strings.Contains(request.TransOutResult, dtmcli.ResultFailure) {
		//这里返回ErrFailure，就代表事务需要回滚
		return dtmcli.ErrFailure
	}

	sql := `
		update dtm_test.user_account set balance = balance + ? where user_id = ?
	`

	return db.Exec(sql, request.Amount, request.Uid).Error
}
