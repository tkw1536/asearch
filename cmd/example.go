package main

import (
	"context"
	"fmt"

	"github.com/tkw1536/fsearch"
)

func main() {
	j := &fsearch.Joiner[int]{
		Cs: []chan int{
			makechanthing(5, 7, 9),
			makechanthing(1, 4, 10),
			makechanthing(2, 3, 8),
			makechanthing(5, 7, 9),
		},
	}

	all := j.Chan(context.Background())

	sliced := fsearch.Sorter[int]{
		C:    all,
		Less: func(a, b int) bool { return a > b },
	}

	fmt.Println(sliced.Slice(1, 2))
}

func makechanthing(values ...int) chan int {
	c := make(chan int)
	go func() {
		defer close(c)
		for _, v := range values {
			c <- v
		}
	}()
	return c
}
