package oldpool

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type BrutePool struct {
	BruteList       []interface{}          //字典
	Concurrency     int                    //线程数
	BruteFunc       func(interface{}) bool //爆破函数
	SuccessCallBack func(interface{})      //爆破成功的回调函数

	success int32            //爆破是否成功
	queues  chan interface{} //数据通道
}

// 默认回调为打印结果
func defaultCallBack(v interface{}) {
	fmt.Printf("[SUCCESS] %v\n", v)
}

func New(list []interface{}, function func(interface{}) bool) *BrutePool {
	return &BrutePool{
		BruteList:       list,
		Concurrency:     3, //建议线程数为3
		BruteFunc:       function,
		SuccessCallBack: defaultCallBack,
		success:         0,
		queues:          make(chan interface{}),
	}
}

func (b *BrutePool) Run() {
	//步骤一：往通道内压入字典数据
	go func() {
		for _, i := range b.BruteList {
			//爆破未成功时，才往通道存数据。如果成功，就不继续存数据了。
			if atomic.LoadInt32(&b.success) == 1 {
				break
			}
			b.queues <- i
		}
		close(b.queues) // 必须关闭,不然阻塞死锁
	}()
	//步骤二：多并发爆破
	var wg sync.WaitGroup //并发控制方式
	wg.Add(b.Concurrency)
	for i := 0; i < b.Concurrency; i++ {
		go func() {
			defer wg.Done()
			//从通道取数据
			for v := range b.queues {
				if b.BruteFunc(v) {
					atomic.StoreInt32(&b.success, 1)
					b.SuccessCallBack(v)
					break
				}
				if atomic.LoadInt32(&b.success) == 1 {
					break
				}
			}
		}()
	}
	wg.Wait()
}
