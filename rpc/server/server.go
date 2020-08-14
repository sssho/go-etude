package server

import (
	"errors"
	"net"
	"net/http"
	"net/rpc"
)

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
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

func Serve() {
	arith := new(Arith)
	rpc.Register(arith)
	// This internally calls http.Handle()
	rpc.HandleHTTP()

	// Unix domain socket can be used
	// path := filepath.Join(os.TempDir(), "rpcetudesocket")
	// l, err := net.Listen("unix", path)

	l, err := net.Listen("tcp", "localhost:8888")
	if err != nil {
		panic(err)
	}
	defer l.Close()
	http.Serve(l, nil)
}
