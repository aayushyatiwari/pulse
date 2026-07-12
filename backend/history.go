package main

import "sync"

type History struct {
	mu    sync.Mutex // mutex for locking shared memory
	lines []string   // actual memory
	max   int        // max limit for memory
}

// go doesn't have constructor
func NewHistory(max int) *History {
	return &History{max: max}
}

func (h *History) Add(lines string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.lines = append(h.lines, lines)
	if len(h.lines) > h.max {
		h.lines = h.lines[1:]
	}
}

func (h *History) All() []string {
	h.mu.Lock()
	defer h.mu.Unlock()

	out := make([]string, len(h.lines)) // make bytes for printing history
	copy(out, h.lines)
	return out
}
