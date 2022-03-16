package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	connection, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		logFatal(err)
	}

	defer connection.Close()
	fmt.Println("Enter your username")

	reader := bufio.NewReader(os.Stdin)
	username, err := reader.ReadString('\n')

	logFatal(err)
	username = strings.Trim(username, "\r\n")

	welcomeMsg := fmt.Sprintf("Welcome %s, to the chat, say hi to your friends.", username)

	fmt.Println(welcomeMsg)
	go read(connection)
	write(connection, username)
}
func read(connection net.Conn) {
	for {
		reader := bufio.NewReader(connection)
		message, err := reader.ReadString('\n')
		if err == io.EOF {
			connection.Close()
			fmt.Println("Connection Closed.")
			os.Exit(0)
		}
		fmt.Println(message)
		fmt.Println("--------------------------")
	}
}
func write(connection net.Conn, username string) {
	for {
		reader := bufio.NewReader(os.Stdin)
		message, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		message = fmt.Sprintf("%s:- %s\n", username, strings.Trim(message, " \r\n"))
		connection.Write([]byte(message))
	}
}
