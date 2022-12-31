package main

import (
	"fmt"
	"strings"
	"time"

	nt "blockchain/network"
)

const (
	ToUpper = iota + 1
	ToLower
)

const (
	Address = ":8080"
)

func main() {
	go nt.Listen(Address, handleServer)
	time.Sleep(time.Second)

	fmt.Println("start client commands")

	res := nt.Send(Address, &nt.Package{
		Command: ToUpper,
		Data:    "Hello World!",
	})
	fmt.Println(res.Data)

	res = nt.Send(Address, &nt.Package{
		Command: ToLower,
		Data:    "Hello World!",
	})
	fmt.Println(res.Data)
}

func handleServer(conn nt.Conn, pack *nt.Package) {
	nt.Handle(ToUpper, conn, pack, handleToUpper)
	nt.Handle(ToLower, conn, pack, handleToLower)
}

func handleToUpper(pack *nt.Package) string {
	return strings.ToUpper(pack.Data)
}

func handleToLower(pack *nt.Package) string {
	return strings.ToLower(pack.Data)
}
