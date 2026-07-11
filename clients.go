package main

import (
	"net"
	"sync"
)

type Clients struct {
	mu    sync.Mutex
	conns map[net.Conn]bool
}

func NewClients() *Clients {
	return &Clients{conns: make(map[net.Conn]bool)}
}

func (c *Clients) Add(conn net.Conn) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.conns[conn] = true
}

func (c *Clients) Remove(conn net.Conn) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.conns, conn)
	conn.Close()
}

func (c *Clients) Broadcast(line string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for conn := range c.conns {
		_, err := conn.Write([]byte(line + "\n"))
		if err != nil {
			delete(c.conns, conn)
			conn.Close()
		}
	}
}

func (c *Clients) Count() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.conns)
}
