package handlers

import (
	"fmt"
	"net"
)

func HandleConnection(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Ошибка при чтении данных:", err)
			return
		}
		fmt.Printf("Получено от клиента: %s\n", string(buffer[:n]))

		_, err = conn.Write([]byte("Сообщение получено\n"))
		if err != nil {
			fmt.Println("Ошибка при отправке данных:", err)
			return
		}
	}
}
