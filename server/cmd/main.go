package main

import (
	"distrubuted-system/server/handlers"
	"fmt"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		fmt.Println("New connection from", conn.RemoteAddr())

		go handlers.HandleConnection(conn)
	}

}
