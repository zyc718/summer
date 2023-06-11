package goredis

import (
	"context"
	"encoding/json"
	"errors"
	rds "github.com/go-redis/redis/v8"
	"log"
	"strings"
	"summer/utils/cache"
	"sync"
	"time"
)

const (
	defaultAddr       = "localhost:6379"
	defaultPassword   = ""
	defaultDB         = 0
	defaultMaxRetries = 3
)

var prefix string

type configType struct {
	Addr       string `json:"addr,omitempty"`
	Password   string `json:"password,omitempty"`
	DB         int    `json:"db,omitempty"`
	Prefix     string `json:"prefix,omitempty"`
	MaxRetries int    `json:"MaxRetries,omitempty"` //放弃前最大重试次数
}
type pool struct {
	client *rds.Client
	once   sync.Once
}

type CacheHandler struct {
	client *rds.Client
	ctx    context.Context
}

func NewPool(jsonconfig json.RawMessage) (cache.CachePool, error) {
	var p = &pool{}
	var err error
	p.once.Do(func() {
		var config configType
		if err = json.Unmarshal(jsonconfig, &config); err != nil {
			return
		}
		if config.Addr == "" {
			config.Addr = defaultAddr
		}

		if config.Password == "" {
			config.Password = defaultPassword
		}
		p.client = rds.NewClient(&rds.Options{
			Addr:       config.Addr,
			Password:   config.Password, // no password set
			DB:         config.DB,       // use default DB
			MaxRetries: config.MaxRetries,
		})
		prefix = config.Prefix

		if !p.IsOpen() {
			err = errors.New("redis连接失败")
			return
		}
	})
	return p, err

}

func (p *pool) Get(ctx context.Context) (cache.CacheHandler, error) {
	if ctx == nil && p != nil && p.client != nil {
		ctx = p.client.Context()
	}

	return &CacheHandler{client: p.client, ctx: ctx}, nil
}

// 关闭消息队列
func (p *pool) Close() error {
	var err error
	if p.client != nil {
		err = p.client.Close()
		p.client = nil
	}
	return err
}

// 消息队列是否连接
func (p *pool) IsOpen() bool {
	var pong string
	if p.client != nil {
		pong, _ = p.client.Ping(p.client.Context()).Result()
	}
	return pong == "PONG"
}

func keyGen(key string) string {
	if prefix != "" {
		key = prefix + key
	}
	return key
}

// 获取缓存值
func (cachehandler *CacheHandler) Get(key string) (string, error) {
	key = keyGen(key)

	value, err := cachehandler.client.Get(cachehandler.ctx, key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

//设置缓存值
func (cachehandler *CacheHandler) Set(key string, value interface{}, expiration time.Duration) error {
	key = keyGen(key)
	err := cachehandler.client.Set(cachehandler.ctx, key, value, expiration).Err()
	if err != nil {
		return errors.New("SET失败，错误信息: " + err.Error())
	}
	return nil
}

// 设置缓存值 setNx
func (cacheHandler *CacheHandler) SetNx(key string, value interface{}, expiration time.Duration) (bool, error) {
	key = keyGen(key)
	flag, err := cacheHandler.client.SetNX(cacheHandler.ctx, key, value, expiration).Result()
	return flag, err

}

// 删除缓存值
func (cacheHandler *CacheHandler) Del(key string) error {
	key = keyGen(key)
	err := cacheHandler.client.Del(cacheHandler.ctx, key).Err()
	if err != nil {
		return errors.New("Del失败，错误信息: " + err.Error())
	}
	return nil
}

// 设置哈希消息
func (cacheHandler *CacheHandler) HSet(channel, key string, value interface{}) error {
	channel = keyGen(channel)
	err := cacheHandler.client.HSet(cacheHandler.ctx, channel, key, value).Err()
	if err != nil {
		return errors.New("HSET失败，错误信息: " + err.Error())
	}
	return nil
}

// 获取哈希消息
func (cacheHandler *CacheHandler) HGet(channel, key string) (string, error) {
	channel = keyGen(channel)
	value, err := cacheHandler.client.HGet(cacheHandler.ctx, channel, key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

// 获取哈希消息
func (cacheHandler *CacheHandler) HGetAll(channel string) (map[string]string, error) {
	channel = keyGen(channel)
	value, err := cacheHandler.client.HGetAll(cacheHandler.ctx, channel).Result()
	if err != nil {
		return nil, errors.New("HGETALL失败，错误信息: " + err.Error())
	}
	return value, nil
}

// 删除哈希消息
func (cacheHandler *CacheHandler) HDel(channel, key string) error {
	channel = keyGen(channel)
	err := cacheHandler.client.HDel(cacheHandler.ctx, channel, key).Err()
	if err != nil {
		return errors.New("HDEL失败，错误信息: " + err.Error())
	}
	return nil
}

// 序列添加or更改元素
func (cacheHandler *CacheHandler) ZAdd(channel, key string, score float64) error {
	channel = keyGen(channel)
	err := cacheHandler.client.ZAdd(cacheHandler.ctx, channel, &rds.Z{Score: score, Member: key}).Err()
	if err != nil {
		return errors.New("ZADD失败，错误信息: " + err.Error())
	}
	return nil
}

// 获取集合长度
func (cacheHandler *CacheHandler) ZCard(channel string) (int64, error) {
	channel = keyGen(channel)
	value, err := cacheHandler.client.ZCard(cacheHandler.ctx, channel).Result()
	if err != nil {
		return 0, errors.New("ZCard失败，错误信息: " + err.Error())
	}
	return value, nil
}

// 获取序列值
func (cacheHandler *CacheHandler) ZRange(channel string, start, stop int64) ([]string, error) {
	channel = keyGen(channel)
	stringArray, err := cacheHandler.client.ZRange(cacheHandler.ctx, channel, start, stop).Result()
	if err != nil {
		return nil, errors.New("ZRANGE失败，错误信息: " + err.Error())
	}
	return stringArray, nil
}

// 获取序列值
func (cacheHandler *CacheHandler) ZRangeWithScores(channel string, start, stop int64) ([]rds.Z, error) {
	channel = keyGen(channel)
	zStruct, err := cacheHandler.client.ZRangeWithScores(cacheHandler.ctx, channel, start, stop).Result()
	if err != nil {
		return nil, errors.New("ZRANGE失败，错误信息: " + err.Error())
	}
	return zStruct, nil
}

// 序列删除元素
func (cacheHandler *CacheHandler) ZRem(channel, key string) error {
	channel = keyGen(channel)
	err := cacheHandler.client.ZRem(cacheHandler.ctx, channel, key).Err()
	if err != nil {
		return errors.New("ZREM失败，错误信息: " + err.Error())
	}
	return nil
}

//执行lua 脚本
func (cacheHandler *CacheHandler) Eval(script *cache.Script, keysAndArgs ...interface{}) (interface{}, error) {
	keys := make([]string, script.KeyCount)
	args := keysAndArgs

	if script.KeyCount > 0 {
		for i := 0; i < script.KeyCount; i++ {
			keys[i] = keyGen(keysAndArgs[i].(string))
		}
		args = keysAndArgs[script.KeyCount:]
	}

	v, err := cacheHandler.client.EvalSha(cacheHandler.ctx, script.Hash, keys, args...).Result()

	if err != nil && strings.HasPrefix(err.Error(), "NOSCRIPT ") {
		v, err = cacheHandler.client.Eval(cacheHandler.ctx, script.Src, keys, args...).Result()
	}
	return v, err

}

func (cacheHandler *CacheHandler) GetLock(key string) error {
	key = keyGen(key)

	tries := 200

	for i := 0; i < tries; i++ {
		select {
		default:
			if cacheHandler.client == nil {
				return errors.New("redis 已关闭")
			}
			if flag, err := cacheHandler.client.SetNX(cacheHandler.ctx, key, time.Now().Unix(), time.Duration(10)*time.Second).Result(); err != nil || !flag {
				time.Sleep(1 * time.Millisecond)
				continue
			} else {
				return nil
			}
		}
	}
	return errors.New("无法获得锁")
}

// tips 持有锁 才有权释放锁
func (cacheHandler *CacheHandler) ReleaseLock(key string) error {
	key = keyGen(key)
	// delete(lockKey, key)
	if err := cacheHandler.client.Del(cacheHandler.ctx, key).Err(); err != nil {
		log.Printf("释放锁失败，错误信息: %v", err)
		return err
	}
	return nil
}

func (cacheHandler *CacheHandler) ClearLock(ctx context.Context) {
	// 已在本机获取到lock的去掉
	// for key, _ := range lockKey {
	// 	cache.ReleaseLock(key)
	// }
}
func noErrNil(err error) error {
	if err == rds.Nil {
		return nil
	} else {
		return err
	}
}

// 清空
func (cacheHandler *CacheHandler) FlushDB() error {
	err := cacheHandler.client.FlushDB(cacheHandler.ctx).Err()
	if err != nil {
		return errors.New("FlushDB失败，错误信息: " + err.Error())
	}
	return nil
}
