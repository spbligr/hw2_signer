package main

import (
	"fmt"
)

func ExecutePipeline(freeFlowJobs... job)  {
	in := make(chan interface{})
	out := make(chan interface{})

	for _, job := range freeFlowJobs {
		go job(in, out)
	}

	fmt.Scanln()

}

var SingleHash = func(in, out chan interface{})  {
	for value := range out {

		hashData := make(map[string]string)

		data := fmt.Sprintf("%s", value)
		result := DataSignerCrc32(hashData["data"]) + "~" + DataSignerCrc32(DataSignerMd5(hashData["data"]))

		fmt.Println("SingleHash data= , resul = ", data, result)

		in <- result
		bufferCh <- data
		fmt.Println("doshol ba ba ba")
	}
}

var MultiHash = func(in, out chan interface{}) {
	LOOP:
	for {
		select {
		case singleHashResult := <- in :
			fmt.Println("MultiHash singleHashResult =", singleHashResult)
		case data := <-bufferCh :
			fmt.Println("MultiHash data =", data)
		default:
			fmt.Println("MultiHash data =", "azino 777")
			break LOOP
		}
	}

}

var CombineResults = func(in, out chan interface{}) {

}