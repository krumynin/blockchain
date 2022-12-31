package network

import (
	"net"
	"strings"
)

// Listen address ipv4:port
func Listen(address string, handle func(Conn, *Package)) Listener {
	addr := strings.Split(address, ":")
	if len(addr) != 2 {
		return nil
	}
	listener, err := net.Listen("tcp", "0.0.0.0:"+addr[1])
	if err != nil {
		return nil
	}

	serve(listener, handle)

	return listener
}

func serve(listener Listener, handle func(Conn, *Package)) {
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			break
		}
		go handleConn(conn, handle)
	}
}

func handleConn(conn Conn, handle func(Conn, *Package)) {
	defer conn.Close()
	pack := readPackage(conn)
	if pack == nil {
		return
	}
	handle(conn, pack)
}
