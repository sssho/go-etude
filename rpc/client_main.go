package main

import (
	"fmt"
	"net/rpc"
	"rpc-etude/server"
)

func main() {
	// Unix domain socket can be used
	// path := filepath.Join(os.TempDir(), "rpcetudesocket")
	// client, err := rpc.DialHTTP("unix", path)

	client, err := rpc.DialHTTP("tcp", "localhost:8888")
	if err != nil {
		panic(err)
	}

	args := &server.Args{7, 8}
	var reply int
	err = client.Call("Arith.Multiply", args, &reply)
	if err != nil {
		panic(err)
	}
	fmt.Println("Arith: %d*%d=%d", args.A, args.B, reply)
}
