package main

import (
	//"fmt"
	//"fmt"
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

var SingleHash = func(in, out chan interface{}) {
	for value := range out {
		data := fmt.Sprintf("SingleHash value = %v", value)
		fmt.Println(data)
		//result := DataSignerCrc32(data)+"~"+DataSignerCrc32(DataSignerMd5(data))
		//fmt.Println("SingleHash result = ", result)
		//in <- result
	}
}

var MultiHash = func(in, out chan interface{}) {
	for value := range out {
		data := fmt.Sprintf("MultiHash value = %v", value)
		fmt.Println(data)
	}
}

var CombineResults = func(in, out chan interface{}) {

}