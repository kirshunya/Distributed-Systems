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

func SendUploadRequest(conn net.Conn) {
	var fileName string
	fmt.Print("Введите имя файла: ")
	fmt.Scan(&fileName)

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Ошибка при получении информации о файле:", err)
		return
	}
	fileSize := fileInfo.Size()

	request := types.Request[types.UploadCommandData]{
		CommandType: types.UPLOAD,
		Data: types.UploadCommandData{
			FileName: fileName,
			FileSize: fileSize,
			Status:   "Sending file",
		},
	}
	sendRequest(conn, request)

	startTime := time.Now()
	_, err = io.Copy(conn, file)
	if err != nil {
		fmt.Println("Ошибка при отправке файла:", err)
		return
	}
	elapsedTime := time.Since(startTime).Seconds()
	bitrate := float64(fileSize) / elapsedTime / 1024

	fmt.Printf("Файл успешно отправлен: %s, Битрейт: %.2f Кб/с\n", fileName, bitrate)
}

func SendDownloadRequest(conn net.Conn) {
	fileName := "text.txt"
	request := types.Request[types.DownloadCommandData]{
		CommandType: types.DOWNLOAD,
		Data: types.DownloadCommandData{
			FileName: fileName,
			Status:   "File request to download",
		},
	}
	sendRequest(conn, request)

	outFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Ошибка при создании файла:", err)
		return
	}
	defer outFile.Close()

	fileCon, _ := net.Dial("tcp", "172.20.10.3:9090")
	fileCon.SetReadDeadline(time.Now().Add(keepAliveTTL * time.Second))
	fmt.Println(fileName)

	_, err = io.Copy(outFile, fileCon)
	if err != nil {
		fmt.Println("Ошибка при записи в файл:", err)
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

	fmt.Printf("Файл успешно получен: %s, %s", fileName, response)
}
