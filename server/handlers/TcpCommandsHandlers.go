package handlers

import (
	"distrubuted-system/shared/types"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

func HandleConnection(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Read error:", err)
		return
	}

	var baseRequest map[string]interface{}
	if err := json.Unmarshal(buffer[:n], &baseRequest); err != nil {
		fmt.Println("Invalid request format:", err)
		return
	}

	commandType := int(baseRequest["command_type"].(float64))

	switch commandType {
	case types.ECHO:
		var req types.Request[types.EchoCommandData]
		sendResponse(conn, "Success", req.Data.Message)
		break
	case types.TIME:
		var _ types.Request[types.TimeCommandData]
		sendResponse(conn, "Success", time.Now().Local().Format(time.DateTime))
		break
	case types.CLOSE:
		sendResponse(conn, "Success", "Closing connection.")
		conn.Close()
	case types.UPLOAD:
		break
	case types.DOWNLOAD:
		break
	default:
		fmt.Println("Unknown command received")
		sendResponse(conn, "Error", "Unknown command received")
		break
	}
}

func sendResponse(conn net.Conn, status, message string) {
	response := types.Response{
		Status:  status,
		Message: message,
	}

	data, _ := json.Marshal(response)
	conn.Write(data)
}
