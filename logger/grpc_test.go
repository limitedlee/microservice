package golog

import (
	"sync"
	"testing"
	"time"
)

var wg sync.WaitGroup

func TestLog(t *testing.T) {
	ch := make(chan int, 10)
	for i := 0; i < 1; i++ {
		wg.Add(1)
		ch <- 1
		go func(ch chan int) {
			defer wg.Done()
			Error(1)
			<-ch
		}(ch)
	}
	wg.Wait()
	time.Sleep(time.Second * 3)
}
