package network

import (
	"fmt"
	"net"
	"time"
)

const OK = "ok"

// Send Отправка пакета на адрес
func Send(address string, pack *Package) *Package {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil
	}
	defer conn.Close()
	_, err = conn.Write([]byte(serializePackage(pack) + EndBytes))
	if err != nil {
		return nil
	}

	var (
		result  = new(Package)
		channel = make(chan bool)
	)
	go func() {
		result = readPackage(conn)
		channel <- true
	}()

	select {
	case <-channel:
		fmt.Println("OK")
	case <-time.After(WaitTime * time.Second):
		fmt.Println("TIMEOUT")
	}

	return result
}
