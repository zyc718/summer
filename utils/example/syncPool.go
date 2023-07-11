package example

import (
	"summer/utils/log"
	"sync"
)

/*
   sync.Pool simple Example
   作用:减少GC，利用重用机制，构建Pool对象池
   特性：临时对象池、并发安全、无通知回收
*/
// 定义一个结构体

type PoolStruct struct {
	Name string
	Age  int
}

// 初始化sync.Pool new 函数就是创建 PoolStruct结构体
func initPool() *sync.Pool {
	return &sync.Pool{
		New: func() interface{} {
			return &PoolStruct{}
		},
	}
}

func SyncPool() {
	pool := initPool()
	poolstruct := pool.Get().(*PoolStruct)

	log.Infof("这是第一次从sync.pool 获取到的poolStruct %+v", poolstruct)
	poolstruct.Name = "zyc"
	poolstruct.Age = 18
	pool.Put(poolstruct)
	log.Infof("设置的poolStruct.name%v", poolstruct.Name)
	log.Infof("设置的poolStruct.age%v", poolstruct.Age)
	log.Infof("Pool 中有一个对象，调用Get方法获取%+v", pool.Get().(*PoolStruct))
	log.Infof("Pool 中没有对象了，再次调用Get方法获取%+v", pool.Get().(*PoolStruct))
}
