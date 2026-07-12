package main

import (
	"net"
	"sync"
)

// creating the clients struct
type Clients struct {
	mu    sync.Mutex          // lock
	conns map[net.Conn]string // conns and their names from the config.
}

// constructor for a client
func NewClients() *Clients {
	return &Clients{conns: make(map[net.Conn]string)}
}

// add a client
func (c *Clients) Add(conn net.Conn, name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.conns[conn] = name // take the name from config when the user registers for the first time. the same username everytime
}

// remove a client meaning close the CLI
func (c *Clients) Remove(conn net.Conn) {
	// we dont need the name to remove connection
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.conns, conn)
	conn.Close()
}

// broadcast a message to all the clients in the conn map
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

// count number of clients -> clients are basically local terminals/CLIs
func (c *Clients) Count() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.conns)
}
