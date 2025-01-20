package main

import (
	"flag"
	"kraots-xa/configs"
)

var flagPath string

func init() {
	//默认加载dev的路径
	//测试是时候go run .\cmd\server\main.go
	//如果是build之后
	flag.StringVar(&flagPath, "conf", "configs/config.dev.yaml", "config path, eg: -conf configs/config.dev.yaml")
}

func main() {
	flag.Parse()

	//初始化配置
	allConfig := configs.ReadConfig(flagPath)

	app := WireApp(allConfig)
	app.RunServer()
}
