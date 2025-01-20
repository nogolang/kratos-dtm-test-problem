package dtmUtils

import (
	"database/sql"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"kraots-xa/internal/utils/gormUtils"
)

// 从Dtm的的db里获取gorm,后面xa可能需要用到
// 这里的配置和gorm要一模一样
func GetGormDbFromDtmConn(db *sql.DB, logger *zap.Logger) *gorm.DB {
	//gorm适配zap
	myGormZap := gormUtils.NewMyGormZap(logger, gormlogger.Info)
	config := &gorm.Config{
		Logger: myGormZap,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		//不自动创建外键
		DisableForeignKeyConstraintWhenMigrating: true,
	}
	dialector := mysql.New(mysql.Config{Conn: db})
	gormDb, err := gorm.Open(dialector, config)
	if err != nil {
		logger.Error("从dtm里获取gorm对象失败", zap.Error(err))
		return nil
	}
	return gormDb
}
