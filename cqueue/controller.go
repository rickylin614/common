package cqueue

import (
	"context"
	"errors"
	"runtime"
	"sync"
	"time"
)

var c = &controller{
	acceptNew: true,
	ctx:       context.Background(),
}

func AddHandler(h Handler) {
	c.AddHandler(h)
}

func RunHandlers(ctx context.Context) {
	c.RunHandlers(ctx)
}

func ShutdownSignal(timeout time.Duration) error {
	c.ShutdownSignal()

	tc := time.Tick(timeout)

	done := make(chan struct{})

	go func() {
		c.lock.Lock()
		count := c.count
		c.lock.Unlock()

		// 等待 ctx 完成或 count 降至 0
		if count == 0 {
			done <- struct{}{}
			return
		}

		runtime.Gosched()
	}()

	for {
		select {
		case <-tc:
			return errors.New("shutdown timeout")
		case <-done:
			return nil
		}
	}
}

type controller struct {
	count     int
	lock      sync.Mutex
	acceptNew bool
	handlers  []Handler
	ctx       context.Context
}

type Handler struct {
	handlerFunc func(ctx context.Context)
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
		handlerFunc: func(ctx context.Context) {
			for {
				c.lock.Lock()
				canAccept := c.acceptNew
				c.lock.Unlock()

				// 只有在 acceptNew 为 true 时才执行
				if canAccept {
					c.Increment()
					h.handlerFunc(ctx)
					c.Decrement()
				} else {
					return
				}
			}
		},
	}

	c.handlers = append(c.handlers, wrappedHandler)
}

func (c *controller) RunHandlers(ctx context.Context) {
	for _, handler := range c.handlers {
		for i := 0; i < handler.worker; i++ {
			go handler.handlerFunc(ctx) // 启动每个处理器
		}
	}
}
