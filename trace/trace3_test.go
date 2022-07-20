package trace

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

func printTrace(id uint64, name, arrow string, indent int) {
	indents := ""
	for i := 0; i < indent; i++ {
		indents += "    "
	}
	fmt.Printf("g[%05d]:%s%s%s\n", id, indents, arrow, name)
}

var mu sync.Mutex
var m = make(map[uint64]int)

func Trace3() func() {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		panic("not found caller")
	}
	fn := runtime.FuncForPC(pc)
	name := fn.Name()
	gid := curGoroutineID()
	mu.Lock()
	indents := m[gid]    // 获取当前gid对应的缩进层次
	m[gid] = indents + 1 // 缩进层次+1后存入map
	mu.Unlock()
	printTrace(gid, name, "->", indents+1)
	return func() {
		mu.Lock()
		indents := m[gid]    // 获取当前gid对应的缩进层次
		m[gid] = indents - 1 // 缩进层次-1后存入map
		mu.Unlock()
		printTrace(gid, name, "<-", indents)
	}
}

func A31() {
	defer Trace3()()
	B31()
}
func B31() {
	defer Trace3()()
	C31()
}
func C31() {
	defer Trace3()()
	D33()
}
func D33() {
	defer Trace3()()
}
func A32() {
	defer Trace3()()
	B32()
}
func B32() {
	defer Trace3()()
	C32()
}
func C32() {
	defer Trace3()()
	D33()
}

func TestTrace31(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		A32()
		wg.Done()
	}()
	time.Sleep(1)
	A31()
	wg.Wait()
}
