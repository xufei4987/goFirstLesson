package main

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

func dumpBytesArray(arr []byte) {
	fmt.Printf("[")
	for _, b := range arr {
		fmt.Printf("%c ", b)
	}
	fmt.Printf("]\n")
}

func TestString(t *testing.T) {
	var s = "hello"
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&s)) // 将string类型变量地址显式转型为reflect.StringHeader
	fmt.Printf("0x%x\n", hdr.Data) // 0x10a30e0
	fmt.Printf("%d\n", hdr.Len) // 0x10a30e0
	p := (*[5]byte)(unsafe.Pointer(hdr.Data)) // 获取Data字段所指向的数组的指针
	dumpBytesArray((*p)[:]) // [h e l l o ]   // 输出底层数组的内容
}

func TestStringFor(t *testing.T) {
	var s = "中国人"
	for i := 0; i < len(s); i++ {
		fmt.Printf("index: %d, value: 0x%x\n", i, s[i])
	}
	for i, v := range s {
		fmt.Printf("index: %d, value: 0x%x\n", i, v)
	}
}

func TestStringTrans(t *testing.T) {
	var s string = "中国人"

	// string -> []rune
	rs := []rune(s)
	fmt.Printf("%x\n", rs) // [4e2d 56fd 4eba] 字符切片 unicode码点

	// string -> []byte
	bs := []byte(s)
	fmt.Printf("%x\n", bs) // e4b8ade59bbde4baba 字节切片 utf8编码

	// []rune -> string
	s1 := string(rs)
	fmt.Println(s1) // 中国人

	// []byte -> string
	s2 := string(bs)
	fmt.Println(s2) // 中国人
}
