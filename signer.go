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

		in <- result //сюда как-то передать data и result
	}
}

var MultiHash = func(in, out chan interface{}) {
	for value := range in {
		fmt.Println("MultiHash", value)
	}

}

var CombineResults = func(in, out chan interface{}) {

}