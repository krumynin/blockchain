package network

import (
	"net"
	"strings"
)

type Listener net.Listener
type Conn net.Conn

// Package Пакет, который передает пользователь (Команда + данные)
type Package struct {
	Command int
	Data    string
}

const (
	EndBytes    = "\000\005\007\001\001\007\005\000"
	WaitTime    = 5
	DataMaxSize = 2 << 20 // (2^20)*2 = 2MiB = (1024^2) * 2 байт
	BufferSize  = 4 << 10 // (2^10)*4 = 4KiB = (1024) * 4 байт
)

func readPackage(conn net.Conn) *Package {
	var (
		dataSize = uint64(0)
		buffer   = make([]byte, BufferSize)
		data     string
	)
	for {
		length, err := conn.Read(buffer)
		if err != nil {
			return nil
		}
		dataSize += uint64(length)
		if dataSize > DataMaxSize {
			return nil
		}
		data += string(buffer[:length])
		if strings.Contains(data, EndBytes) {
			data = strings.Split(data, EndBytes)[0]
			break
		}
	}

	return deserializePackage(data)
}
