package handlers

import (
	"distrubuted-system/shared/types"
	"encoding/json"
	"fmt"
	"net"
)

func SendTimeRequest(conn net.Conn) {
	request := types.Request[types.TimeCommandData]{
		CommandType: types.TIME,
	}

	sendRequest(conn, request)
}

func SendEchoRequest(conn net.Conn) {
	message := "Hello, server!"
	request := types.Request[types.EchoCommandData]{
		CommandType: types.ECHO,
		Data: types.EchoCommandData{
			Message: message,
		},
	}

	sendRequest(conn, request)
}

func sendRequest[T any](conn net.Conn, request types.Request[T]) {
	data, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error marshaling request:", err)
		return
	}

	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	var response types.Response
	err = json.Unmarshal(buffer[:n], &response)
	if err != nil {
		fmt.Println("Error unmarshaling response:", err)
		return
	}

	fmt.Printf("Server response: Status=%s, Message=%s\n", response.Status, response.Message)
}
