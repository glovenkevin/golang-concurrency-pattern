package pipeline

import (
	"math"
	"sync"
)

type IntegerReducerParam struct {
	Amount int
	Power  int
}

func LaunchIntegerReducer(in *IntegerReducerParam) int {
	out := generator(in.Amount)

	// distribute the job into 2 workers
	o1 := power(out, in.Power)
	o2 := power(out, in.Power)
	out = merge(o1, o2)

	out = sum(out)
	return <-out
}

func generator(in int) <-chan int {
	ret := make(chan int, in)

	go func() {
		for i := 0; i < in; i++ {
			ret <- i
		}
		close(ret)
	}()

	return ret
}

func power(in <-chan int, pow int) <-chan int {
	ret := make(chan int)
	go func() {
		for i := range in {
			ret <- int(math.Pow(float64(i), float64(pow)))
		}
		close(ret)
	}()
	return ret
}

func sum(in <-chan int) <-chan int {
	ret := make(chan int, 1)
	go func() {
		var sum int
		for v := range in {
			sum += v
		}
		ret <- sum
		close(ret)
	}()

	return ret
}

func merge(ii ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	output := func(c <-chan int) {
		defer wg.Done()
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
