package main

import (
	"distrubuted-system/client/handlers"
	"fmt"
	"net"
)

func main() {
	ipAddress := "localhost:8080" // type here ip address of host machine
	conn, err := net.Dial("tcp", ipAddress)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	for {
		fmt.Println("Choose a command:")
		fmt.Println("1. ECHO")
		fmt.Println("2. TIME")
		fmt.Println("3. UPLOAD")
		fmt.Println("4. DOWNLOAD")
		fmt.Println("5. CLOSE")

		var choice int
		_, err := fmt.Scan(&choice)
		if err != nil {
			fmt.Println("Invalid input:", err)
			continue
		}

		switch choice {
		case 1:
			handlers.SendEchoRequest(conn)
			break
		case 2:
			handlers.SendTimeRequest(conn)
			break
		case 3:
			handlers.SendFileResponse(conn)
			break
		case 4:
			break
		case 5:
			fmt.Println("Closing connection.")
			break
		default:
			fmt.Println("Invalid choice, please try again.")
			break
		}
	}
}
