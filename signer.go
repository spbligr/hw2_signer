package main

import (
	"sync"
	"fmt"
	"sort"
	"strings"
)

type multiHashResult struct {
	number int
	hash string
}

const TH_COUNT  = 6

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
			//@todo надеюсь за этот костыль я не попаду в ад
			if data == "8" {
				defer close(hashChan)
			}
			defer wg.Done()
			crt32DataChan := getCrt32Data(data)
			right32 := DataSignerCrc32(hashMd5)
			left32 := <- crt32DataChan
			hashChan <- left32 + "~" + right32
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

	for input := range in {
		wg.Add(1)

		wgTemp := &sync.WaitGroup{}
		data := input.(string)
		inCh := make(chan multiHashResult)

		wgTemp.Add(TH_COUNT)
		for i := 0; i < TH_COUNT; i++ {
			go GetMultiHashProcess(inCh, wgTemp, data, i)
		}
		go func(wgInner *sync.WaitGroup, c chan multiHashResult) {
			defer close(c)
			wgInner.Wait()
		}(wgTemp, inCh)

		go sortMultiResults(inCh, outTemp, wg)

	}

	go func(wgOut *sync.WaitGroup, c chan string) {
		defer close(c)
		wgOut.Wait()
	}(wg, outTemp)

	for hash := range outTemp {
		out <- hash
	}

}

func GetMultiHashProcess(hashResultChan chan multiHashResult, wg *sync.WaitGroup, singleHash interface{}, i int)  {
	defer wg.Done()
	hashResultChan <- multiHashResult{number: i, hash: DataSignerCrc32(fmt.Sprintf("%v%v", i, singleHash))}
}

func sortMultiResults(hashResults chan multiHashResult, out chan string, wg *sync.WaitGroup)  {
	defer wg.Done()
	result := map[int]string{}
	var data []int

	for hashResult := range hashResults {
		result[hashResult.number] = hashResult.hash
		data = append(data, hashResult.number)
	}
	sort.Ints(data)

	var results []string
	for i := range data {
		results = append(results, result[i])
	}

	out <- strings.Join(results, "")
}

func getCrt32Data(data string) chan string  {
	result := make(chan string, 1)
	go func(out chan <- string) {
		out <- DataSignerCrc32(data)
	}(result)

	return result
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