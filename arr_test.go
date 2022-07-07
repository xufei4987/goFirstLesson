package main

import (
	"fmt"
	"testing"
	"unsafe"
)

func TestArr(t *testing.T) {
	var arr1 [5]int
	//var arr2 [6]int
	//var arr3 [5]string
	foo(arr1)
	//foo(arr2) //数组的元素类型+长度决定是否是同一种类型
	//foo(arr3) //数组的元素类型+长度决定是否是同一种类型
}

func foo(arr [5]int){

}

func TestArr1(t *testing.T) {
	var arr = [6]int{1, 2, 3, 4, 5, 6}
	fmt.Println("数组长度：", len(arr))           // 6
	fmt.Println("数组大小：", unsafe.Sizeof(arr)) // 48
}

func TestArr2(t *testing.T){
	var arr2 = [6]int {
		11, 12, 13, 14, 15, 16,
	} // [11 12 13 14 15 16]
	//用“…”替代，Go 编译器会根据数组元素的个数，自动计算出数组长度
	var arr3 = [...]int {
		21, 22, 23,
	} // [21 22 23]
	fmt.Printf("%T\n", arr2) // [6]int
	fmt.Printf("%T\n", arr3) // [3]int

	var arr4 = [...]int{
		99: 39, // 将第100个元素(下标值为99)的值赋值为39，其余元素值均为0
	}
	fmt.Printf("%T\n", arr4) // [100]int

	var arr = [6]int{11, 12, 13, 14, 15, 16}
	fmt.Println(arr[0], arr[5]) // 11 16
	//fmt.Println(arr[-1])        // 错误：下标值不能为负数
	//fmt.Println(arr[8])         // 错误：小标值超出了arr的长度范围
}