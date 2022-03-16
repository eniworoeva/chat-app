package main

import (
	"bufio"
	"log"
	"net"
)

//this function is called everytime we encounter
//an error, and we have to terminate the program
func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var openConnections = make(map[net.Conn]bool)
var newConnection = make(chan net.Conn)
var deadConnection = make(chan net.Conn)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	logFatal(err)

	defer ln.Close()

	go func() {
		for {
			conn, err := ln.Accept()
			logFatal(err)

			openConnections[conn] = true
			newConnection <- conn
		}
	}()

	for {
		select {
		case conn := <-newConnection:
			go broadcastMessage(conn)
		case conn := <-deadConnection:
			for item := range openConnections {
				if item == conn {
					break
				}
			}
			delete(openConnections, conn)
		}
	}
}
func broadcastMessage(conn net.Conn) {
	for {
		reader := bufio.NewReader(conn)
		message, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		for item := range openConnections {
			if item != conn {
				item.Write([]byte(message))
			}
		}
	}
	deadConnection <- conn
}
