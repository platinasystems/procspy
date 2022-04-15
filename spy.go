// Package procspy lists TCP connections, and optionally tries to find the
// owning processes. Works on Linux (via /proc) and Darwin (via `lsof -i` and
// `netstat`). You'll need root to use Processes().
package procspy

import (
	"fmt"
	"net"
)

const (
	TcpEstablished = 1  // according to /include/net/tcp_states.h
	TcpListen      = 10 // according to /include/net/tcp_states.h
	TcpTransport   = "tcp"
)

// Connection is a (TCP) connection. The Proc struct might not be filled in.
type Connection struct {
	Transport     string
	LocalAddress  net.IP
	LocalPort     uint16
	RemoteAddress net.IP
	RemotePort    uint16
	Inode         uint64
	Proc
}

// Connection is a (TCP) connection. The Proc struct might not be filled in.
type ConnectionImmutable struct {
	Transport     string
	LocalAddress  string
	LocalPort     uint16
	RemoteAddress string
	RemotePort    uint16
	Inode         uint64
	Proc
}

func (c *Connection) Immutable() ConnectionImmutable {
	return ConnectionImmutable{
		Transport:     c.Transport,
		LocalAddress:  fmt.Sprintf("%v", c.LocalAddress),
		LocalPort:     c.LocalPort,
		RemoteAddress: fmt.Sprintf("%v", c.RemoteAddress),
		RemotePort:    c.RemotePort,
		Inode:         c.Inode,
		Proc:          c.Proc,
	}
}

// Proc is a single process with PID and process name.
type Proc struct {
	PID        uint
	Name       string
	NsNetInode uint64
}

// ConnIter is returned by Connections().
type ConnIter interface {
	Next() *Connection
}

// Connections returns all established (TCP) connections.  If processes is
// false we'll just list all TCP connections, and there is no need to be root.
// If processes is true it'll additionally try to lookup the process owning the
// connection, filling in the Proc field. You will need to run this as root to
// find all processes.
func Connections(processes bool, state uint) (ConnIter, error) {
	return cbConnections(processes, state)
}
