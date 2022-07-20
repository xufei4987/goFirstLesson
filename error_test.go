package main

import (
	"errors"
	"fmt"
	"testing"
)

type MyError struct {
	Msg1 string
	Msg2 string
	Err  error
}

func (e MyError) Error() string {
	return e.Err.Error()
}

var (
	err1 = errors.New("err1")
	err2 = MyError{
		Msg1: "m1",
		Msg2: "m2",
		Err:  err1,
	}
)

func TestError1(t *testing.T) {
	fmt.Println(err2.Error())
	fmt.Println(err2.Msg1)
	fmt.Println(err2.Msg2)
}

func isErr2Error(err error) bool {
	if myErr, ok := err.(MyError); ok && myErr.Msg1 == "m1" {
		return true
	}
	return false
}

func TestError2(t *testing.T) {
	if isErr2Error(err2) {
		fmt.Println("err2 is MyError")
	} else {
		fmt.Println("err2 is not MyError")
	}
}

var ErrSentinel = errors.New("the underlying sentinel error")

func TestError3(t *testing.T) {
	err1 := fmt.Errorf("wrap sentinel: %w", ErrSentinel)
	err2 := fmt.Errorf("wrap err1: %w", err1)
	println(err1 == ErrSentinel) //false
	println(err2 == ErrSentinel) //false
	if errors.Is(err2, ErrSentinel) {
		println("err2 is ErrSentinel")
		return
	}
	println("err2 is not ErrSentinel")
}

type MyError1 struct {
	e string
}

func (e *MyError1) Error() string {
	return e.e
}
func TestError4(t *testing.T) {
	var err = &MyError1{"MyError error demo"}
	err1 := fmt.Errorf("wrap err: %w", err)
	err2 := fmt.Errorf("wrap err1: %w", err1)
	var e *MyError1
	if errors.As(err2, &e) {
		println("MyError is on the chain of err2")
		println(e)
		println(err)
		println(e == err)
		return
	}
	println("MyError is not on the chain of err2")
}

type NotFoundError struct {
	File string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("file %q not found", e.File)
}

func open(file string) error {
	return &NotFoundError{File: file}
}

func TestError5(t *testing.T) {
	if err := open("testfile.txt"); err != nil {
		var notFound *NotFoundError
		if errors.As(err, &notFound) {
			// handle the error
			fmt.Println("error hapend:", notFound)
		} else {
			panic("unknown error")
		}
	}
}
