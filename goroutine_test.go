package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

/*
资源占用小，每个 goroutine 的初始栈大小仅为 2k；
由 Go 运行时而不是操作系统调度，goroutine 上下文切换在用户层完成，开销更小；
在语言层面而不是通过标准库提供。goroutine 由go关键字创建，一退出就会被回收或销毁，开发体验更佳；
语言内置 channel 作为 goroutine 间通信原语，为并发设计提供了强大支撑。
*/
func TestGoroutine1(t *testing.T) {
	c := make(chan int)
	go func(a, b int) {
		time.Sleep(1 * time.Second)
		c <- a + b
	}(1, 2)
	println(<-c)
}

func deadloop() {
	for {
	}
}
func TestGoroutine2(t *testing.T) {
	go deadloop()
	for {
		time.Sleep(time.Second * 1)
		fmt.Println("I got scheduled!")
	}
}

/*
由于无缓冲 channel 的运行时层实现不带有缓冲区，所以 Goroutine 对无缓冲 channel 的接收和发送操作是同步的。
也就是说，对同一个无缓冲 channel，只有对它进行接收操作的 Goroutine 和对它进行发送操作的 Goroutine 都存在的情况下，通信才能得以进行，否则单方面的操作会让对应的 Goroutine 陷入挂起状态
*/
func TestGoroutine3(t *testing.T) {
	ch1 := make(chan int)
	ch1 <- 13 // fatal error: all goroutines are asleep - deadlock!
	n := <-ch1
	println(n)
}

func TestGoroutine4(t *testing.T) {
	ch1 := make(chan int)
	go func() {
		ch1 <- 13
	}()
	n := <-ch1
	println(n)
}

/*
对一个带缓冲 channel 来说，在缓冲区未满的情况下，对它进行发送操作的 Goroutine 并不会阻塞挂起；在缓冲区有数据的情况下，对它进行接收操作的 Goroutine 也不会阻塞挂起。
但当缓冲区满了的情况下，对它进行发送操作的 Goroutine 就会阻塞挂起；当缓冲区为空的情况下，对它进行接收操作的 Goroutine 也会阻塞挂起
*/
func TestGoroutine5(t *testing.T) {
	ch2 := make(chan int, 1)
	n := <-ch2 // 由于此时ch2的缓冲区中无数据，因此对其进行接收操作将导致goroutine挂起
	println(n)
	ch3 := make(chan int, 1)
	ch3 <- 17 // 向ch3发送一个整型数17
	ch3 <- 27 // 由于此时ch3中缓冲区已满，再向ch3发送数据也将导致goroutine挂起
}

func TestGoroutine6(t *testing.T) {
	//ch1 := make(chan<- int, 1)  // 只发送channel类型
	//ch2 := make(<-chan int, 1) // 只接收channel类型

	//<-ch1       // invalid operation: <-ch1 (receive from send-only type chan<- int)
	//ch2 <- 13   // invalid operation: ch2 <- 13 (send to receive-only type <-chan int)
}

func produce(ch chan<- int) {
	for i := 0; i < 10; i++ {
		ch <- i
		time.Sleep(time.Second)
	}
	//channel 关闭后，所有等待从这个 channel 接收数据的操作都将返回
	//channel 的一个使用惯例，那就是发送端负责关闭 channel , why?
	//1. 发送端没有像接受端那样的、可以安全判断 channel 是否被关闭了的方法
	//2. 一旦向一个已经关闭的 channel 执行发送操作，这个操作就会引发 panic
	close(ch)
}

func consume(ch <-chan int) {
	//for range 会阻塞在对 channel 的接收操作上，直到 channel 中有数据可接收或 channel 被关闭循环，才会继续向下执行。channel 被关闭后，for range 循环也就结束了
	for n := range ch {
		println(n)
	}
}

/*
n := <- ch      // 当ch被关闭后，n将被赋值为ch元素类型的零值
m, ok := <-ch   // 当ch被关闭后，m将被赋值为ch元素类型的零值, ok值为false
for v := range ch { // 当ch被关闭后，for range循环结束
    ... ...
}
*/
func TestGoroutine7(t *testing.T) {
	ch := make(chan int, 5)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		produce(ch)
		wg.Done()
	}()
	go func() {
		consume(ch)
		wg.Done()
	}()
	wg.Wait()
}

type signal struct{}

func work() {
	println("worker is working")
	time.Sleep(time.Second)
}

func spawn(f func()) <-chan signal {
	c := make(chan signal)
	go func() {
		println("woker start to work")
		f()
		c <- signal{}
	}()
	return c
}

/*
 1 对 1 的通知机制
*/
func TestGoroutine8(t *testing.T) {
	println("start a worker...")
	c := spawn(work)
	<-c
	println("worker work done!")
}

func worker(i int) {
	fmt.Printf("worker %d: is working...\n", i)
	time.Sleep(1 * time.Second)
	fmt.Printf("worker %d: works done\n", i)
}

func spawnGroup(f func(i int), num int, groupSignal <-chan signal) <-chan signal {
	c := make(chan signal)
	var wg sync.WaitGroup

	for i := 0; i < num; i++ {
		wg.Add(1)
		go func(i int) {
			<-groupSignal
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

/*
 1 对 n 的“广播”机制
*/
func TestGoroutine9(t *testing.T) {
	fmt.Println("start a group of workers...")
	groupSignal := make(chan signal)
	c := spawnGroup(worker, 5, groupSignal)
	time.Sleep(time.Second)
	fmt.Println("the group of workers start to work...")
	close(groupSignal)
	<-c
	fmt.Println("the group of workers work done!")
}

type counter struct {
	c chan int
	i int
}

func NewCounter() *counter {
	cter := &counter{
		c: make(chan int),
	}
	go func() {
		for {
			cter.i++
			cter.c <- cter.i
		}
	}()
	return cter
}

func (cter *counter) Increase() int {
	return <-cter.c
}

/*
将计数器操作全部交给一个独立的 Goroutine 去处理，并通过无缓冲 channel 的同步阻塞特性，实现了计数器的控制
这种并发设计逻辑更符合 Go 语言所倡导的“不要通过共享内存来通信，而是通过通信来共享内存”的原则
*/
func TestGoroutine10(t *testing.T) {
	cter := NewCounter()
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			v := cter.Increase()
			fmt.Printf("goroutine-%d: current counter value is %d\n", i, v)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

var active = make(chan struct{}, 3)
var jobs = make(chan int, 10)

/*
用作计数信号量（counting semaphore）,控制同时运行的goroutine数量
*/
func TestGoroutine11(t *testing.T) {
	go func() {
		for i := 0; i < 8; i++ {
			jobs <- i + 1
		}
		close(jobs) //关闭一个通道，该通道必须是双向的或仅发送的。 它应该只由发送者执行(生产者)。
		// 关闭通道后无法再向通道中发送数据，关闭后的通道的最后一个值被接收后，再次接收会立即返回零值（x, ok := <-c，c为零值，ok为false）
	}()
	wg := sync.WaitGroup{}
	for job := range jobs {
		wg.Add(1)
		go func(j int) {
			active <- struct{}{}
			fmt.Printf("handle job: %d\n", j)
			time.Sleep(2 * time.Second)
			<-active
			wg.Done()
		}(job)
	}
	wg.Wait()
}

func trySend(c chan<- int, i int) bool {
	select {
	case c <- i:
		return true
	default:
		return false
	}
}

func tryRecv(c <-chan int) (int, bool) {
	select {
	case i := <-c:
		return i, true
	default:
		return 0, false
	}
}

func producer(c chan<- int) {
	i := 1
	for {
		time.Sleep(2 * time.Second)
		if ok := trySend(c, i); ok {
			fmt.Printf("[producer]: send [%d] to channel\n", i)
			i++
			continue
		}
		fmt.Printf("[producer]: try send [%d], but channel is full\n", i)
	}
}

func consumer(c <-chan int) {
	for {
		if i, ok := tryRecv(c); !ok {
			fmt.Println("[consumer]: try to recv from channel, but the channel is empty")
			time.Sleep(1 * time.Second)
			continue
		} else {
			fmt.Printf("[consumer]: recv [%d] from channel\n", i)
			if i >= 3 {
				fmt.Println("[consumer]: exit")
				return
			}
		}
	}
}

/*
由于用到了 select 原语的 default 分支语义，当 channel 空的时候，tryRecv 不会阻塞；当 channel 满的时候，trySend 也不会阻塞
有多个发送者，但有且只有一个接收者。在这样的场景下，我们可以在接收 goroutine 中使用len(channel)是否大于0来判断是否 channel 中有数据需要接收
有多个接收者，但有且只有一个发送者。在这样的场景下，我们可以在发送 Goroutine 中使用len(channel)是否小于cap(channel)来判断是否可以执行向 channel 的发送操作
*/
func TestGoroutine12(t *testing.T) {
	var wg sync.WaitGroup
	c := make(chan int, 3)
	wg.Add(2)
	go func() {
		producer(c)
		wg.Done()
	}()
	go func() {
		consumer(c)
		wg.Done()
	}()
	wg.Wait()
}

/*
如果一个 channel 类型变量的值为 nil，我们称它为 nil channel。nil channel 有一个特性，那就是对 nil channel 的读写都会发生阻塞
*/
func TestGoroutine13(t *testing.T) {
	ch1, ch2 := make(chan int), make(chan int)
	go func() {
		time.Sleep(5 * time.Second)
		ch1 <- 5
		close(ch1)
	}()

	go func() {
		time.Sleep(7 * time.Second)
		ch2 <- 7
		close(ch2)
	}()

	for {
		//这里已经被置为 nil 的 c1 或 c2 的分支，将再也不会被 select 选中执行
		select {
		case x, ok := <-ch1:
			if !ok {
				ch1 = nil
			} else {
				fmt.Println(x)
			}
		case x, ok := <-ch2:
			if !ok {
				ch2 = nil
			} else {
				fmt.Println(x)
			}
		}
		if ch1 == nil && ch2 == nil {
			break
		}
	}
	fmt.Println("program end")
}

/*
channel 与 select 结合使用的一些惯用法
第一种用法：利用 default 分支避免阻塞
select 语句的 default 分支的语义，就是在其他非 default 分支因通信未就绪，而无法被选择的时候执行的，这就给 default 分支赋予了一种“避免阻塞”的特性
func tryRecv(c <-chan int) (int, bool) {
  select {
  case i := <-c:
    return i, true
  default: // channel为空
    return 0, false
  }
}
func trySend(c chan<- int, i int) bool {
  select {
  case c <- i:
    return true
  default: // channel满了
    return false
  }
}
第二种用法：实现超时机制
func worker() {
  select {
  case <-c:
       // ... do some stuff
  case <-time.After(30 *time.Second):
      return
  }
}
第三种用法：实现心跳机制
func worker() {
  heartbeat := time.NewTicker(30 * time.Second)
  defer heartbeat.Stop()
  for {
    select {
    case <-c:
      // ... do some stuff
    case <- heartbeat.C:
      //... do heartbeat stuff
    }
  }
}
*/

func TestGoroutine14(t *testing.T) {
	c := make(chan struct{})
	go func() {
		time.Sleep(4 * time.Second)
		c <- struct{}{}
	}()
	select {
	case <-c:
		// ... do some stuff
		fmt.Println("do some stuff")
	case <-time.After(5 * time.Second):
		fmt.Println("timeout")
	}
}
