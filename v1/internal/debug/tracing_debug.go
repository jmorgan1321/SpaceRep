// +build debug

package debug

import (
	"fmt"
	"runtime"
	"strings"
)

func Trace() {
	IndentationLevel.Increment()

	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		fmt.Println("what??")
		return
	}
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		fmt.Println("what??")
		return
	}

	s := strings.Split(fn.Name(), "/")
	fmt.Printf("%s%s()\n", IndentationLevel, s[len(s)-1])
}

func UnTrace() {
	IndentationLevel.Decrement()
}
