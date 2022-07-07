package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	//var s = `         ,_---~~~~~----._
  //  _,,_,*^____      _____*g*\"*,--,
  // / __/ /'     ^.  /      \ ^@q   f
  //[  @f | @))    |  | @))   l  0 _/
  // \/   \~____ / __ \_____/     \
  //  |           _l__l_           I
  //  }          [______]           I
  //  ]            | | |            |
  //  ]             ~ ~             |
  //  |                            |
  //   |                           |`
	//fmt.Println(s)

	var s = "中国人"
	fmt.Printf("the length of s = %d\n", len(s)) // 9
	for i := 0; i < len(s); i++ {
		fmt.Printf("0x%x ", s[i]) // 0xe4 0xb8 0xad 0xe5 0x9b 0xbd 0xe4 0xba 0xba
	}
	fmt.Printf("\n")

	fmt.Println("the character count in s is", utf8.RuneCountInString(s)) // 3
	for _, c := range s {
		fmt.Printf("0x%x ", c) // 0x4e2d 0x56fd 0x4eba
		fmt.Printf("%c ", c) // 0x4e2d 0x56fd 0x4eba
	}
	fmt.Printf("\n")
	fmt.Printf("%c \n", '\u4e2d')
	encodeRune()
	decodeRune()
}

// rune -> []byte
func encodeRune() {
	var r rune = 0x4E2D
	fmt.Printf("the unicode charactor is %c\n", r) // 中
	buf := make([]byte, 3)
	_ = utf8.EncodeRune(buf, r) // 对rune进行utf-8编码
	fmt.Printf("utf-8 representation is 0x%X\n", buf) // 0xE4B8AD
}

// []byte -> rune
func decodeRune() {
	var buf = []byte{0xE4, 0xB8, 0xAD}
	r, _ := utf8.DecodeRune(buf) // 对buf进行utf-8解码
	fmt.Printf("the unicode charactor after decoding [0xE4, 0xB8, 0xAD] is %s\n", string(r)) // 中
}