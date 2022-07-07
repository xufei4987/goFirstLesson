package main

import "testing"

func TestSwitch(t *testing.T) {
	ext := "md"
	switch ext {
	case "json":
		println("read json file")
	case "jpg", "jpeg", "png", "gif":
		println("read image file")
	case "txt", "md":
		println("read text file")
	case "yml", "yaml":
		println("read yaml file")
	case "ini":
		println("read ini file")
	default:
		println("unsupported file extension:", ext)
	}
}

type person struct {
	name string
	age  int
}

func TestSwitch2(t *testing.T) {
	p := person{"tom", 13}
	switch p {
	case person{"tony", 33}:
		println("match tony")
	case person{"tom", 13}:
		println("match tom")
	case person{"lucy", 23}:
		println("match lucy")
	default:
		println("no match")
	}
}

func case1() int {
	println("eval case1 expr")
	return 1
}
func case2() int {
	println("eval case2 expr")
	return 2
}
func switchexpr() int {
	println("eval switch expr")
	return 1
}

/*
使用fallthrough关键字并不会使下一个case的表达式执行，而是直接执行其代码分支
*/
func TestSwitch3(t *testing.T) {
	switch switchexpr() {
	case case1():
		println("exec case1")
		fallthrough
	case case2():
		println("exec case2")
		fallthrough
	default:
		println("exec default")
	}
}

func TestSwitch4(t *testing.T) {
	var x interface{} = 13
	switch x.(type) {
	case nil:
		println("x is nil")
	case int:
		println("the type of x is int")
	case string:
		println("the type of x is string")
	case bool:
		println("the type of x is string")
	default:
		println("don't support the type")
	}

	// v表示的是x的值，而不是x的类型
	switch v := x.(type) {
	case nil:
		println("v is nil")
	case int:
		println("the type of v is int, v =", v)
	case string:
		println("the type of v is string, v =", v)
	case bool:
		println("the type of v is bool, v =", v)
	default:
		println("don't support the type")
	}
}

type I interface {
	M()
}

type T struct {
}

func (T) M() {
}

type S struct {
}

func (S) M() {
}

func TestSwitch5(t1 *testing.T) {
	var t T
	var i I = t
	switch i.(type) {
	case T:
		println("it is type T")
	case S:
		println("it is type int")
		// case 后面的类型都只能是实现了接口 I 的类型
	//case int:
	//	println("it is type int")
	}
}
/*
同一函数内 break 语句所在的最内层的 for、switch 或 select
如果需要跳出for循环，可以利用带 label 的 break 语句
 */
func TestSwitch6(t *testing.T) {
	var sl = []int{5, 19, 6, 3, 8, 12}
	var firstEven int = -1
	// find first even number of the interger slice
	for i := 0; i < len(sl); i++ {
		switch sl[i] % 2 {
		case 0:
			firstEven = sl[i]
			break
		case 1:
			// do nothing
		}
	}
	println(firstEven)
}
