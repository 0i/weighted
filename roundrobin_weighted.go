package weighted

import "slices"

// rrWeighted is a wrapped weighted item that is used to implement LVS weighted round robin algorithm.
type rrWeighted[T comparable] struct {
	Item   T
	Weight int
}

// RRW is a struct that contains weighted items implement LVS weighted round robin algorithm.
//
// http://kb.linuxvirtualitem.org/wiki/Weighted_Round-Robin_Scheduling
// http://zh.linuxvirtualitem.org/node/37
type RRW[T comparable] struct {
	items []*rrWeighted[T]
	n     int
	gcd   int
	maxW  int
	i     int
	cw    int
}

// Add a weighted item.
func (w *RRW[T]) Add(item T, weight int) {
	weighted := &rrWeighted[T]{Item: item, Weight: weight}
	if weight > 0 {
		if w.gcd == 0 {
			w.gcd = weight
			w.maxW = weight
			w.i = -1
			w.cw = 0
		} else {
			w.gcd = gcd(w.gcd, weight)
			if w.maxW < weight {
				w.maxW = weight
			}
		}
	}
	w.items = append(w.items, weighted)
	w.n++
}

// All returns all items.
func (w *RRW[T]) All() map[T]int {
	m := make(map[T]int)
	for _, i := range w.items {
		m[i.Item] = i.Weight
	}
	return m
}

// RemoveAll removes all weighted items.
func (w *RRW[T]) RemoveAll() {
	w.items = w.items[:0]
	w.n = 0
	w.gcd = 0
	w.maxW = 0
	w.i = -1
	w.cw = 0
}

// Reset resets all current weights.
func (w *RRW[T]) Reset() {
	w.i = -1
	w.cw = 0
}

// Next returns next selected item.
func (w *RRW[T]) Next(exclusions ...T) T {
	if w.n == 0 {
		var t T
		return t
	}

	if len(exclusions) == 0 {
		if w.n == 1 {
			return w.items[0].Item
		}

		for {
			w.i = (w.i + 1) % w.n
			if w.i == 0 {
				w.cw = w.cw - w.gcd
				if w.cw <= 0 {
					w.cw = w.maxW
					if w.cw == 0 {
						var t T
						return t
					}
				}
			}

			if w.items[w.i].Weight >= w.cw {
				return w.items[w.i].Item
			}
		}
	} else {
		w2 := &RRW[T]{}
		for _, i := range w.items {
			if !slices.Contains(exclusions, i.Item) {
				w2.Add(i.Item, i.Weight)
			}
		}

		return w2.Next()
	}
}

func gcd(x, y int) int {
	var t int
	for {
		t = (x % y)
		if t > 0 {
			x = y
			y = t
		} else {
			return y
		}
	}
}
