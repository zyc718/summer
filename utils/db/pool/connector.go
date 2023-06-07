package pool

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type PoolStruct struct {
	config     map[string]*configType
	connectors map[string]*sqlx.DB
}

var Pool PoolStruct

type configType struct {
	Host     string `json:"hostname"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Database string `json:"database,omitempty"`
	Port     string `json:"port,omitempty"`
	Dbdriver string `json:"dbdriver,omitempty"`
	Charset  string `json:"char_set,omitempty"`
}

//json信息初始化
func (pool *PoolStruct) Initiate(jsonconfig json.RawMessage) error {
	var mapper map[string]interface{}

	err := json.Unmarshal(jsonconfig, &mapper)
	if err != nil {
		return err
	}

	if pool.config == nil {
		pool.config = make(map[string]*configType)
	}

	for k, v := range mapper["adapters"].(map[string]interface{}) {
		v := v.(map[string]interface{})

		pool.config[k] = &configType{
			Host:     v["hostname"].(string),
			Username: v["username"].(string),
			Password: v["password"].(string),
			Database: v["database"].(string),
			Port:     v["port"].(string),
			Dbdriver: v["dbdriver"].(string),
			Charset:  v["char_set"].(string),
		}
	}

	return nil
}

//初始化数据库连接
func (pool *PoolStruct) Open(name string) (*sqlx.DB, error) {
	var err error
	var config *configType

	config = pool.config[name]

	//如果存在已经建立的连接 直接返回
	if pool.connectors[name] != nil {
		return pool.connectors[name], nil
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local&collation=utf8mb4_unicode_ci", config.Username, config.Password, config.Host, config.Port, config.Database)

	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if pool.connectors == nil {
		pool.connectors = make(map[string]*sqlx.DB)
	}
	pool.connectors[name] = db

	return db, nil
}

//关闭连接
func (pool *PoolStruct) Close(name string) error {
	var err error
	if pool.connectors[name] != nil {
		//sqlx 关闭的方法
		err = pool.connectors[name].Close()
		delete(pool.connectors, name)
	}
	return err
}

//全部关闭

func (pool *PoolStruct) AllClose(name string) error {
	var err error
	for _, v := range pool.connectors {
		err = v.Close()
		if err != nil {
			break
		}
	}
	return err

}
