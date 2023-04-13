package pipeline

import (
	"math"
	"sync"
)

func LaunchIntegerReducerWithCancellation(in *IntegerReducerParam, doneFlag <-chan struct{}) int {
	out := generator(in.Amount)

	// distribute the job into 2 workers
	o1 := powerWithDone(doneFlag, out, in.Power)
	o2 := powerWithDone(doneFlag, out, in.Power)
	out = mergeWithDone(doneFlag, o1, o2)

	out = sum(out)
	return <-out
}

func powerWithDone(done <-chan struct{}, in <-chan int, pow int) <-chan int {
	ret := make(chan int)
	go func() {
		select {
		case <-done:
			return
		default:
		}

		for i := range in {
			ret <- int(math.Pow(float64(i), float64(pow)))
		}
		close(ret)
	}()
	return ret
}

func mergeWithDone(done <-chan struct{}, ii ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	output := func(c <-chan int) {
		defer wg.Done()
		select {
		case <-done:
			return
		default:
		}

		for n := range c {
			out <- n
		}
	}

	wg.Add(len(ii))
	for _, i := range ii {
		go output(i)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
