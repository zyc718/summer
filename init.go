package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"summer/utils/cache"
	"summer/utils/cache/goredis"
	_ "summer/utils/db/mysql"
	"summer/utils/db/pool"
	"summer/utils/log"
	"summer/utils/store"
)

type configType struct {
	Log   json.RawMessage `json:"log_config"`
	Store json.RawMessage `json:"store_config"`
	Cache json.RawMessage `json:"cache_config"`
}

func init() {
	// 读取配置文件
	var configfile = flag.String("config", "config.conf", "Path to config file.")

	//var tags = flag.String("tags", "mysql", "")

	flag.Parse()
	var config configType
	fileContent, err := ioutil.ReadFile(*configfile)
	if err != nil {
		log.Fatalf("Failed to read config file: %s", err.Error())
	}
	if err = json.Unmarshal(fileContent, &config); err != nil {
		log.Fatalf("Failed to decode config file: %s", err.Error())
	}
	fmt.Printf("config success\n")

	// 初始化日志
	if err = log.Init(config.Log); err != nil {
		panic(any("初始化日志失败" + err.Error()))
	}
	fmt.Printf("login success\n")

	err = pool.Pool.Initiate(config.Store)
	if err != nil {
		panic(any("初始化连接池失败" + err.Error()))
	}

	fmt.Printf("pool success\n")

	cache.Pool, err = goredis.NewPool(config.Cache)
	if err != nil {
		panic(any("连接redis失败:" + err.Error()))
	}

	fmt.Printf("redis success\n")

	//db, _ := cache.Pool.Get(nil)
	//flag, err := db.Get("name")
	//
	//if flag == "" || err != nil {
	//	fmt.Printf("name 缓存不存在\n")
	//} else {
	//
	//}
	err = store.OpenAdapter(1, config.Store)
	if err != nil {
		panic(any("初始化数据库失败" + err.Error()))
	}
	fmt.Printf("db success\n")

	//注册路由
	route()
	log.Infof("init done")

}
