package main

import (
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if len(os.Args) != 2 {
		log.Printf("Usage: %s <port>\n", os.Args[0])
		os.Exit(1)
	}
	port := os.Args[1]

	listener, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Println("Failed to open port on "+port+": ", err)
		os.Exit(1)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("\nReceived SIGINT or SIGTERM, terminating...")
		listener.Close()
		os.Exit(0)
	}()

	log.Println("Listening on 0.0.0.0:" + port)

	conn, err := listener.Accept()
	if err != nil {
		log.Println("Error accepting connection: ", err)
		os.Exit(1)
	}
	handleConnection(conn)
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	_, err := io.Copy(os.Stdout, conn)
	if err != nil {
		log.Println("Error reading from connection: ", err)
		return
	}
}
