package asearch

import (
	"golang.org/x/exp/slices"
)

// Sorter can reorder elements from C, and return a finite subset of results.
type Sorter[T any] struct {
	// C is the channel where initial results are delivered to
	C chan T

	// Less is the order of the output
	Less func(a, b T) bool
}

// Slice0 returns count number of elements from the slicer.
// Elements are the count smallest elements contained in the entire result set.
func (s Sorter[T]) Slice0(count int) []T {
	results := make([]T, 0, count)
	for i := 0; i < count; i++ {
		value, ok := <-s.C
		if !ok {
			break
		}

		// append the next element
		results = append(results, value)
	}

	slices.SortFunc(results, s.Less)

	// insert elements one-by-one, in sorted order
	// evicting elements from the buffer that are no longer needed
	for next := range s.C {
		for i, e := range results {
			if s.Less(next, e) {
				// element is smaller and should be inserted before the current index
				// so shift everything one-by-one, and insert the current element
				copy(results[i+1:], results[i:])
				results[i] = next
				break
			}
		}
	}

	return results
}
func (s Sorter[T]) Slice(skip, limit int) []T {
	return s.Slice0(skip + limit)[skip:]
}
