package handlers

import (
	"distrubuted-system/shared/types"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
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

func SendFileRequest(conn net.Conn) {
	var fileName string
	fmt.Print("Введите имя файла: ")
	fmt.Scan(&fileName)
	request := types.Request[types.UploadCommandData]{
		CommandType: types.UPLOAD,
		Data: types.UploadCommandData{
			FileName: fileName,
			Status:   "Sending file",
		},
	}
	sendRequest(conn, request)

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
	}
	defer file.Close()

	// Отправляем имя файла
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
	request = types.Request[types.UploadCommandData]{
		CommandType: types.UPLOAD,
		Data: types.UploadCommandData{
			FileName: fileName,
			Status:   "File sent",
		},
	}
	sendRequest(conn, request)
	//sendRequest(conn, request)
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

	_, err = io.Copy(outFile, conn)
	if err != nil {
		fmt.Println("Ошибка при записи в файл:", err)
	}

	fmt.Println("Файл успешно получен:", fileName)

	request = types.Request[types.DownloadCommandData]{
		CommandType: types.DOWNLOAD,
		Data: types.DownloadCommandData{
			FileName: fileName,
			Status:   "File accepted",
		},
	}

}
