package workers_pool

import (
	"fmt"
	"log"
	"strings"
	"sync"
)

func NewStringRequest(s string, id int, wg *sync.WaitGroup) Request {
	return Request{
		Data: s,
		Handler: func(i interface{}) {
			defer wg.Done()

			s, ok := i.(string)
			if !ok {
				log.Fatal("Data is not string")
			}
			fmt.Println(s)
		},
	}
}

type PreffixSuffixWorker struct {
	Id     int
	Prefix string
	Suffix string
}

func (p *PreffixSuffixWorker) LaunchWorker(in chan Request) {
	p.appendPrefix(p.appendSuffix(p.uppercase(in)))
}

func (p *PreffixSuffixWorker) uppercase(in <-chan Request) <-chan Request {
	out := make(chan Request)
	go func() {
		for msg := range in {
			s, ok := msg.Data.(string)
			if !ok {
				msg.Handler(nil)
				continue
			}
			msg.Data = strings.ToUpper(s)
			out <- msg
		}
		close(out)
	}()

	return out
}

func (p *PreffixSuffixWorker) appendSuffix(in <-chan Request) <-chan Request {
	out := make(chan Request)
	go func() {
		for msg := range in {
			uppercaseString, ok := msg.Data.(string)
			if !ok {
				msg.Handler(nil)
				continue
			}
			msg.Data = fmt.Sprintf("%s-%s", uppercaseString, p.Suffix)
			out <- msg
		}
		close(out)
	}()

	return out
}

func (p *PreffixSuffixWorker) appendPrefix(in <-chan Request) {
	go func() {
		for msg := range in {
			uppercaseStringWithSuffix, ok := msg.Data.(string)
			if !ok {
				msg.Handler(nil)
				continue
			}
			msg.Handler(fmt.Sprintf("%s-%s", p.Prefix, uppercaseStringWithSuffix))
		}
	}()
}
