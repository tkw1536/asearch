package main

import (
	"fmt"

	"github.com/tkw1536/fsearch"
)

func main() {
	/*
		j := &fsearch.Joiner[int]{
			Cs: []chan int{
				makechanthing(5, 7, 9),
				makechanthing(1, 4, 10),
				makechanthing(2, 3, 8),
				makechanthing(5, 7, 9),
			},
		}

		all := j.Chan(context.Background())
	*/
	grouper := &fsearch.GroupBuffer[int, int]{
		C:     makechanthing(1, 2, 3, 4, 5, 6, 7, 8, 9, 10),
		Group: func(i int) int { return i % 10 },
		Less:  func(a, b int) bool { return a < b },
	}

	fmt.Println(grouper.Slice(1, 2))
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
