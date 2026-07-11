package main

import "sync"

type History struct {
	mu sync.Mutex
	lines []string
	max int 
}

func NewHistory(max int) *History {
	return &History{max: max}
}

func (h *History) Add (line string){
	h.mu.Lock()
	defer h.mu.Unlock()
	h.lines = append(h.lines, line)
	if len(h.lines) > h.max {
		h.lines = h.lines[1:]
	}
}

func (h *History) All() []string {
	h.mu.Lock()
	defer h.mu.Unlock()
	out := make([]string, len(h.lines))
	copy (out, h.lines)
	return out
}

