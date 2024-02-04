package brutepool

import (
	"context"
	"fmt"
	"sync"
)

type BrutePool struct {
	BruteList   []interface{}    //字典
	Concurrency int              //线程数
	queues      chan interface{} //数据通道

	BruteFunc       func(interface{}) bool //爆破函数
	SuccessCallBack func(interface{})      //爆破成功的回调函数
}

// 默认回调为打印结果
func defaultCallBack(v interface{}) {
	fmt.Printf("[SUCCESS] %v\n", v)
}

func New(list []interface{}, function func(interface{}) bool) *BrutePool {
	return &BrutePool{
		BruteList:       list,
		Concurrency:     4, //建议线程数为4
		BruteFunc:       function,
		SuccessCallBack: defaultCallBack,
		queues:          make(chan interface{}),
	}
}

// 往通道内压入字典数据
func (p *BrutePool) setDataToChan(ctx context.Context) {
	go func(ctx context.Context) {
		defer close(p.queues) //关闭通道
		for _, i := range p.BruteList {
			select {
			case <-ctx.Done(): //如果context取消，就结束该线程
				return
			default:
				p.queues <- i //压入数据
			}
		}
	}(ctx)
}

// 调用爆破函数，处理通道内的数据
func (p *BrutePool) doDataFunc(ctx context.Context, cancel context.CancelFunc) {
	var wg sync.WaitGroup //并发控制方式
	wg.Add(p.Concurrency)
	for i := 0; i < p.Concurrency; i++ {
		go func(ctx context.Context) {
			defer wg.Done()
			for v := range p.queues {
				select {
				case <-ctx.Done(): //如果context取消，就结束该线程
					return
				default:
					//调用爆破函数
					if p.BruteFunc(v) {
						p.doCallBack(ctx, v)
						cancel() //如果爆破成功，调用context取消函数
					}
				}
			}
		}(ctx)
	}
	wg.Wait()
}

// 调用回调函数
func (p *BrutePool) doCallBack(ctx context.Context, v interface{}) {
	select {
	case <-ctx.Done(): //如果context取消，就结束该线程
		return
	default:
		go p.SuccessCallBack(v)
	}
}

// 伙计，跑起来！
func (p *BrutePool) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	p.setDataToChan(ctx)
	p.doDataFunc(ctx, cancel)
}
