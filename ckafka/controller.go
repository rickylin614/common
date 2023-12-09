package ckafka

import (
	"context"
	"sync"
	"time"
)

var c = &controller{
	acceptNew: true,
}

func AddHandler(h Handler) {
	c.AddHandler(h)
}

func RunHandlers() {
	c.RunHandlers()
}

func ShutdownSignal(ctx context.Context) {
	c.ShutdownSignal()

	// 等待 ctx 完成或 count 降至 0
	for {
		c.lock.Lock()
		count := c.count
		c.lock.Unlock()

		if count == 0 {
			return // 如果 count 为 0，退出函数
		}

		select {
		case <-ctx.Done():
			return // 如果 context 完成，退出函数
		default:
			time.Sleep(time.Millisecond * 100) // 短暂等待再次检查
		}
	}
}

type controller struct {
	count     int
	lock      sync.Mutex
	acceptNew bool
	handlers  []Handler
}

type Handler struct {
	handlerFunc func()
	worker      int
}

func (c *controller) Increment() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.count++
}

func (c *controller) Decrement() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.count--
}

func (c *controller) ShutdownSignal() {
	c.lock.Lock()
	c.acceptNew = false
	c.lock.Unlock()
}

func (c *controller) AddHandler(h Handler) {
	wrappedHandler := Handler{
		worker: h.worker,
		handlerFunc: func() {
			for {
				c.lock.Lock()
				canAccept := c.acceptNew
				c.lock.Unlock()

				// 只有在 acceptNew 为 true 时才执行
				if canAccept {
					c.Increment()
					h.handlerFunc()
					c.Decrement()
				} else {
					return
				}
			}
		},
	}

	c.handlers = append(c.handlers, wrappedHandler)
}

func (c *controller) RunHandlers() {
	for _, handler := range c.handlers {
		for i := 0; i < handler.worker; i++ {
			go handler.handlerFunc() // 启动每个处理器
		}
	}
}
