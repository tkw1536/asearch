package fsearch

import (
	"sync"

	"golang.org/x/exp/constraints"
)

type GroupBuffer[T any, U constraints.Ordered] struct {
	m sync.Mutex

	C     chan T            // C is the channel where initial results are delivered to
	Group func(t T) U       // Group returns the group of U
	Less  func(a, b U) bool // Less is the order of groups

	data []Group[T, U]
	size int
}

func (gb *GroupBuffer[T, U]) Slice0(count int) []Group[T, U] {
	gb.m.Lock()
	defer gb.m.Unlock()

	gb.data = nil
	gb.size = count

	for c := range gb.C {
		gb.insert(gb.Group(c), c)
	}

	return gb.data
}

func (gb *GroupBuffer[T, U]) Slice(skip, limit int) []Group[T, U] {
	return gb.Slice0(skip + limit)[skip:]
}

type Group[T any, U constraints.Ordered] struct {
	Group    U
	Elements []T
}

func (gb *GroupBuffer[T, U]) insert(group U, data T) {
	// if we already have it in the group, insert into it
	for i, g := range gb.data {
		if g.Group == group {
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
