package main

import (
	"distrubuted-system/server/handlers"
	"fmt"
	"net"
)

func main() {
	requestListener, err := net.Listen("tcp", ":8081")
	if err != nil {
		panic(err)
	}

	defer requestListener.Close()

	for {
		conn, err := requestListener.Accept()
		if err != nil {
			panic(err)
		}

		fmt.Println("New connection from", conn.RemoteAddr())

		go handlers.HandleConnection(conn)
	}
}
