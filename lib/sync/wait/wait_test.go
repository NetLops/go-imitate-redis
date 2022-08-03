package wait

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func work() {
	w := Wait{}
	defer w.Done()
	w.Add(1)
	w.WaitWithTimeout(100 * time.Microsecond)
}

func TestWait(t *testing.T) {
	go func() {
		for true {
			fmt.Printf("NumGoroutine = %d\n", runtime.NumGoroutine())
			time.Sleep(time.Second)
		}
	}()
	for {
		work()
	}
}
