// Clock Server is a concurrent TCP server that periodically writes the time.
package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
	"time"
)

func handleConn(c net.Conn, channel chan string, zone string) {
	defer c.Close()
	for {
		time.Sleep(time.Second)
		_, err := io.WriteString(c, zone+"\t: "+<-channel)
		if err != nil {
			return // e.g., client disconnected
		}
	}
}
func clockWall(Location string, channelOut chan string) {
	for {
		local, err := time.LoadLocation(Location)
		t := time.Now()
		if err == nil {
			t = t.In(local)
		}
		channelOut <- t.Format("15:04:05\n")
	}
}
func main() {
	var port string
	flag.StringVar(&port, "port", "", "port to be opened")
	flag.Parse()
	listener, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		log.Fatal(err)
	}
	channel := make(chan string)
	TimeZone := os.Getenv("TZ")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		go clockWall(TimeZone, channel)
		go handleConn(conn, channel, TimeZone) // handle connections concurrently
	}
	close(channel)
}
