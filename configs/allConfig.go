package configs

import (
	"github.com/dtm-labs/dtm/client/dtmcli"
	"github.com/spf13/viper"
	"log"
)

type AllConfig struct {
	Mode string `json:"mode"`

	//服务配置
	Server struct {
		ServerName string `json:"serverName"`
		HttpPort   int    `json:"httpPort"`
		GrpcPort   int    `json:"grpcPort"`
	}

	//日志配置
	Log struct {
		Level string `json:"level"`
	}

	//数据库配置
	Gorm struct {
		Url string `json:"url"`
	}

	//redis配置
	Redis struct {
		Single     bool     `json:"single"`
		SingleUrl  string   `json:"singleUrl"`
		ClusterUrl []string `json:"ClusterUrl"`
	}

	//etcd配置
	Etcd struct {
		Url []string `json:"url"`
	}
	Nsq struct {
		LookupdAddr string `json:"lookupdAddr"`
		NsqdAddr    string `json:"nsqdAddr"`
		Topic       string `json:"topic"`
	}

	//DtmConf struct {
	//	Driver   string `json:"Driver"`
	//	Host     string `json:"Host"`
	//	Port     int64  `json:"Port"`
	//	User     string `json:"User"`
	//	Password string `json:"Password"`
	//	Db       string `json:"Db"`
	//	Schema   string `json:"Schema"`
	//} `json:"dtmConf"`

	DtmConf dtmcli.DBConf
}

// 判断开发还是生产环境
func (receiver *AllConfig) IsDev() bool {
	if receiver.Mode == "dev" {
		return true
	} else if receiver.Mode == "prod" {
		return false
	} else if receiver.Mode == "" {
		return true
	}
	return true
}

// ReadConfig 读取所有配置
func ReadConfig(configPath string) *AllConfig {
	viper.SetConfigFile(configPath)

	//读取所有配置
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("配置文件格式不正确:", err)
	}

	var allConfig AllConfig
	err = viper.Unmarshal(&allConfig)
	if err != nil {
		log.Fatal("配置文件解析失败:", err)
	}
	return &allConfig
}
