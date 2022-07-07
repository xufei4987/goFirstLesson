package main

import (
	"fmt"
	"io"
	"os"
	"testing"
)

func myAppend(sl []int, elems ...int) []int {
	fmt.Printf("%T\n", elems) // []int
	if len(elems) == 0 {
		println("no elems to append")
		return sl
	}
	sl = append(sl, elems...)
	return sl
}

func TestFunction(t *testing.T) {
	sl := []int{1, 2, 3}
	sl = myAppend(sl) // no elems to append
	fmt.Println(sl)   // [1 2 3]
	sl = myAppend(sl, 4, 5, 6)
	fmt.Println(sl) // [1 2 3 4 5 6]
}

var myFprintf = func(w io.Writer, format string, a ...interface{}) (int, error) {
	return fmt.Fprintf(w, format, a...)
}

func TestFunction1(t *testing.T) {
	myFprintf(os.Stdout, "%s\n", "hello go")
}

//具名返回值 在返回值上声明 可以直接返回
func myfunc(a, b int) (c, d int) {
	c = a + b
	d = a - b
	return
}
func TestFunction2(t *testing.T) {
	r1, r2 := myfunc(1, 2)
	fmt.Println(r1, r2)
}

func setup(task string) func() {
	fmt.Println("do some setup stuff for", task)
	return func() {
		println("do some teardown stuff for", task)
	}
}

func TestFunction3(t *testing.T) {
	teardown := setup("demo")
	defer teardown()
	fmt.Println("do something")
}

//通过闭包的特性，减少入参的个数
func partialTimes(x int) func(int) int {
	return func(y int) int {
		return x * y
	}
}

func TestFunction4(t *testing.T) {
	timesTwo := partialTimes(2)
	fmt.Println(timesTwo(10))
	fmt.Println(timesTwo(20))
}

func change(s *string) {
	*s = "123"
}

func TestFunction5(t *testing.T) {
	s := "111"
	change(&s)
	fmt.Println(s)
}
