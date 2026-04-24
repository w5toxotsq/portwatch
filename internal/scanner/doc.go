// Package scanner provides functionality for probing network ports
// to determine whether they are open or closed.
//
// Basic usage:
//
//	ports := []int{22, 80, 443, 8080}
//	states, err := scanner.Scan("localhost", "tcp", ports, time.Second)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, s := range scanner.OpenPorts(states) {
//		fmt.Printf("%s/%d is open\n", s.Protocol, s.Port)
//	}
//
// The scanner performs connection-based probing: a port is considered open
// if a connection can be established within the specified timeout. UDP probing
// is supported but results may be unreliable due to the connectionless nature
// of the protocol.
package scanner
