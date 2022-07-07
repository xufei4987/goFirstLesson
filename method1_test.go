package main

import (
	"fmt"
	"testing"
	"time"
)

/*
receiver 参数的基类型本身不能为 指针类型 或 接口类型
 */
//type Myint *int
//
//func (r Myint) String() string {
//	return "123"
//}
//
//type MyReader io.Reader
//
//func (r MyReader) Read(p []byte) (int, error) {
//	return r.Read(p)
//}
/*
方法声明要与 receiver 参数的基类型声明放在同一个包内:
1、我们不能为原生类型（诸如 int、float64、map 等）添加方法
2、不能跨越 Go 包为其他包的类型声明新方法
 */
//func (i int) Foo() string { // 编译器报错：cannot define new methods on non-local type int
//	return fmt.Sprintf("%d", i)
//}
//
//func (s http.Server) Foo() { // 编译器报错：cannot define new methods on non-local type http.Server
//}

type T1 struct {
	a int
}
func (t T1) Get() int {
	return t.a
}
func (t *T1) Set(a int) int {
	t.a = a
	return t.a
}

// 类型T的方法Get的等价函数
func Get(t T1) int {
	return t.a
}
// 类型*T的方法Set的等价函数
func Set(t *T1, a int) int {
	t.a = a
	return t.a
}

func TestMethod1(t *testing.T) {
	var t1 T1
	println(t1.Get())
	t1.Set(1)
	println(t1.a)

	//等价替换
	var t2 T1
	println(T1.Get(t2))
	(*T1).Set(&t2,1)
	println(t2.a)
}

func TestMethod2(t *testing.T) {
	var t1 T1
	f1 := (*T1).Set // f1的类型，也是T类型Set方法的类型：func (t *T, int)int
	f2 := T1.Get    // f2的类型，也是T类型Get方法的类型：func(t T)int
	fmt.Printf("the type of f1 is %T\n", f1) // the type of f1 is func(*main.T, int) int
	fmt.Printf("the type of f2 is %T\n", f2) // the type of f2 is func(main.T) int
	f1(&t1, 3)
	fmt.Println(f2(t1)) // 3
}

type field1 struct {
	name string
}
func (p *field1) print() {
	fmt.Println(p.name)
}
/*
迭代 data1 时，由于 data1 中的元素类型是 field 指针 (*field)，因此赋值后 v 就是元素地址，与 print 的 receiver 参数类型相同，
每次调用 (*field).print 函数时直接传入的 v 即可，实际上传入的也是各个 field 元素的地址；

迭代 data2 时，由于 data2 中的元素类型是 field（非指针），与 print 的 receiver 参数类型不同，因此需要将其取地址后再传入 (*field).print 函数。
这样每次传入的 &v 实际上是变量 v 的地址，而不是切片 data2 中各元素的地址。
 */
func TestMethod3(t *testing.T) {
	data1 := []*field1{{"one"}, {"two"}, {"three"}}
	for _, v := range data1 {
		//go v.print() //等价于
		go (*field1).print(v)
	}
	data2 := []field1{{"four"}, {"five"}, {"six"}}
	for _, v := range data2 {
		//go v.print() //等价于
		go (*field1).print(&v)
	}
	time.Sleep(3 * time.Second)
}


type field2 struct {
	name string
}
func (p field2) print() {
	fmt.Println(p.name)
}
func TestMethod4(t *testing.T) {
	data1 := []*field2{{"one"}, {"two"}, {"three"}}
	for _, v := range data1 {
		//go v.print() //等价于
		go field2.print(*v)
	}
	data2 := []field2{{"four"}, {"five"}, {"six"}}
	for _, v := range data2 {
		//go v.print() //等价于
		go field2.print(v)
	}
	time.Sleep(3 * time.Second)
}