package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
)

func handleConnection(conn net.Conn) {
	fmt.Printf("Accept %v\n", conn.RemoteAddr())
	// Read a request
	request, err := http.ReadRequest(bufio.NewReader(conn))
	if err != nil {
		panic(err)
	}
	dump, err := httputil.DumpRequest(request, true)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(dump))
	// Write a response
	response := http.Response{
		StatusCode: 200,
		ProtoMajor: 1,
		ProtoMinor: 0,
		Body:       ioutil.NopCloser(strings.NewReader("Hello world\n")),
	}
	response.Write(conn)
	conn.Close()
}

func main() {
	// Initialize signal
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT)

	path := filepath.Join(os.TempDir(), "unixdomainsocket-sample")
	os.Remove(path)
	listener, err := net.Listen("unix", path)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	fmt.Println("Server is running at " + path)

	receiveSigint := false

	go func() {
		<-signals
		receiveSigint = true
		listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			if receiveSigint {
				break
			} else {
				panic(err)
			}
		}
		go handleConnection(conn)
	}
}
