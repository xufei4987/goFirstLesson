package main

import (
	"fmt"
	"testing"
)

const (
	Apple, Banana = 11, 22
	Strawberry, Grape
	Pear, Watermelon
)
/*
等价于
const (
    Apple, Banana = 11, 22
    Strawberry, Grape  = 11, 22 // 使用上一行的初始化表达式
    Pear, Watermelon  = 11, 22 // 使用上一行的初始化表达式
)

同一行的 iota 即便出现多次，多个 iota 的值也是一样的
const (
    Apple, Banana = iota, iota + 10 // 0, 10 (iota = 0)
    Strawberry, Grape // 1, 11 (iota = 1)
    Pear, Watermelon  // 2, 12 (iota = 2)
)
 */
func TestConst(t *testing.T) {
	fmt.Printf("%d, %d\n",Apple,Banana)
	fmt.Printf("%d, %d\n",Strawberry,Grape)
	fmt.Printf("%d, %d\n",Pear,Watermelon)
}

/*
iota:代表枚举中的偏移量，第一行为0，第二行为1，以此类推
 */
const (
	mutexLocked = 1 << iota //1 << 0 = 1
	mutexWoken //1<<1 = 2
	mutexStarving //1<<2 =4
	mutexWaiterShift = iota //3
	starvationThresholdNs = 1e6
)

func TestConstIota(t *testing.T) {
	fmt.Printf("%d\n",mutexLocked)
	fmt.Printf("%d\n",mutexWoken)
	fmt.Printf("%d\n",mutexStarving)
	fmt.Printf("%d\n",mutexWaiterShift)
	fmt.Printf("%f\n",starvationThresholdNs)
}

const (
	_ = iota // 0
	Pin1
	Pin2
	Pin3
	_
	Pin5    // 5
)

func TestConstIota1(t *testing.T) {
	fmt.Printf("%d\n",Pin1)
	fmt.Printf("%d\n",Pin2)
	fmt.Printf("%d\n",Pin3)
	fmt.Printf("%d\n",Pin5)
}
