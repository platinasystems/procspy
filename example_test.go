package procspy_test

import (
	"fmt"

	"github.com/platinasystems/procspy"
)

func Example() {
	lookupProcesses := true
	cs, err := procspy.Connections(lookupProcesses, procspy.TcpEstablished)
	if err != nil {
		panic(err)
	}

	fmt.Printf("TCP Connections:\n")
	for c := cs.Next(); c != nil; c = cs.Next() {
		fmt.Printf(" - %v\n", c)
	}
}
