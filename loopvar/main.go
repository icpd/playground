package main

import (
	"fmt"
	"runtime"
	"time"
)

// Go1.22: 乱序输出
// Go1.22之前：全是9
func main() {
	runtime.GOMAXPROCS(1)
	ch := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
		}
		time.Sleep(time.Second)
		close(ch)
	}()

	for i := range ch {
		go func() {
			runtime.Gosched()
			fmt.Println(i)
		}()
	}
}
