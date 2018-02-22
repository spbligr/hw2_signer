package main

import (
	"sync"
	"fmt"
	"sort"
	"strings"
)

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
	wg := &sync.WaitGroup{}
	hashChan := make(chan string)

	for data := range in {
		wg.Add(1)
		data := fmt.Sprintf("%v", data)
		hashMd5 := DataSignerMd5(data)
		go func(data string, hashMd5 string) {
			hashChan <- DataSignerCrc32(data)+ "~" + DataSignerCrc32(hashMd5)
			//@todo избавиться от костыля
			if data == "8" {
				close(hashChan)
			}
			defer wg.Done()
		}(data, hashMd5)
	}

	for hashResult := range hashChan {
		out <- hashResult
	}

	defer wg.Wait()
}

func MultiHash(in chan interface{}, out chan interface{})  {
	wg := &sync.WaitGroup{}
	outTemp := make(chan string)

	count := 0
	for singleHash := range in {
		wg.Add(1)
		count++
		go func(outTemp chan string, singleHash interface{}, count int) {
			defer wg.Done()
			//@todo избавиться от костыля
			if count == 7 {
				defer close(outTemp)
			}

			var hashResult string
			for i:=0; i < 6; i ++ {
				hashResult = hashResult + DataSignerCrc32(fmt.Sprintf("%v%v", i, singleHash))
			}

			outTemp <- hashResult
		}(outTemp, singleHash, count)
	}

	for hashResult := range outTemp {
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