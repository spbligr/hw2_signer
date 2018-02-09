package main

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func ExecutePipeline(hashSignJobs ...job) {
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	in := make(chan interface{})

	for _, jobItem := range hashSignJobs {

		wg.Add(1)

		out := make(chan interface{})

		go func(job job, in chan interface{}, out chan interface{}, wg *sync.WaitGroup) {
			job(in, out)
			defer wg.Done()
			defer close(out)
		}(jobItem, in, out, wg)

		in = out
	}

}

func SingleHash(in chan interface{}, out chan interface{}) {

	wg := &sync.WaitGroup{}

	tmp := make(chan string)

	for data := range in {

		fmt.Printf(" in %v", data.(int))


		wg.Add(1)

		raw, _ := data.(int)

		data := strconv.Itoa(raw)
		md5 := DataSignerMd5(data)

		go func(wg *sync.WaitGroup, out chan string, d string, md5 string) {
			hash := make(chan string)
			hashMd5 := make(chan string) // make second because there is no guarante which will calculated first, so must take control of order.

			go func(out chan string, input string) {
				out <- DataSignerCrc32(input)
			}(hash, d)

			go func(out chan string, input string) {
				out <- DataSignerCrc32(input)
			}(hashMd5, md5)

			out <- fmt.Sprintf("%v~%v", <-hash, <-hashMd5)
			wg.Done()
		}(wg, tmp, data, md5)
	}

	go func(wg *sync.WaitGroup, c chan string) {
		defer close(c)
		wg.Wait()
	}(wg, tmp)

	for hash := range tmp {
		out <- hash
	}
}

func MultiHash(in chan interface{}, out chan interface{}) {
	type hashNode struct {
		id    int
		value string
	}

	wgOut := &sync.WaitGroup{}
	outCh := make(chan string)

	for input := range in {
		wgIn := &sync.WaitGroup{}
		data, ok := input.(string)
		if !ok {
			log.Fatalf("can't convert %T to string", input)
		}
		inCh := make(chan hashNode)

		wgOut.Add(1)
		wgIn.Add(6)
		for i := 0; i < 6; i++ {
			go func(wg *sync.WaitGroup, i int, inp string, out chan hashNode) {
				defer wg.Done()
				out <- hashNode{i, DataSignerCrc32(fmt.Sprintf("%v%v", i, inp))}
			}(wgIn, i, data, inCh)
		}
		go func(wgInner *sync.WaitGroup, c chan hashNode) {
			defer close(c)
			wgInner.Wait()
		}(wgIn, inCh)

		go func(wg *sync.WaitGroup, in chan hashNode, out chan string) {
			data := map[int]string{}
			var dataKeys []int

			for o := range in {
				data[o.id] = o.value
				dataKeys = append(dataKeys, o.id)
			}
			sort.Ints(dataKeys)

			var results []string
			for i := range dataKeys {
				results = append(results, data[i])
			}

			out <- strings.Join(results, "")
			wg.Done()
		}(wgOut, inCh, outCh)
	}

	go func(wgOut *sync.WaitGroup, c chan string) {
		defer close(c)
		wgOut.Wait()
	}(wgOut, outCh)

	for hash := range outCh {
		out <- hash
	}
}

func CombineResults(in, out chan interface{}) {
	var result []string

	for input := range in {
		data, _ := input.(string)
		result = append(result, data)
	}

	sort.Strings(result)
	out <- strings.Join(result, "_")
}