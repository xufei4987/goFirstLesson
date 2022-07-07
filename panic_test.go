package main

import (
	"fmt"
	"testing"
)

func foo1() {
	defer func() {
		if e := recover(); e!=nil{
			fmt.Println("recover the panic:",e)
		}
	}()
	println("call foo")
	bar()
	println("exit foo")
}

func bar() {
	//defer func() {
	//	if e := recover(); e!=nil{
	//		fmt.Println("recover the panic:",e)
	//	}
	//}()
	println("call bar")
	panic("panic occurs in bar")
	zoo()
	println("exit bar")
}

func zoo() {
	println("call zoo")
	println("exit zoo")
}

func TestPanic(t *testing.T) {
	println("call main")
	foo1()
	println("exit main")
}
