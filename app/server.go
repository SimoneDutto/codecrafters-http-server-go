package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"

	// Uncomment this block to pass the first stage
	"net"
	"os"
)

type StatusLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

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
	httpRequestParts := strings.Split(string(b), "\r\n")
	if len(httpRequestParts) > 0 {
		statusLine, err := parseStatusLine(httpRequestParts[0])
		header := parseHeader(httpRequestParts[1:])
		_ = header
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if statusLine.Method == "GET" && statusLine.RequestTarget == "/" {
			response := http.Response{
				Status:     "200 OK",
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
			}
			err := response.Write(conn)
			if err != nil {
				fmt.Println("Error during write", err)
			}
		} else if statusLine.Method == "GET" && strings.HasPrefix(statusLine.RequestTarget, "/echo") {
			echoString := path.Base(statusLine.RequestTarget)
			stringReader := strings.NewReader(echoString)
			stringReadCloser := io.NopCloser(stringReader)
			h := http.Header{}
			h.Add("Content-Type", "text/plain")
			h.Add("Content-Length", fmt.Sprint(len(echoString)))
			response := http.Response{
				Status:     "200 OK",
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				Header:     h,
				Body:       stringReadCloser,
			}
			err := response.Write(conn)
			if err != nil {
				fmt.Println("Error during write", err)
			}
		} else {
			response := http.Response{
				Status:     "404 Not Found",
				StatusCode: 404,
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
			}
			err := response.Write(conn)
			if err != nil {
				fmt.Println("Error during write", err)
			}
		}
	} else {
		fmt.Println("Bad formatted request")
		os.Exit(1)
	}
}

func parseHeader(header []string) http.Header {
	h := http.Header{}
	for _, v := range header {
		hFields := strings.Split(v, ": ")
		if len(hFields) == 2 {
			h.Add(hFields[0], hFields[1])
		}
	}
	return h
}

func parseStatusLine(statusLine string) (*StatusLine, error) {
	statusLineFields := strings.Fields(statusLine)
	if len(statusLineFields) < 3 {
		return nil, errors.New("Error trying to parse the status line")
	}
	return &StatusLine{
		Method:        statusLineFields[0],
		RequestTarget: statusLineFields[1],
		HttpVersion:   statusLineFields[2],
	}, nil
}
