package asearch

import (
	"context"
	"sync"

	"golang.org/x/exp/constraints"
)

// Joiner implements the joining of underlying channels C.
//
// Joiner is safe to be used concurrently, unless Cs is modified.
type Joiner[T constraints.Ordered] struct {
	m  sync.Mutex
	wg sync.WaitGroup

	// A slice of channels, which may or may not be closed.
	//
	// Elements are assumed to be returned in sorted order for each channel.
	// The same element may be returned by multiple channels.
	Cs []chan T

	ok   []bool
	data []T
}

// SendTo returns a new channel that behaves like Next().
// the channel is closed once no more elements are available, or the context is closed.
func (j *Joiner[T]) Chan(ctx context.Context) chan T {
	c := make(chan T)
	go func() {
		defer close(c)
		for {
			ok, _, next := j.Next(ctx)
			if !ok {
				return
			}
			c <- next
		}
	}()
	return c
}

// Fill ensures that the data buffer is full
// That is, it performs a receive on all channels where no buffered data is available.
func (j *Joiner[T]) fill(ctx context.Context) {
	// initialize ok and data buffers, if one of them is the wrong size
	if len(j.ok) != len(j.Cs) || len(j.data) != len(j.Cs) {
		j.ok = make([]bool, len(j.Cs))
		j.data = make([]T, len(j.Cs))
	}

	// fill all the unfilled elements
	// (typically only one element but that's ok)
	for i, ok := range j.ok {
		if ok {
			continue
		}

		j.wg.Add(1)
		go func(i int) {
			defer j.wg.Done()
			select {
			case j.data[i], j.ok[i] = <-j.Cs[i]:
			case <-ctx.Done():
			}

		}(i)
	}

	// wait for everything to be done
	j.wg.Wait()
}

// Next returns the next smallest element received from any of the underlying channels.
// When such an element exists, returns ok = True, source as the id of the channel it came from, and the value.
//
// When all channels have closed, returns ok = False
func (j *Joiner[T]) Next(ctx context.Context) (ok bool, source int, value T) {
	j.m.Lock()
	defer j.m.Unlock()

	j.fill(ctx) // fill the buffer

	if ctx.Err() != nil { // if the context closed, return!
		return
	}

	// check the first filled element
	first := -1
	for i, ok := range j.ok {
		if ok {
			first = i
			break
		}
	}

	// no element is filled
	if first == -1 {
		return
	}

	// pick the smallest element
	smallest := first
	value = j.data[smallest]
	for i := first + 1; i < len(j.data); i++ {
		if !j.ok[i] { // element is invalid!
			continue
		}
		if j.data[i] < value {
			smallest = i
			value = j.data[i]
		}
	}

	// invalidate the element in the buffer
	var zero T
	j.ok[smallest] = false
	j.data[smallest] = zero

	// and return it!
	return true, smallest, value
}
