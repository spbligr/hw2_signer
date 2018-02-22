package main

import (
	"fmt"
	"time"
	"sync"
)

func main() {
	tempChan := make(chan int)

	wg := &sync.WaitGroup{}

	for i:=0; i < 5; i++ {
		wg.Add(1)
		go func(in chan int, i int) {
			fmt.Println("Отравботал джоб", i)
			time.Sleep(2 * time.Second)
			wg.Done()
		}(tempChan, i)
	}

	wg.Wait()


}
