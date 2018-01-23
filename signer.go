package main

import (
	//"fmt"
)

func ExecutePipeline(freeFlowJobs... job)  {
	ChIn := make(chan interface{})
	ChOut := make(chan interface{})

	for _, job := range freeFlowJobs {
		go job(ChIn, ChOut)
	}


}

var SingleHash = func(in, out chan interface{}) {
	value := "1"
	out <- DataSignerCrc32(value)+"~"+DataSignerCrc32(DataSignerMd5(value))
}

var MultiHash = func(in, out chan interface{}) {

}

var CombineResults = func(in, out chan interface{}) {

}