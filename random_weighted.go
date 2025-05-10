package weighted

import (
	"math/rand/v2"
	"time"
)

// randWeighted is a wrapped weighted item that is used to implement weighted random algorithm.
type randWeighted[T comparable] struct {
	Item   T
	Weight int
}

// RandW is a struct that contains weighted items implement weighted random algorithm.
type RandW[T comparable] struct {
	items        []*randWeighted[T]
	n            int
	sumOfWeights int
	r            *rand.Rand
}

// NewRandW creates a new RandW with a random object.
func NewRandW[T comparable]() *RandW[T] {
	rw := &RandW[T]{}
	rw.Reset()
	return rw
}

// Next returns next selected item.
func (rw *RandW[T]) Next() (item T) {
	if rw.n == 0 {
		var t T
		return t
	}
	if rw.sumOfWeights <= 0 {
		var t T
		return t
	}
	randomWeight := rw.r.IntN(rw.sumOfWeights) + 1
	for _, item := range rw.items {
		randomWeight = randomWeight - item.Weight
		if randomWeight <= 0 {
			return item.Item
		}
	}

	return rw.items[len(rw.items)-1].Item
}

// Add adds a weighted item for selection.
func (rw *RandW[T]) Add(item T, weight int) {
	rItem := &randWeighted[T]{Item: item, Weight: weight}
	rw.items = append(rw.items, rItem)
	rw.sumOfWeights += weight
	rw.n++
}

// All returns all items.
func (rw *RandW[T]) All() map[T]int {
	m := make(map[T]int)
	for _, i := range rw.items {
		m[i.Item] = i.Weight
	}
	return m
}

// RemoveAll removes all weighted items.
func (rw *RandW[T]) RemoveAll() {
	rw.items = make([]*randWeighted[T], 0)
	rw.n = 0
	rw.sumOfWeights = 0
	rw.Reset()
}

// Reset resets the balancing algorithm.
func (rw *RandW[T]) Reset() {
	seed := uint64(time.Now().UnixNano())
	rw.r = rand.New(rand.NewPCG(seed, seed+1))
}
