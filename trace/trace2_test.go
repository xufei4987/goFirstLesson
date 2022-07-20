package trace

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"testing"
)

var goroutineSpace = []byte("goroutine ")

func curGoroutineID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	// Parse the 4707 out of "goroutine 4707 ["
	b = bytes.TrimPrefix(b, goroutineSpace)
	i := bytes.IndexByte(b, ' ')
	if i < 0 {
		panic(fmt.Sprintf("No space found in %q", b))
	}
	b = b[:i]
	n, err := strconv.ParseUint(string(b), 10, 64)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse goroutine ID out of %q: %v", b, err))
	}
	return n
}

func Trace2() func() {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		panic("not found caller")
	}
	fn := runtime.FuncForPC(pc)
	name := fn.Name()
	gid := curGoroutineID()
	fmt.Printf("g[%05d]: enter: [%s]\n", gid, name)
	return func() { fmt.Printf("g[%05d]: exit: [%s]\n", gid, name) }
}

func A1() {
	defer Trace2()()
	B1()
}
func B1() {
	defer Trace2()()
	C1()
}
func C1() {
	defer Trace2()()
	D()
}
func D() {
	defer Trace2()()
}
func A2() {
	defer Trace2()()
	B2()
}
func B2() {
	defer Trace2()()
	C2()
}
func C2() {
	defer Trace2()()
	D()
}

func TestTrace21(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		A2()
		wg.Done()
	}()
	A1()
	wg.Wait()
}
