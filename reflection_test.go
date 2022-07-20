package main

import (
	"fmt"
	"reflect"
	"testing"
)

/*
反射第一定律：反射可以将“接口类型变量”转换为“反射类型对象”( reflect.TypeOf(x) reflect.ValueOf(x) )
反射第二定律：反射可以将“反射类型对象”转换为“接口类型变量”
反射第三定律：如果要修改“反射类型对象”其值必须是“可写的”
种类（Kind）指的是对象归属的品种，在 reflect 包中有如下定义：
   Invalid Kind = iota  // 非法类型
   Bool                 // 布尔型
   Int                  // 有符号整型
   Int8                 // 有符号8位整型
   Int16                // 有符号16位整型
   Int32                // 有符号32位整型
   Int64                // 有符号64位整型
   Uint                 // 无符号整型
   Uint8                // 无符号8位整型
   Uint16               // 无符号16位整型
   Uint32               // 无符号32位整型
   Uint64               // 无符号64位整型
   Uintptr              // 指针
   Float32              // 单精度浮点数
   Float64              // 双精度浮点数
   Complex64            // 64位复数类型
   Complex128           // 128位复数类型
   Array                // 数组
   Chan                 // 通道
   Func                 // 函数
   Interface            // 接口
   Map                  // 映射
   Ptr                  // 指针
   Slice                // 切片
   String               // 字符串
   Struct               // 结构体
   UnsafePointer        // 底层指针
Map、Slice、Chan 属于引用类型，使用起来类似于指针，但是在种类常量定义中仍然属于独立的种类，不属于 Ptr
type A struct{} 定义的结构体属于 Struct 种类，*A 属于 Ptr
*/
// 定义一个Enum类型
type Enum int

const (
	Zero Enum = 0
)

func TestReflection1(z *testing.T) {
	// 声明一个空结构体
	type cat struct {
	}
	// 获取结构体实例的反射类型对象
	typeOfCat := reflect.TypeOf(cat{})
	// 显示反射类型对象的名称和种类
	fmt.Println(typeOfCat.Name(), typeOfCat.Kind())
	// 获取Zero常量的反射类型对象
	typeOfA := reflect.TypeOf(Zero)
	// 显示反射类型对象的名称和种类
	fmt.Println(typeOfA.Name(), typeOfA.Kind())
}

func TestReflection2(z *testing.T) {
	// 声明一个空结构体
	type cat struct {
	}
	// 创建cat的实例
	ins := &cat{}
	// 获取结构体实例的反射类型对象
	typeOfCat := reflect.TypeOf(ins)
	// 显示反射类型对象的名称和种类
	fmt.Printf("name:'%v' kind:'%v'\n", typeOfCat.Name(), typeOfCat.Kind())
	// 取类型的元素
	typeOfCat = typeOfCat.Elem()
	// 显示反射类型对象的名称和种类
	fmt.Printf("element name: '%v', element kind: '%v'\n", typeOfCat.Name(), typeOfCat.Kind())

}

func TestReflection3(t *testing.T) {
	type cat struct {
		Name string
		Type int `json:"type" id:"100"`
	}
	ins := cat{
		Name: "mm",
		Type: 1,
	}
	typeOfCat := reflect.TypeOf(ins)
	// 遍历结构体所有成员
	for i := 0; i < typeOfCat.NumField(); i++ {
		// 获取每个成员的结构体字段类型
		fieldType := typeOfCat.Field(i)
		// 输出成员名和tag
		fmt.Printf("name: %v  tag: '%v'\n", fieldType.Name, fieldType.Tag)
	}
	// 通过字段名, 找到字段类型信息
	if catType, ok := typeOfCat.FieldByName("Type"); ok {
		// 从tag中取出需要的tag
		fmt.Println(catType.Tag.Get("json"), catType.Tag.Get("id"))
	}
}

/*
编写 Tag 时，必须严格遵守键值对的规则。结构体标签的解析代码的容错能力很差，一旦格式写错，编译和运行时都不会提示任何错误
*/
func TestReflection4(t *testing.T) {
	type cat struct {
		Name string
		Type int `json: "type" id:"100"`
	}
	typeOfCat := reflect.TypeOf(cat{})
	if catType, ok := typeOfCat.FieldByName("Type"); ok {
		fmt.Println(catType.Tag.Get("json"))
	}
}

func TestReflection5(t *testing.T) {
	type cat struct {
		Name string
		Type int `json:"type" id:"100"`
	}
	ins := cat{
		Name: "mm",
		Type: 1,
	}
	val := reflect.ValueOf(ins)
	if v, ok := val.Interface().(cat); ok {
		fmt.Println(v.Name, v.Type)
	} else {
		println("val is not cat")
	}
	fmt.Println(val.Interface())
}

func TestReflection6(t *testing.T) {
	type cat struct {
		Name string
		Type int `json:"type" id:"100"`
	}
	var c cat
	typeOfC := reflect.TypeOf(c)
	ins := reflect.New(typeOfC) //创建这个类型的实例值，值以 reflect.Value 类型返回; 这步操作等效于：new(int)，因此返回的是 *int 类型的实例
	fmt.Println(ins.Type(), ins.Kind(), ins.Elem().Interface())
}

func mul(a, b int) int {
	return a * b
}

func TestReflection7(t *testing.T) {
	funcValue := reflect.ValueOf(mul)

	params := []reflect.Value{reflect.ValueOf(10), reflect.ValueOf(20)}

	rets := funcValue.Call(params)

	for _, ret := range rets {
		fmt.Println("ret is :", ret.Int())
	}

}
