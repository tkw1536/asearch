package asearch

import (
	"sync"
)

// GroupBy implements grouping a stream of individuals by their groups
// and returning an ordered slice
type GroupBy[I any, K any] struct {
	C    <-chan I             // C is a channel that returns individuals
	Kind func(individual I) K // Kind returns the kind of an individual

	Less  func(a, b K) bool // Less implements an ordering on kinds
	Equal func(a, b K) bool // Equal implements equality on K. When nil, assumed to be Less(a, b) == Less(b, a).

	m    sync.Mutex
	data []Group[I, K]
	size int
}

func (gb *GroupBy[I, U]) Slice0(count int) []Group[I, U] {
	gb.m.Lock()
	defer gb.m.Unlock()

	gb.data = nil
	gb.size = count

	for c := range gb.C {
		gb.insert(gb.Kind(c), c)
	}

	return gb.data
}

func (gb *GroupBy[T, U]) Slice(skip, limit int) []Group[T, U] {
	return gb.Slice0(skip + limit)[skip:]
}

type Group[T any, U any] struct {
	Group    U
	Elements []T
}

func (gb *GroupBy[T, U]) insert(group U, data T) {
	if gb.Equal == nil {
		gb.Equal = func(a, b U) bool {
			return gb.Less(a, b) == gb.Less(b, a)
		}
	}

	// if we already have it in the group, insert into it
	for i, g := range gb.data {
		if gb.Equal(g.Group, group) {
			gb.data[i].Elements = append(gb.data[i].Elements, data)
			return
		}
	}

	// find the index to insert the new group into!
	insert := len(gb.data)
	for i, g := range gb.data {
		if gb.Less(group, g.Group) {
			insert = i
			break
		}
	}

	// add a new element
	gb.data = Insert(gb.data, Group[T, U]{
		Group:    group,
		Elements: []T{data},
	}, insert, len(gb.data) < gb.size)

}

// Insert inserts element as index i
func Insert[T any](slice []T, elem T, i int, grow bool) []T {
	switch {
	case i < 0:
		panic("i < 0")
	case i < len(slice):
		s := slice
		if grow {
			s = append(s, elem)
		}
		copy(s[i+1:], s[i:])
		s[i] = elem
		return s
	case i == len(slice):
		if grow {
			return append(slice, elem)
		}
		return slice
	default:
		panic("i > len(slice)")
	}

}
