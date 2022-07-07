package main

import (
	"fmt"
	"testing"
)
/*
defer 将它们注册到其所在 Goroutine 中，用于存放 deferred 函数的栈数据结构中，这些 deferred 函数将在执行 defer 的函数退出前，按后进先出（LIFO）的顺序被程序调度执行
deferred 函数是一个可以在任何情况下为函数进行收尾工作的好“伙伴”
 */
func TestDefer1(t *testing.T) {
	defer func() {
		println("defer1")
	}()
	defer func() {
		println("defer2")
	}()
	defer func() {
		println("defer3")
	}()
	println("main end")
}

func TestDefer2(t *testing.T) {
	for i := 0; i <= 3; i++ {
		defer fmt.Println(i)
	}
}
/*
当 foo3 返回后，deferred 函数被调度执行时，上述压入栈的 deferred 函数将以 LIFO 次序出栈执行。匿名函数会以闭包的方式访问外围函数的变量 i，并通过 Println 输出 i 的值，此时 i 的值为 4
 */
func TestDefer3(t *testing.T) {
	for i := 0; i <= 3; i++ {
		defer func() {
			fmt.Println(i)
		}()
	}
}

func sum(max int) int {
	total := 0
	for i := 0; i < max; i++ {
		total += i
	}
	return total
}
func fooWithDefer() {
	defer func() {
		sum(10)
	}()
}
func fooWithoutDefer() {
	sum(10)
}
func BenchmarkFooWithDefer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fooWithDefer()
	}
}
func BenchmarkFooWithoutDefer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fooWithoutDefer()
	}
}