package server

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"
)

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

type Result struct {
	Value int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *Result) error {
	reply.Value = args.A * args.B
	return nil
}

func (t *Arith) Divide(args *Args, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}

	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}

func inspect(conn io.ReadWriteCloser) {
	fmt.Println("---> inspect start")

	var req rpc.Request
	req = rpc.Request{}

	dec := json.NewDecoder(conn)
	fmt.Println("-------> decode start")
	if err := dec.Decode(&req); err != nil {
		fmt.Println("decode error")
		panic(err)
	}
	fmt.Println("-------> decode OK")

	codec := jsonrpc.NewServerCodec(conn)
	fmt.Println("---> ReadRequestHeader start")
	err := codec.ReadRequestHeader(&req)
	if err != nil {
		if err != io.EOF {
			log.Println("rpc:", err)
		}
	}
	fmt.Println("---> ReadRequestHeader end")
	fmt.Println("---> inspect end")
}

func printHTTPConn(conn net.Conn) {
	fmt.Println("---> dump HTTP Conn start")
	request, err := http.ReadRequest(bufio.NewReader(conn))
	if err != nil {
		fmt.Println("printConn: ReadRequest error")
		panic(err)
	}
	dump, err := httputil.DumpRequest(request, true)
	if err != nil {
		fmt.Println("printConn: DumpRequest error")
		panic(err)
	}
	fmt.Println(string(dump))
	fmt.Println("---> dump HTTP Conn end")
}

func printHTTPRequest(r *http.Request) {
	fmt.Println("----> request:")
	fmt.Printf("%v\n", *r)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("----> body:")
	fmt.Printf("%s\n", body)
}

func handleConnection(conn net.Conn) {
	fmt.Printf("Accept %v\n", conn.RemoteAddr())
	request, err := http.ReadRequest(bufio.NewReader(conn))
	if err != nil {
		fmt.Println("printConn: ReadRequest error")
		panic(err)
	}
	body := TmpReadWriteCloser{Rw: new(bytes.Buffer)}
	_, err = io.Copy(body, request.Body)
	if err != nil {
		log.Fatal(err)
	}
	// timeout does not work
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	rpc.ServeCodec(jsonrpc.NewServerCodec(body))

	response := http.Response{
		StatusCode: 200,
		ProtoMajor: 1,
		ProtoMinor: 0,
		Body:       body,
	}
	response.Write(conn)
	conn.Close()
}

type TmpReadWriteCloser struct {
	Rw *bytes.Buffer
}

func (f TmpReadWriteCloser) Read(b []byte) (n int, err error) {
	nn, err := f.Rw.Read(b)
	if err != nil {
		return nn, err
	}
	return nn, nil
}

func (f TmpReadWriteCloser) Write(b []byte) (n int, err error) {
	nn, err := f.Rw.Write(b)
	if err != nil {
		return nn, err
	}
	return nn, nil
}

func (f TmpReadWriteCloser) Close() error {
	return nil
}

func handleJSONRPC(w http.ResponseWriter, r *http.Request) {
	body := TmpReadWriteCloser{Rw: new(bytes.Buffer)}
	_, err := io.Copy(body, r.Body)
	if err != nil {
		log.Fatal(err)
	}

	rpc.ServeCodec(jsonrpc.NewServerCodec(body))

	result, err := ioutil.ReadAll(body)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, string(result))
}

func Serve() {
	arith := new(Arith)
	rpc.Register(arith)

	// s := &http.Server{
	// 	Addr:           "localhost:8888",
	// 	Handler:        http.HandlerFunc(handleJSONRPC),
	// 	ReadTimeout:    5 * time.Second,
	// 	WriteTimeout:   5 * time.Second,
	// 	MaxHeaderBytes: 1 << 20,
	// }

	// log.Println(s.ListenAndServe())

	listener, err := net.Listen("tcp", "localhost:8888")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	// http.HandleFunc("/jrpc", handleJSONRPC)
	// http.Serve(listener, nil)

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go handleConnection(conn)
	}
}
