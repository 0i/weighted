package weighted

import "slices"

// smoothWeighted is a wrapped weighted item.
type smoothWeighted[T comparable] struct {
	Item            T
	Weight          int
	CurrentWeight   int
	EffectiveWeight int
}

/*
SW (Smooth Weighted) is a struct that contains weighted items and provides methods to select a weighted item.
It is used for the smooth weighted round-robin balancing algorithm. This algorithm is implemented in Nginx:
https://github.com/phusion/nginx/commit/27e94984486058d73157038f7950a0a36ecc6e35.

Algorithm is as follows: on each peer selection we increase current_weight
of each eligible peer by its weight, select peer with greatest current_weight
and reduce its current_weight by total number of weight points distributed
among peers.

In case of { 5, 1, 1 } weights this gives the following sequence of
current_weight's: (a, a, b, a, c, a, a)

*/
type SW[T comparable] struct {
	items []*smoothWeighted[T]
	n     int
}

// func (w *smoothWeighted) fail() {
// 	w.EffectiveWeight -= w.Weight
// 	if w.EffectiveWeight < 0 {
// 		w.EffectiveWeight = 0
// 	}
// }

// Add a weighted server.
func (w *SW[T]) Add(item T, weight int) {
	weighted := &smoothWeighted[T]{Item: item, Weight: weight, EffectiveWeight: weight}
	w.items = append(w.items, weighted)
	w.n++
}

// RemoveAll removes all weighted items.
func (w *SW[T]) RemoveAll() {
	w.items = w.items[:0]
	w.n = 0
}

// Reset resets all current weights.
func (w *SW[T]) Reset() {
	for _, s := range w.items {
		s.EffectiveWeight = s.Weight
		s.CurrentWeight = 0
	}
}

// All returns all items.
func (w *SW[T]) All() map[interface{}]int {
	m := make(map[interface{}]int)
	for _, i := range w.items {
		m[i.Item] = i.Weight
	}
	return m
}

// Next returns next selected server.
func (w *SW[T]) Next(exclusions ...T) T {
	i := w.nextWeighted(exclusions...)
	if i == nil {
		var t T
		return t
	}
	return i.Item
}

// nextWeighted returns next selected weighted object.
func (w *SW[T]) nextWeighted(exclusions ...T) *smoothWeighted[T] {
	if w.n == 0 {
		return nil
	}

	if es := len(exclusions); es == 0 {
		if w.n == 1 {
			return w.items[0]
		}

		return nextSmoothWeighted[T](w.items)
	} else {
		items := make([]*smoothWeighted[T], 0, len(w.items)-es)
		for _, i := range w.items {
			if !slices.Contains(exclusions, i.Item) {
				items = append(items, i)
			}
		}

		return nextSmoothWeighted[T](items)
	}
}

// https://github.com/phusion/nginx/commit/27e94984486058d73157038f7950a0a36ecc6e35
func nextSmoothWeighted[T comparable](items []*smoothWeighted[T]) (best *smoothWeighted[T]) {
	total := 0

	for i := 0; i < len(items); i++ {
		w := items[i]

		if w == nil {
			continue
		}

		w.CurrentWeight += w.EffectiveWeight
		total += w.EffectiveWeight
		if w.EffectiveWeight < w.Weight {
			w.EffectiveWeight++
		}

		if best == nil || w.CurrentWeight > best.CurrentWeight {
			best = w
		}

	}

	if best == nil {
		return nil
	}

	best.CurrentWeight -= total
	return best
}
