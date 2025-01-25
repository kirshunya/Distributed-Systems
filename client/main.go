package main

import (
	"distrubuted-system/server/types"
	"encoding/json"
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	msg := types.Request[types.EchoCommandData]{
		CommandType: types.ECHO,
		Data: types.EchoCommandData{
			Message: "Hello, Server!",
		},
	}

	encoder := json.NewEncoder(conn)
	err = encoder.Encode(msg)
	if err != nil {
		fmt.Println("Ошибка при кодировании сообщения:", err)
		return
	}
}
