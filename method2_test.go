package main

import (
	"fmt"
	"reflect"
	"testing"
)

/*
方法和函数的等价替换，receiver作为第一个参数传入函数（相当于this）
func (t T) M1() <=> F1(t T)
func (t *T) M2() <=> F2(t *T)
 */
type T21 struct {
	a int
}
func (t T21) M1() {
	t.a = 10
}
func (t *T21) M2() {
	t.a = 11
}
/*
首先，当 receiver 参数的类型为 T 时：
当我们选择以 T 作为 receiver 参数类型时，M1 方法等价转换为F1(t T)。我们知道，Go 函数的参数采用的是值拷贝传递，也就是说，F1 函数体中的 t 是 T 类型实例的一个副本。
这样，我们在 F1 函数的实现中对参数 t 做任何修改，都只会影响副本，而不会影响到原 T 类型实例。

第二，当 receiver 参数的类型为 *T 时：
当我们选择以 *T 作为 receiver 参数类型时，M2 方法等价转换为F2(t *T)。同上面分析，我们传递给 F2 函数的 t 是 T 类型实例的地址，这样 F2 函数体中对参数 t 做的任何修改，都会反映到原 T 类型实例上。
 */
func TestMethod21(t *testing.T) {
	var t21 T21
	println(t21.a) // 0
	t21.M1()
	println(t21.a) // 0
	p := &t21
	p.M2()
	//t21.M2()
	println(t21.a) // 11
}

/*
原则1：如果 Go 方法要把对 receiver 参数代表的类型实例的修改，反映到原类型实例上，那么我们应该选择 *T 作为 receiver 参数的类型

无论是 T 类型实例，还是 *T 类型实例，都既可以调用 receiver 为 T 类型的方法，也可以调用 receiver 为 *T 类型的方法。
这样，我们在为方法选择 receiver 参数的类型的时候，就不需要担心这个方法不能被与 receiver 参数类型不一致的类型实例调用了

原则2：一般情况下，我们通常会为 receiver 参数选择 T 类型，因为这样可以缩窄外部修改类型实例内部状态的“接触面”，也就是尽量少暴露可以修改类型内部状态的方法。
不过也有一个例外需要你特别注意。考虑到 Go 方法调用时，receiver 参数是以值拷贝的形式传入方法中的。
那么，如果 receiver 参数类型的 size 较大，以值拷贝形式传入就会导致较大的性能开销，这时我们选择 *T 作为 receiver 类型可能更好些
 */

func TestMethod22(t *testing.T) {
	var t1 T21
	println(t1.a) // 0
	t1.M1()
	println(t1.a) // 0
	//(&t1).M2()
	t1.M2() //t1.M2() 这种用法是 Go 提供的“语法糖”：Go 判断 t1 的类型为 T，也就是与方法 M2 的 receiver 参数类型 *T 不一致后，会自动将t1.M2()转换为(&t1).M2()
	println(t1.a) // 11

	var t2 = &T21{}
	println(t2.a) // 0
	//(*t2).M1()
	t2.M1() //Go 判断 t2 的类型为 *T，与方法 M1 的 receiver 参数类型 T 不一致，就会自动将t2.M1()转换为(*t2).M1()
	println(t2.a) // 0
	t2.M2()
	println(t2.a) // 11
}


type Interface interface {
	M1()
	M2()
}

type TT struct{}

func (t TT) M1()  {}
func (t *TT) M2() {}

func TestMethod23(t *testing.T) {
	//var tt TT
	//var pt *TT
	//var i Interface
	//
	//i = pt
	//TT 没有实现 Interface 类型方法列表中的 M2，因此类型 TT 的实例 tt 不能赋值给 Interface 变量
	//i = tt
}

func dumpMethodSet(i interface{}) {
	dynTyp := reflect.TypeOf(i)
	if dynTyp == nil {
		fmt.Printf("there is no dynamic type\n")
		return
	}
	n := dynTyp.NumMethod()
	if n == 0 {
		fmt.Printf("%s's method set is empty!\n", dynTyp)
		return
	}
	fmt.Printf("%s's method set:\n", dynTyp)
	for j := 0; j < n; j++ {
		fmt.Println("-", dynTyp.Method(j).Name)
	}
	fmt.Printf("\n")
}

type SS TT
type AA = TT
type II Interface
/*
Go 语言规定，*TT 类型的方法集合包含所有以 *TT 为 receiver 参数类型的方法，以及所有以 TT 为 receiver 参数类型的方法
 */
func TestMethod24(t *testing.T) {
	var s string
	dumpMethodSet(s)
	dumpMethodSet(&s)
	var tt TT
	dumpMethodSet(tt)
	dumpMethodSet(&tt)
	var ss SS
	dumpMethodSet(ss)
	dumpMethodSet(&ss)
	var aa AA
	dumpMethodSet(aa)
	dumpMethodSet(&aa)
	var ii II
	dumpMethodSet(ii)
	dumpMethodSet(&ii)

}