package main

import (
	//"log"
	//"fmt"
	//"fmt"
	//"fmt"
	//"fmt"
)



func main() {

	// println("run as\n\ngo test -v -race")
	inputData := []int{0, 1, 1, 2, 3, 5, 8}
	//testResult := "NOT_SET"

	hashSignJobs := []job{
		job(func(in, out chan interface{}) {
			for _, fibNum := range inputData {
				out <- fibNum
			}
		}),
		job(SingleHash),
		job(MultiHash),
		//job(CombineResults),
		//job(func(in, out chan interface{}) {
		//	dataRaw := <-in
		//	data, ok := dataRaw.(string)
		//	if !ok {
		//		log.Fatal("cant convert result data to string")
		//	}
		//	testResult = data
		//}),
	}

	ExecutePipeline(hashSignJobs...)

	//fmt.Scanln()
}

