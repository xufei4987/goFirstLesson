package main

import (
	"errors"
	"fmt"
	"testing"
)

func TestInterface(t *testing.T) {
	var i interface{}
	i = "astring"
	/*
		类型断言: v, ok := i.(T)
		如果接口类型变量 i 之前被赋予的值确为 T 类型的值，那么这个语句执行后，左侧“comma, ok”语句中的变量 ok 的值将为 true，变量 v 的类型为 T，它值会是之前变量 i 的右值。
		如果 i 之前被赋予的值不是 T 类型的值，那么这个语句执行后，变量 ok 的值为 false，变量 v 的类型还是那个要还原的类型，但它的值是类型 T 的零值
	*/
	if v, ok := i.(string); ok {
		fmt.Printf("type:%T value:%v\n", v, v)
	}
}

func TestInterface1(t *testing.T) {

	var a int64 = 13
	var i interface{} = a
	v1, ok := i.(int64)
	fmt.Printf("v1=%d, the type of v1 is %T, ok=%t\n", v1, v1, ok) // v1=13, the type of v1 is int64, ok=true
	v2, ok := i.(string)
	fmt.Printf("v2=%s, the type of v2 is %T, ok=%t\n", v2, v2, ok) // v2=, the type of v2 is string, ok=false
	v3 := i.(int64)
	fmt.Printf("v3=%d, the type of v3 is %T\n", v3, v3) // v3=13, the type of v3 is int64
	v4 := i.([]int)                                     // panic: interface conversion: interface {} is int64, not []int
	fmt.Printf("the type of v4 is %T\n", v4)

}

type MyErrorx struct {
	error
}

var ErrBad = MyErrorx{
	error: errors.New("bad things happened"),
}

func bad() bool {
	return false
}

//func returnsError() error {
//	var p *MyErrorx = nil
//	if bad() {
//		p = &ErrBad
//	}
//	return p //返回的是 tab._type为MyErrorx，data为nil的指针变量
//}
func returnsError() error {
	if bad() {
		return &ErrBad
	}
	return nil
}

func TestInterface2(t *testing.T) {
	err := returnsError()
	if err != nil {
		fmt.Printf("error occur: %+v\n", err)
		return
	}
	fmt.Println("ok")
}

/*
第一种：nil 接口变量
无论是空接口类型还是非空接口类型变量，一旦变量值为 nil，那么它们内部表示均为(0x0,0x0)，也就是类型信息、数据值信息均为空
*/
func TestInterface3(t *testing.T) {
	// nil接口变量
	var i interface{} // 空接口类型
	var err error     // 非空接口类型
	println(i)
	println(err)
	println("i = nil:", i == nil)
	println("err = nil:", err == nil)
	println("i = err:", i == err)
}

/*
第二种：空接口类型变量
对于空接口类型变量，只有 _type 和 data 所指数据内容一致的情况下，两个空接口类型变量之间才能划等号
*/
func TestInterface4(t *testing.T) {
	var eif1 interface{} // 空接口类型
	var eif2 interface{} // 空接口类型
	var n, m int = 17, 18

	eif1 = n
	eif2 = m
	println("eif1:", eif1)
	println("eif2:", eif2)
	println("eif1 = eif2:", eif1 == eif2) // false

	eif2 = 17
	println("eif1:", eif1)
	println("eif2:", eif2)
	println("eif1 = eif2:", eif1 == eif2) // true

	eif2 = int64(17)
	println("eif1:", eif1)
	println("eif2:", eif2)
	println("eif1 = eif2:", eif1 == eif2) // false
}

/*
第三种：非空接口类型变量
空接口类型变量一样，只有 tab 和 data 指的数据内容一致的情况下，两个非空接口类型变量之间才能划等号
*/
type Tx int

func (t Tx) Error() string {
	return "bad error"
}
func TestInterface5(t *testing.T) {
	var err1 error // 非空接口类型
	var err2 error // 非空接口类型
	err1 = (*Tx)(nil)
	println("err1:", err1)
	println("err2:", err2)
	println("err1 = nil:", err1 == nil)
	err1 = Tx(5)
	err2 = Tx(6)
	println("err1:", err1)
	println("err2:", err2)
	println("err1 = err2:", err1 == err2)
	err2 = fmt.Errorf("%d\n", 5)
	println("err1:", err1)
	println("err2:", err2)
	println("err1 = err2:", err1 == err2)
}

/*
第四种：空接口类型变量与非空接口类型变量的等值比较
Go 在进行等值比较时，类型比较使用的是 eface 的 _type 和 iface 的 tab._type，因此就像我们在这个例子中看到的那样，当 eif 和 err 都被赋值为T(5)时，两者之间是划等号的
*/
func TestInterface6(t *testing.T) {
	var eif interface{} = Tx(5)
	var err error = Tx(5)
	println("eif:", eif)
	println("err:", err)
	println("eif = err:", eif == err)
	err = Tx(6)
	println("eif:", eif)
	println("err:", err)
	println("eif = err:", eif == err)
}

type Txx struct {
	n int
	s string
}

func (Txx) Mxx1() {}
func (Txx) Mxx2() {}

type NonEmptyInterface interface {
	Mxx1()
	Mxx2()
}

func TestInterface7(test *testing.T) {
	var t = Txx{
		n: 17,
		s: "hello, interface",
	}
	var ei interface{}
	ei = t

	var i NonEmptyInterface
	i = t
	fmt.Println(ei)
	fmt.Println(i)

	println("eif:", ei)
	println("err:", i)
	println("eif = err:", ei == i)

	var txx NonEmptyInterface
	txx = (*Txx)(nil)
	println("txx:", txx)
	println("txx == nil :", txx == nil)
}

type Add func(a, b int) int

func anotheradd(a, b int) int {
	return a + b
}

func adder(a, b int, add Add) int {
	return add(a, b)
}

func TestInterface8(test *testing.T) {
	res := adder(1, 2, Add(anotheradd))
	println(res)
}
