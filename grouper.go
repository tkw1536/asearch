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
	if len(gb.data) < gb.size {
		new := gb.data[:insert+1]
		new = append(new, Group[T, U]{
			Group:    group,
			Elements: []T{data},
		})
		gb.data = append(new, gb.data[insert:]...)
		return
	}

	// replace and remove the last element
	new := gb.data[:insert+1]
	new = append(new, Group[T, U]{
		Group:    group,
		Elements: []T{data},
	})
	gb.data = append(new, gb.data[insert:len(gb.data)]...)

}
