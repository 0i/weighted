package weighted

import "math/rand"

type weighted3 struct {
	Server interface{}
	Weight int
}

type W3 struct {
	servers []*weighted3
	totalW  int
}

// Add a weighted server.
func (w *W3) Add(server interface{}, weight int) {
	if weight <= 0 {
		return
	}

	weighted := &weighted3{Server: server, Weight: w.totalW + weight}
	w.totalW += weight

	w.servers = append(w.servers, weighted)
}

// RemoveAll removes all weighted servers.
func (w *W3) RemoveAll() {
	w.servers = nil
	w.totalW = 0
}

//Reset resets all current weights.
func (w *W3) Reset() {
	w.RemoveAll()
}

// Next returns next selected server.
func (w *W3) Next() interface{} {
	if w.totalW == 0 {
		return nil
	}

	randint := rand.Intn(w.totalW)
	for _, weighted := range w.servers {
		if randint < weighted.Weight {
			return weighted.Server
		}
	}

	return nil
}
