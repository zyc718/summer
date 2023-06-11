package cache

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	rds "github.com/go-redis/redis/v8"
	"io"
	"time"
)

var Pool CachePool

type CachePool interface {
	//获取链接
	Get(ctx context.Context) (CacheHandler, error)
	//	关闭缓存
	Close() error
	//	检查缓存是否开启
	IsOpen() bool
}

type CacheHandler interface {
	// General
	// 获取缓存值
	Get(key string) (string, error)
	// 设置缓存值
	Set(key string, value interface{}, expiration time.Duration) error
	// 设置缓存值
	SetNx(key string, value interface{}, expiration time.Duration) (bool, error)
	// 删除缓存值
	Del(key string) error
	// 设置哈希消息
	HSet(channel, key string, value interface{}) error
	// 获取哈希消息
	HGet(channel, key string) (string, error)
	// 获取全部哈希消息
	HGetAll(channel string) (map[string]string, error)
	// 删除哈希数据
	HDel(channel, key string) error
	// 有序队列添加元素
	ZAdd(channel, key string, score float64) error
	// 有序集合长度
	ZCard(channel string) (int64, error)
	// 获取序列值
	ZRange(channel string, start, stop int64) ([]string, error)
	// 获取序列的值和分数
	ZRangeWithScores(channel string, start, stop int64) ([]rds.Z, error)
	// 删除序列值
	ZRem(channel, key string) error
	// 执行lua脚本
	Eval(script *Script, keysAndArgs ...interface{}) (interface{}, error)
	// 获得锁
	GetLock(key string) error
	// 释放锁
	ReleaseLock(key string) error
	// 释放本地派发的锁
	ClearLock(ctx context.Context)
	// 清空
	FlushDB() error
}

type Script struct {
	KeyCount int
	Src      string
	Hash     string
}

func NewScript(keyCount int, src string) *Script {
	h := sha1.New()
	_, _ = io.WriteString(h, src)
	return &Script{keyCount, src, hex.EncodeToString(h.Sum(nil))}
}

// 消息队列配置
type cacheConfig struct {
	// 消息队列具体配置
	Adapters map[string]json.RawMessage `json:"adapters"`
}

func CacheClose(ctx context.Context) error {
	// db, _ := Pool.Get(ctx)
	// // 清理掉所有本机派发的锁
	// db.ClearLock(ctx)
	if err := Pool.Close(); err != nil {
		return err
	}
	return nil
}
