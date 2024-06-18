package main

import (
	"fmt"
	"strings"

	// Uncomment this block to pass the first stage
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	conn, err := l.Accept()
	defer conn.Close()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	b := make([]byte, 4048)
	n, err := conn.Read(b)
	if err != nil || n == 0 {
		fmt.Println("Error reading from connection")
		os.Exit(1)
	}
	lines := strings.Split(string(b), "\r\n")
	if len(lines) > 0 {
		reqLines := strings.Fields(lines[0])
		fmt.Println(reqLines)
		fmt.Println(len(reqLines))
		if len(reqLines) < 2 {
			fmt.Println("req line bad formatted")
			os.Exit(1)
		}
		if reqLines[0] == "GET" && reqLines[1] == "/" {
			conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		} else {
			conn.Write([]byte("HTTP/1.1 404 NOT FOUND\r\n\r\n"))
		}
	} else {
		fmt.Println("Bad formatted request")
		os.Exit(1)
	}
}
