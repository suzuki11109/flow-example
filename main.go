package main

import (
	"fmt"
	"sync"
)

func gen(start int, end int) <-chan int {
	out := make(chan int)
	go func() {
		for i := start; i <= end; i++ {
			out <- i
		}
		close(out)
	}()
	return out
}

func sq(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func even(in <-chan int) <-chan int {
	return filter(in, func(n int) bool { return (n%2 == 0) })
}

func odd(in <-chan int) <-chan int {
	return filter(in, func(n int) bool { return (n%2 != 0) })
}

func filter(in <-chan int, fn func(int) bool) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			if fn(n) {
				out <- n
			}
		}
		close(out)
	}()
	return out
}

func merge(ins ...<-chan int) <-chan int {
	out := make(chan int)
	go func() {
		var wg sync.WaitGroup
		wg.Add(len(ins))

		for _, in := range ins {
			go func(inn <-chan int) {
				for n := range inn {
					out <- n
				}
				wg.Done()
				//Done
			}(in)
		}

		wg.Wait()
		close(out)
	}()
	return out
}

func broadcast(in <-chan int, out ...chan<- int) {
	go func() {
		var wg sync.WaitGroup

		for n := range in {
			wg.Add(len(out))

			for _, och := range out {
				go func(och chan<- int, n int) {
					och <- n
					wg.Done()
				}(och, n)
			}
		}
		wg.Wait()
		for _, och := range out {
			close(och)
		}
	}()
}

func dis(in <-chan int, out ...chan<- int) {
	go func() {
		var wg sync.WaitGroup

		for n := range in {
			wg.Add(len(out))
			next := make(chan struct{})
			for _, och := range out {
				go func(o chan<- int, n int, next chan struct{}) {
					select {
					case <-next:
					case o <- n:
						next <- struct{}{}
					}
					wg.Done()
				}(och, n, next)
			}
		}
		wg.Wait()
		for _, och := range out {
			close(och)
		}
	}()
}

func main() {
	// out := gen(1, 10)
	// out2 := gen(50, 60)
	// m := merge(out, out2)

	o1 := make(chan int)
	o2 := make(chan int)
	dis(gen(1, 3), o1, o2)
	m := merge(o1, sq(o2))

	for n := range m {
		fmt.Println(n)
	}
}
