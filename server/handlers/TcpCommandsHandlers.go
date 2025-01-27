package handlers

import (
	"distrubuted-system/shared/types"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
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
		req.Data.Message = baseRequest["data"].(map[string]interface{})["message"].(string)
		sendResponse(conn, "Success", req.Data.Message)

	case types.TIME:
		sendResponse(conn, "Success", time.Now().Local().Format(time.DateTime))
		break

	case types.CLOSE:
		sendResponse(conn, "Success", "Closing connection.")
		conn.Close()

	case types.UPLOAD:
		var req types.Request[types.UploadCommandData]
		req.Data.FileName = baseRequest["data"].(map[string]interface{})["file_name"].(string)
		req.Data.Status = baseRequest["data"].(map[string]interface{})["status"].(string)
		fmt.Println("Uploading file:", req.Data.FileName)
		sendResponse(conn, "Connected", "Загрузка файла началась..")
		sendResponse(conn, "Success", req.Data.FileName)
		sendFileResponse(conn, req.Data.FileName)

	case types.DOWNLOAD:
		var req types.Request[types.DownloadCommandData]
		req.Data.FileName = baseRequest["data"].(map[string]interface{})["file_name"].(string)
		sendFile(conn, req.Data.FileName)
		sendResponse(conn, "Success", req.Data.FileName)

	default:
		fmt.Println("Unknown command received")
		sendResponse(conn, "Error", "Unknown command received")
	}
}

func sendResponse(conn net.Conn, status, message string) {
	response := types.Response{
		Status:  status,
		Message: message,
	}

	data, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error marshaling response:", err)
		return
	}

	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("Error sending response:", err)
	}
}

func sendFileResponse(conn net.Conn, fileName string) {

	outFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Ошибка при создании файла:", err)
		return
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, conn)
	if err != nil {
		fmt.Println("Ошибка при записи в файл:", err)
	}

	fmt.Println("Файл успешно получен:", fileName)

}

func sendFile(conn net.Conn, fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
	}
	defer file.Close()

	_, err = conn.Write([]byte(fileName))
	if err != nil {
		fmt.Println("Ошибка при отправке имени файла:", err)
		return
	}

	_, err = io.Copy(conn, file)
	if err != nil {
		fmt.Println("Ошибка при отправке файла:", err)
	}

	fmt.Println("Файл успешно отправлен:", fileName)
}
