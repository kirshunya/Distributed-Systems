package handlers

import (
	"distrubuted-system/server/types"
	"encoding/json"
	"fmt"
	"net"
)

func HandleConnection(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Read error:", err)
		return
	}

	// Определяем команду
	var baseRequest map[string]interface{}
	if err := json.Unmarshal(buffer[:n], &baseRequest); err != nil {
		fmt.Println("Invalid request format:", err)
		return
	}

	commandType := int(baseRequest["command_type"].(float64)) // Приведение типа из JSON

	switch commandType {
	case types.ECHO:
		// Обрабатываем команду ECHO
		var req types.Request[types.EchoCommandData]
		if err := json.Unmarshal(buffer[:n], &req); err == nil {
			fmt.Println("ECHO received:", req.Data.Message)
			response := fmt.Sprintf("Server received: %s", req.Data.Message)
			conn.Write([]byte(response))
		}
	case types.UPLOAD:
		// Обрабатываем команду UPLOAD
		var req types.Request[types.UploadCommandData]
		if err := json.Unmarshal(buffer[:n], &req); err == nil {
			fmt.Println("UPLOAD received for file:", req.Data.FileName)
			conn.Write([]byte("File uploaded successfully"))
		}
	case types.CLOSE:
		fmt.Println("CLOSE command received. Closing connection.")
		conn.Write([]byte("Goodbye!"))
	default:
		fmt.Println("Unknown command received")
		conn.Write([]byte("Unknown command"))
	}
}
