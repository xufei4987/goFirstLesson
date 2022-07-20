package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

/*
无论是在单 Goroutine 情况下，还是在并发测试情况下，sync.Mutex实现的同步机制的性能，都要比 channel 实现的高出三倍多
通常在需要高性能的临界区（critical section）同步机制的情况下，sync 包提供的低级同步原语更为适合
*/
var cs = 0 // 模拟临界区要保护的数据
var mu sync.Mutex
var c = make(chan struct{}, 1)

func criticalSectionSyncByMutex() {
	mu.Lock()
	cs++
	mu.Unlock()
}
func criticalSectionSyncByChan() {
	c <- struct{}{}
	cs++
	<-c
}
func BenchmarkCriticalSectionSyncByMutex(b *testing.B) {
	for n := 0; n < b.N; n++ {
		criticalSectionSyncByMutex()
	}
}
func BenchmarkCriticalSectionSyncByMutexInParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			criticalSectionSyncByMutex()
		}
	})
}
func BenchmarkCriticalSectionSyncByChan(b *testing.B) {
	for n := 0; n < b.N; n++ {
		criticalSectionSyncByChan()
	}
}
func BenchmarkCriticalSectionSyncByChanInParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			criticalSectionSyncByChan()
		}
	})
}

/*
如果对使用过的、sync 包中的类型的示例进行复制，并使用了复制后得到的副本，将导致不可预期的结果。
所以，在使用 sync 包中的类型的时候，我们推荐通过闭包方式，或者是传递类型实例（或包裹该类型的类型实例）的地址（指针）的方式进行。
这就是使用 sync 包时最值得我们注意的事项
sync.Mutex //互斥锁
sync.RWMutex //读写锁
sync.Cond //条件变量

sync 包中的低级同步原语各有各的擅长领域，你可以记住：
在具有一定并发量且读多写少的场合使用 RWMutex；
在需要“等待某个条件成立”的场景下使用 Cond；
当你不确定使用什么原语时，那就使用 Mutex 吧。
如果你对同步的性能有极致要求，且并发量较大，读多写少，那么可以考虑一下 atomic 包提供的原子操作函数。
*/

var ready bool

func worker1(i int) {
	fmt.Printf("worker %d: is working...\n", i)
	time.Sleep(1 * time.Second)
	fmt.Printf("worker %d: works done\n", i)
}

func spawnGroup1(f func(i int), num int, groupSignal *sync.Cond) <-chan signal {
	c := make(chan signal)
	var wg sync.WaitGroup

	for i := 0; i < num; i++ {
		wg.Add(1)
		go func(i int) {
			groupSignal.L.Lock()
			for !ready {
				groupSignal.Wait()
			}
			groupSignal.L.Unlock()
			fmt.Printf("worker %d: start to work...\n", i)
			f(i)
			wg.Done()
		}(i + 1)
	}

	go func() {
		wg.Wait()
		c <- signal{}
	}()
	return c
}

func TestSync1(t *testing.T) {
	fmt.Println("start a group of workers...")
	groupSignal := sync.NewCond(&sync.Mutex{})
	c := spawnGroup1(worker1, 5, groupSignal)

	time.Sleep(5 * time.Second) // 模拟ready前的准备工作
	fmt.Println("the group of workers start to work...")

	groupSignal.L.Lock()
	ready = true
	groupSignal.Broadcast()
	groupSignal.L.Unlock()
	<-c
	fmt.Println("the group of workers work done!")
}
