package main

import (
	"fmt"
	"go-concurency-pattern/pipeline"
	"go-concurency-pattern/workers_pool"
	"sync"
)

func main() {

	// Pipeline Pattern
	// Simple Source: https://go.dev/blog/pipelines
	a, p := 10, 2
	res := pipeline.LaunchIntegerReducer(&pipeline.IntegerReducerParam{
		Amount: a,
		Power:  p,
	})
	fmt.Printf("Total amount of %d powered to %d is %d \n\n", a, p, res)

	doneFlags := make(chan struct{}, 4)
	resCancel := pipeline.LaunchIntegerReducerWithCancellation(&pipeline.IntegerReducerParam{
		Amount: a,
		Power:  p,
	}, doneFlags)

	for i := 0; i < 4; i++ {
		doneFlags <- struct{}{}
	}
	close(doneFlags)
	fmt.Printf("Total amount of %d powered to %d is %d \n\n", a, p, resCancel)

	// Workers Pool Pattern
	var buffSize = 100
	var workerSize = 5
	dispatcher := workers_pool.NewDispatcher(buffSize)
	for i := 0; i < workerSize; i++ {
		wl := &workers_pool.PreffixSuffixWorker{
			Id:     i,
			Prefix: fmt.Sprintf("Worker ID: %d ", i),
			Suffix: "KA",
		}
		dispatcher.LaunchWorker(wl)
	}

	var reqSize = 10
	var wg sync.WaitGroup
	wg.Add(reqSize)
	for i := 0; i < reqSize; i++ {
		req := workers_pool.NewStringRequest("aaa", i, &wg)
		dispatcher.MakeRequest(req)
	}
	dispatcher.Stop()
	wg.Wait()
}
