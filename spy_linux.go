package procspy

import (
	"bytes"
	"sync"
)

var bufPool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, 5000))
	},
}

type pnConnIter struct {
	transport string
	pn        *ProcNet
	buf       *bytes.Buffer
	procs     map[uint64]Proc
}

func (c *pnConnIter) Next() *Connection {
	n := c.pn.Next()
	if n == nil {
		// Done!
		bufPool.Put(c.buf)
		return nil
	}
	if proc, ok := c.procs[n.Inode]; ok {
		n.Proc = proc
	}
	n.Transport = c.transport
	return n
}

// cbConnections sets Connections()
var cbConnections = func(processes bool, state uint) (ConnIter, error) {
	// buffer for contents of /proc/<pid>/net/tcp
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()

	var procs map[uint64]Proc
	if processes {
		var err error
		if procs, err = walkProcPid(buf); err != nil {
			return nil, err
		}
	}

	if buf.Len() == 0 {
		readFile(procRoot+"/net/tcp", buf)
		readFile(procRoot+"/net/tcp6", buf)
	}

	return &pnConnIter{
		transport: TcpTransport,
		pn:        NewProcNet(buf.Bytes(), state),
		buf:       buf,
		procs:     procs,
	}, nil
}
