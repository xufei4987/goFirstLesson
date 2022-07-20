package workerpool2

import (
	"errors"
	"fmt"
	"sync"
)

type Task func()

type Pool struct {
	capacity int // workerpool大小
	active   chan struct{}
	tasks    chan Task
	wg       sync.WaitGroup // 用于在pool销毁时等待所有worker退出
	quit     chan struct{}  // 用于通知各个worker退出的信号channel
	block    bool           // 没有worker后是否阻塞
}

var (
	defaultCapacity       = 2
	maxCapacity           = 16
	ErrWorkerPoolFreed    = errors.New("workerpool1 freed")      // workerpool已终止运行
	ErrNoIdleWorkerInPool = errors.New("no idle worker in pool") // workerpool中任务已满，没有空闲goroutine用于处理新任务
)

func New(capacity int, options ...Option) *Pool {
	if capacity <= 0 {
		capacity = defaultCapacity
	}
	if capacity > maxCapacity {
		capacity = maxCapacity
	}

	p := &Pool{
		capacity: capacity,
		tasks:    make(chan Task),
		quit:     make(chan struct{}),
		active:   make(chan struct{}, capacity),
	}

	for _, opt := range options {
		opt(p)
	}

	fmt.Println("workerpool2 start")

	go p.run()

	return p
}

func (p *Pool) run() {
	idx := 0

	for {
		select {
		case <-p.quit:
			return
		case p.active <- struct{}{}:
			idx++
			p.newWorker(idx)
		}
	}
}

func (p *Pool) newWorker(idx int) {
	p.wg.Add(1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("worker[%03d]: recover panic[%s] and exit\n", idx, err)
				<-p.active
			}
			p.wg.Done()
		}()

		fmt.Printf("worker[%03d]: start\n", idx)

		for {
			select {
			case <-p.quit:
				fmt.Printf("worker[%03d]: exit\n", idx)
				<-p.active
				return
			case t := <-p.tasks:
				fmt.Printf("worker[%03d]: receive a task\n", idx)
				t()
			}
		}
	}()
}

func (p *Pool) Schedule(t Task) error {
	select {
	case <-p.quit:
		return ErrWorkerPoolFreed
	case p.tasks <- t:
		return nil
	default:
		if p.block {
			p.tasks <- t
			return nil
		}
		return ErrNoIdleWorkerInPool
	}
}

func (p *Pool) Free() {
	close(p.quit)
	p.wg.Wait()
	fmt.Printf("workerpool1 freed\n")
}
