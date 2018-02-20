package main

func main() {
	inputData := []int{0, 1, 1, 2, 3, 5, 8}
	hashSignJobs := []job{
		job(func(in, out chan interface{}) {
			for _, fibNum := range inputData {
				out <- fibNum
			}
		}),
		job(SingleHash),
		job(MultiHash),
		job(CombineResults),
		//job(func(in, out chan interface{}) {
		//	dataRaw := <-in
		//	data, ok := dataRaw.(string)
		//	if !ok {
		//		t.Error("cant convert result data to string")
		//	}
		//	testResult = data
		//}),
	}

	ExecutePipeline(hashSignJobs...)
}
