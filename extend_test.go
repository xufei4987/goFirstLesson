package main

import (
	"fmt"
	"io"
	"strings"
	"testing"
)

type A interface {
	A1()
	A2()
}

//type B interface {
//	A1()
//	A2()
//	B1()
//}
//用接口类型 A 替代上面接口类型 B 定义中 M1 和 M2
type B interface {
	A //接口类型的类型嵌入
	B1()
}


type MyInt int
func (n *MyInt) Add(m int) {
	*n = *n + MyInt(m)
}
type t struct {
	a int
	b int
}
type s struct {
	*MyInt
	t
	io.Reader
	s string
	n int
}
/*
结构体类型嵌入 的底层原理  实际上是静态代理 Read 代理给 Reader.Read
 */
func TestExtend1(a *testing.T) {
	m := MyInt(17)
	r := strings.NewReader("hello,go")
	s := s{
		MyInt:&m,
		t: t{
			a: 1,
			b: 2,
		},
		Reader: r,
		s:      "demo",
	}

	var sl = make([]byte, len("hello, go"))
	//s.Reader.Read(sl)  //类型嵌入“实现继承”
	s.Read(sl)
	fmt.Println(string(sl)) // hello, go
	//s.MyInt.Add(5)  //类型嵌入“实现继承”
	s.Add(5)
	fmt.Println(*(s.MyInt)) // 22
}


type e1 interface {
	E1()
	E2()
	E3()
}

type e2 interface {
	E1()
	E2()
	E4()
}

type e3 struct {
	e1
	e2
}

func (e3) E1(){
	println("e3 E1")
}
func (e3) E2(){
	println("e3 E2")
}

func TestExtend2(a *testing.T) {
	e := e3{}
	// 方法集合存在交集,编译器会因无法做出选择而报错,需要e3实现冲突的方法
	e.E1()
	e.E2()
}

type Tt1 struct{}
func (Tt1) Tt1M1()   { println("Tt1's M1") }
func (*Tt1) PTt1M2() { println("PTt1's M2") }
type Tt2 struct{}
func (Tt2) Tt2M1()   { println("Tt2's M1") }
func (*Tt2) PTt2M2() { println("PTt2's M2") }
type Tt3 struct {
	Tt1
	*Tt2
}

func TestExtend3(t *testing.T) {
	tt3 := Tt3{
		Tt1: Tt1{},
		Tt2: &Tt2{},
	}
	dumpMethodSet(tt3)
	//*Tt3 类型的方法集合，它包含的可不是 Tt1 类型的方法集合，而是 *Tt1 类型的方法集合
	dumpMethodSet(&tt3)
}


type BT1 int
type bt2 struct{
	n int
	m int
}
type BI interface {
	X1()
}
type S1 struct {
	BT1
	*bt2
	BI
	a int
	b string
}
type S2 struct {
	BT1 BT1
	bt2 *bt2
	BI  BI
	a  int
	b  string
}

func TestExtend4(t *testing.T) {
	bt2 := bt2{
		n:1,m:2,
	}
	var bi BI
	s1 := S1{
		BT1:1,
		bt2:&bt2,
		BI:bi,
		a:10,
		b:"10",
	}
	s2 := S2{
		BT1:1,
		bt2:&bt2,
		BI:bi,
		a:10,
		b:"10",
	}
	dumpMethodSet(s1)
	dumpMethodSet(&s1)

	dumpMethodSet(s2)
	dumpMethodSet(&s2)
}