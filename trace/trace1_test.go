package trace

import (
	"fmt"
	"runtime"
	"testing"
)

func Trace(name string) func() {
	println("enter:", name)
	return func() {
		println("exit:", name)
	}
}
func foo() {
	defer Trace("foo")()
	bar()
}
func bar() {
	defer Trace("bar")()
}

func TestTrace1(t *testing.T) {
	defer Trace("main")()
	foo()
}

func Trace1() func() {
	caller, file, line, ok := runtime.Caller(1)
	if !ok {
		panic("not found caller")
	}
	fn := runtime.FuncForPC(caller)
	name := fn.Name()
	fmt.Printf("file:%s line:%d enter:%s\n", file, line, name)
	return func() {
		fmt.Println("exit:", name)
	}
}
func foo1() {
	defer Trace1()()
	bar1()
}
func bar1() {
	defer Trace1()()
}

func TestTrace(t *testing.T) {
	defer Trace1()()
	foo1()
}
