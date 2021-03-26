package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func getData(connection net.Conn, channel chan string) {
	defer connection.Close()
	buffer := make([]byte, 1024)
	for {
		bytes, err := connection.Read(buffer)
		if err != nil {
			fmt.Println("error reading connection: ", err.Error())
		}
		if bytes > 0 {
			channel <- string(buffer[:])
		}
	}
}
func main() {
	connection := make([]string, len(os.Args[1:]))
	channel := make(chan string, len(os.Args[1:]))
	for num, part := range os.Args[1:] {
		split := strings.Split(part, "=")
		connection[num] = split[1]
	}
	for _, con := range connection {
		conn, err := net.Dial("tcp", con)
		if err != nil {
			fmt.Println("error connecting:", err.Error())
		}
		go getData(conn, channel)
	}
	for {
		for line := range channel {
			fmt.Printf("\r%v", line)
		}
	}
	close(channel)
}
