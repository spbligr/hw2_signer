package main

import (
	"sync"
	"fmt"
	"sort"
	"strings"
)

type SingleHashResult struct {
	Hash string
	Number string
}

func ExecutePipeline(hashSignJobs ...job) {
	wg := &sync.WaitGroup{}
	in := make(chan interface{})

	for _, jobItem := range hashSignJobs {
		wg.Add(1)
		out := make(chan interface{})
		go func(jobFunc job, in chan interface{}, out chan interface{}, wg *sync.WaitGroup) {
			defer wg.Done()
			defer close(out)
			jobFunc(in, out)
		}(jobItem, in, out, wg)
		in = out
	}

	defer wg.Wait()
}

func SingleHash(in chan interface{}, out chan interface{}) {
	//wg := &sync.WaitGroup{}
	//defer wg.Wait()
	//
	for data := range in {
		number := fmt.Sprintf("%v", data)
		result := DataSignerCrc32(number)+ "~" + DataSignerCrc32(DataSignerMd5(number))
		out <- SingleHashResult{Hash: result, Number: number}
	}
}

func MultiHash(in chan interface{}, out chan interface{})  {
	for th := range in {
		var hashResult string
		for i:=0; i < 6; i ++ {
			hashResult = hashResult + DataSignerCrc32(fmt.Sprintf("%v%v", i, (th).(SingleHashResult).Hash))
		}
		out <- hashResult
	}
}

func CombineResults(in, out chan interface{}){

	var hashResults []string
	var result string

	for hashResult := range in {
		hashResults = append(hashResults, (hashResult).(string))
	}

	sort.Strings(hashResults)

	result = strings.Join(hashResults, "_")

	out <- result
}