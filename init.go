package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"summer/utils/log"
)

type configType struct {
	Log json.RawMessage `json:"log_config"`
}

func init() {
	//读取配置文件
	var configfile = flag.String("config", "config.conf", "Path to config file.")
	flag.Parse()
	var config configType
	fileContent, err := ioutil.ReadFile(*configfile)
	if err != nil {
		log.Fatalf("Failed to read config file: %s", err.Error())
	}
	if err = json.Unmarshal(fileContent, &config); err != nil {
		log.Fatalf("Failed to decode config file: %s", err.Error())
	}
	//初始化日志
	if err = log.Init(config.Log); err != nil {
		panic(any("初始化日志失败" + err.Error()))
	}
	log.Infof("这是日志debug测试 %v\n", "zzzz")
	route()
}
