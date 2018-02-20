package main

import (
	"sync"
	"fmt"
)

type SingleHashResult struct {
	Hash string
	Number string
}

func ExecutePipeline(hashSignJobs ...job) {
	wg := &sync.WaitGroup{}

	in := make(chan interface{})
	out := make(chan interface{})

	wg.Add(len(hashSignJobs))
	for _, jobItem := range hashSignJobs {
		go func(jobFunc job, in chan interface{}, out chan interface{}, wg *sync.WaitGroup) {
			jobFunc(in, out)
			defer wg.Done()
		}(jobItem, in, out, wg)
	}

	defer wg.Wait()
}

func SingleHash(in chan interface{}, out chan interface{}) {
	for data := range out {
		number := fmt.Sprintf("%v", data)
		result := DataSignerCrc32(number)+ "~" + DataSignerCrc32(DataSignerMd5(number))
		in <- SingleHashResult{Hash: result, Number: number}
	}
}

func MultiHash(in chan interface{}, out chan interface{})  {
	for th := range in {
		hashResult := (th).(SingleHashResult)
		fmt.Printf("number %v, hash %v ", hashResult.Number, hashResult.Hash)
		fmt.Printf("MultiHash %v \n", DataSignerCrc32(hashResult.Hash + hashResult.Number))

		//out <- DataSignerCrc32(hashResult.Hash + hashResult.Number)
	}
}

func CombineResults(in, out chan interface{}){

	//for multiHashResut := range out {
	//	fmt.Println("out", multiHashResut)
	//}
}


/*
func SingleHash(in chan interface{}, out chan interface{}) {

	wg := &sync.WaitGroup{}

	tmpOut := make(chan string)

	for data := range in {

		wg.Add(1)

		raw, _ := data.(int)

		data := strconv.Itoa(raw)
		md5 := DataSignerMd5(data)

		go func(wg *sync.WaitGroup, out chan string, d string, md5 string) {
			hash := make(chan string)
			hashMd5 := make(chan string)

			go func(out chan string, input string) {
				out <- DataSignerCrc32(input)
			}(hash, d)

			go func(out chan string, input string) {
				out <- DataSignerCrc32(input)
			}(hashMd5, md5)

			result := fmt.Sprintf("%v~%v", <-hash, <-hashMd5)
			out <- result

			wg.Done()

		}(wg, tmpOut, data, md5)
	}

	go func(wg *sync.WaitGroup, c chan string) {
		defer close(c)
		wg.Wait()
	}(wg, tmpOut)


	for hash := range tmpOut {
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

	wgInCount := 6

	for input := range in {
		wgOut.Add(1)

		wgIn := &sync.WaitGroup{}
		data, _ := input.(string)

		inCh := make(chan hashNode)

		wgIn.Add(wgInCount)

		for i := 0; i < wgInCount; i++ {
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

*/