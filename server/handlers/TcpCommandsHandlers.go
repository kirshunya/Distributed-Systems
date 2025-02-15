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

const (
	keepAliveTTL                 = 5
	fileTransferPort      string = ":9090"
	fileTransferChunkSize        = 1024
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

	case types.CLOSE:
		sendResponse(conn, "Success", "Closing connection.")
		conn.Close()

	case types.UPLOAD:
		var req types.Request[types.UploadCommandData]
		sendResponse(conn, "Connected", "Загрузка файла началась..")
		req.Data.FileName = baseRequest["data"].(map[string]interface{})["file_name"].(string)
		req.Data.Status = baseRequest["data"].(map[string]interface{})["status"].(string)
		fileSizeFloat := baseRequest["data"].(map[string]interface{})["file_size"].(float64)
		req.Data.FileSize = int64(fileSizeFloat)
		fmt.Println("Uploading file:", req.Data.FileName)
		sendFileResponse(conn, req.Data.FileName, req.Data.FileSize)

	case types.DOWNLOAD:
		var req types.Request[types.DownloadCommandData]
		req.Data.FileName = baseRequest["data"].(map[string]interface{})["file_name"].(string)
		sendResponse(conn, "Connected", "Загрузка файла началась..")
		startTime := time.Now()
		sendFile(conn, req.Data.FileName)
		fileInfo, _ := os.Stat(req.Data.FileName)
		elapsedTime := time.Since(startTime).Seconds()
		bitrate := float64(fileInfo.Size()*8) / elapsedTime / 1024
		sendResponse(conn, "Success", fmt.Sprintf("%s, Битрейт: %.2f Кб/с", req.Data.FileName, bitrate))

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

func sendFileResponse(conn net.Conn, fileName string, fileSize int64) {
	outFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Ошибка при создании файла:", err)
		return
	}
	defer outFile.Close()

	startTime := time.Now()
	_, err = io.CopyN(outFile, conn, fileSize)
	if err != nil {
		fmt.Println("Ошибка при записи в файл:", err)
		return
	}
	elapsedTime := time.Since(startTime).Seconds()
	bitrate := float64(fileSize) / elapsedTime / 1024 // в Кб/c

	fmt.Printf("Файл успешно получен: %s, Битрейт: %.2f Кб/с\n", fileName, bitrate)
	sendResponse(conn, "Success", fileName)
}

func sendFile(conn net.Conn, fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
	}
	defer file.Close()

	fileListener, err := net.Listen("tcp", fileTransferPort)
	if err != nil {
		panic(err)
	}

	defer fileListener.Close()

	fileTransferConn, err := fileListener.Accept()
	fileTransferConn.SetReadDeadline(time.Now().Add(keepAliveTTL * time.Second))

	// Отправляем содержимое файла
	_, err = io.Copy(fileTransferConn, file)
	if err != nil {
		fmt.Println("Ошибка при отправке файла:", err)
		return
	}

	fmt.Println("Файл успешно отправлен:", fileName)
}
