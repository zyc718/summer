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
	// 初始化日志
	if err = log.Init(config.Log); err != nil {
		panic(any("初始化日志失败" + err.Error()))
	}
	log.Infof("日志测试111%v", "aa")

	err = pool.Pool.Initiate(config.Store)
	if err != nil {
		panic(any("初始化连接池失败" + err.Error()))
	}

	cache.Pool, err = goredis.NewPool(config.Cache)
	if err != nil {
		panic(any("连接redis失败:" + err.Error()))
	}

	db, _ := cache.Pool.Get(nil)
	flag, err := db.Get("name")

	if flag == "" || err != nil {
		fmt.Printf("name 缓存不存在\n")
	} else {
		fmt.Printf("name 缓存是 %v\n", flag)

	}
	err = store.OpenAdapter(1, config.Store)
	if err != nil {
		panic(any("初始化数据库失败" + err.Error()))
	}
	//fmt.Printf("获取的缓存设置%v\n", string(config.Cache))
	//res, err := store.Topics.Get(*tags)
	//if err != nil {
	//	fmt.Printf("错误结果是%v\n", err.Error())
	//}
	//fmt.Printf("结果是%+v\n", res)
	//注册路由
	route()

	log.Infof("这是init 函数")

}
