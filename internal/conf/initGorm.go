package conf

import (
	"github.com/google/wire"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"kraots-xa/configs"
	"kraots-xa/internal/utils/gormUtils"
)

var GormProvider = wire.NewSet(NewGormConfig)

// NewGormConfig logger由外部注入进来
func NewGormConfig(logger *zap.Logger, allConfig *configs.AllConfig) *gorm.DB {

	//gorm适配zap
	myGormZap := gormUtils.NewMyGormZap(logger, gormlogger.Info)

	//初始化gorm的配置
	config := &gorm.Config{
		Logger: myGormZap,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		//不自动创建外键
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	//gormDb无需使用.session，它Open出来就是一个链式安全的实例
	db, err := gorm.Open(mysql.Open(allConfig.Gorm.Url), config)
	if err != nil {
		logger.Fatal("gorm连接数据库失败", zap.Error(err))
		return nil
	}

	logger.Info("连接mysql成功")

	//迁移一些表
	//db.AutoMigrate(&model.User{})
	return db
}
