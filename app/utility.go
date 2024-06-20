package main

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func encodeBody(body string, encoding string) (string, string) {
	for _, enc := range strings.Split(encoding, ", ") {
		if enc == "gzip" {
			var b bytes.Buffer
			gz := gzip.NewWriter(&b)
			gz.Write([]byte(body))
			gz.Close()
			return string(b.Bytes()[:]), enc
		}
	}
	return body, ""
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

func parseStatusLine(statusLine string) (StatusLine, error) {
	statusLineFields := strings.Fields(statusLine)
	if len(statusLineFields) < 3 {
		return StatusLine{}, errors.New("Error trying to parse the status line")
	}
	return StatusLine{
		Method:        statusLineFields[0],
		RequestTarget: statusLineFields[1],
		HttpVersion:   statusLineFields[2],
	}, nil
}

func parseData(conn net.Conn) (StatusLine, http.Header, string) {
	h := http.Header{}
	var statusLine StatusLine
	bufInit := make([]byte, 8*1024)
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	conn.Read(bufInit)
	lines := strings.Split(string(bufInit), "\r\n")
	endHeaders := len(lines[0]) + 2
	// read status line
	statusLine, _ = parseStatusLine(lines[0])
	// scan headers
	for _, line := range lines[1:] {
		endHeaders += 2
		if line == "" {
			break
		}
		hFields := strings.Split(line, ": ")
		if len(hFields) == 2 {
			h.Add(hFields[0], hFields[1])
		}
		endHeaders += len(line)
	}
	fmt.Println(h)
	// reading body
	contentLength, err := strconv.Atoi(h.Get("Content-Length"))
	fmt.Println("Expecting to read ", contentLength)
	if err != nil {
		contentLength = 0
		return statusLine, h, string([]byte{})
	}
	remaining := contentLength - (len(bufInit) - endHeaders)
	buf := make([]byte, remaining) // big buffer
	conn.Read(buf)
	buf = append(bufInit[endHeaders:], buf...)

	return statusLine, h, string(buf)
}
