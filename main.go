package main

import (
	"context"
	"fmt"
	"io"

	// "time"
	"golang.org/x/net/websocket"

	// "io"
	"net"
	"net/http"
	"net/rpc"

	// "flag"
	// "fmt"
	"strings"

	// "github.com/butuhanov/trading-helpers/server"

	// "github.com/butuhanov/trading-helpers/news"
	// ps "github.com/butuhanov/trading-helpers/proto"
	log "github.com/sirupsen/logrus"

	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"

	grpchello "google.golang.org/grpc/examples/helloworld/helloworld"

	"github.com/butuhanov/trading-helpers/news"
	ps "github.com/butuhanov/trading-helpers/proto"
)


type exampleHTTPHandler struct{}

func (h *exampleHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug("ServeHTTP...")
    fmt.Fprintf(w, "example http response\n")
}

func serveHTTP(l net.Listener) {
	log.Debug("ServeHTTP-2...")
    s := &http.Server{
        Handler: &exampleHTTPHandler{},
    }
    if err := s.Serve(l); err != cmux.ErrListenerClosed {
        panic(err)
    }
}

func EchoServer(ws *websocket.Conn) {
	log.Debug("EchoServer...")
    if _, err := io.Copy(ws, ws); err != nil {
        panic(err)
    }
}

func serveWS(l net.Listener) {
	log.Debug("serveWS...")
    s := &http.Server{
        Handler: websocket.Handler(EchoServer),
    }
    if err := s.Serve(l); err != cmux.ErrListenerClosed {
        panic(err)
    }
}

type ExampleRPCRcvr struct{}

func (r *ExampleRPCRcvr) Cube(i int, j *int) error {
	log.Debug("Cube...")
    *j = i * i
    return nil
}

func serveRPC(l net.Listener) {
	log.Debug("serveRPC...")
    s := rpc.NewServer()
    if err := s.Register(&ExampleRPCRcvr{}); err != nil {
        panic(err)
    }
    for {
        conn, err := l.Accept()
        if err != nil {
            if err != cmux.ErrListenerClosed {
                panic(err)
            }
            return
        }
        go s.ServeConn(conn)
    }
}

type grpcServer struct{}

func (s *grpcServer) SayHello(ctx context.Context, in *grpchello.HelloRequest) (
    *grpchello.HelloReply, error) {
			log.Debug("SayHello...")
    return &grpchello.HelloReply{Message: "Hello " + in.Name + " from cmux"}, nil
}

type NewsServiceServer struct {
}


func (s *NewsServiceServer) GetNews(ctx context.Context,
	req *ps.MessageParams) (*ps.LastNews, error) {

		log.Debug("GetNews...")

	var err error
	response := new(ps.LastNews)

	a := req.Sources
	b := req.Keywords

	res, err := news.CheckNews(a, b)

	// stringByte := "\x00" + strings.Join(res, "\x20\x00") // x20 = space and x00 = null

	// response.News = []byte(stringByte)

	response.News = []byte(res)

	return response, err
}

func serveGRPC(l net.Listener) {
	log.Debug("serveGRPC...")
    grpcs := grpc.NewServer()
		// grpchello.RegisterGreeterServer(grpcs, &grpcServer{})
		ps.RegisterNewsServiceServer(grpcs, &NewsServiceServer{})
    if err := grpcs.Serve(l); err != cmux.ErrListenerClosed {
        panic(err)
		}



}

func main() {

	// var sources = flag.String("sources", "sources.txt", "The 'sources' option value")
	// var keywords = flag.String("keywords", "keywords.txt", "The 'keywords' option value")

	// // parse flag's options
	// flag.Parse()

	// go server.StartServer()

	// time.Sleep(1000 * time.Millisecond)

	// log.Info("starting client...")

	// conn, _ := grpc.Dial("127.0.0.1:8080", grpc.WithInsecure())
	// client := ps.NewNewsServiceClient(conn)



	// log.Debug("start HTTP server...")

	// helloHandler := func(w http.ResponseWriter, req *http.Request) {
	// 	io.WriteString(w, "Hello, world!\n")
	// }

	// http.HandleFunc("/hello", helloHandler)

	// log.Fatal(http.ListenAndServe(":8082", nil))


	// resp, err := client.GetNews(context.Background(),
	// 	&ps.MessageParams{Sources: *sources, Keywords: *keywords})

	// if err != nil {
	// 	log.Fatal("could not get answer: ", err)
	// }

	// log.Println("Output:\n", string((resp.News)), err)


	// res, err := news.CheckNews(*sources, *keywords)
	// if err != nil {
	// 	log.Warn("При выполнении операции произошла ошибка:", err)
	// 	fmt.Println("{error:", err, "}")
	// }

	// fmt.Printf("%v", trimSuffix(string(res), ","))

	// fmt.Println(string(res))

// Create the main listener.
l, err := net.Listen("tcp", ":23456")
if err != nil {
	log.Fatal(err)
}

// Create a cmux.
m := cmux.New(l)

// Match connections in order:
// First grpc, then HTTP, and otherwise Go RPC/TCP.
grpcL := m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
httpL := m.Match(cmux.HTTP1Fast())
trpcL := m.Match(cmux.Any()) // Any means anything that is not yet matched.

// Create your protocol servers.
grpcS := grpc.NewServer()
// grpchello.RegisterGreeterServer(grpcS, &server{})
ps.RegisterNewsServiceServer(grpcS, &NewsServiceServer{})

httpS := &http.Server{
	Handler: &exampleHTTPHandler{},
}

trpcS := rpc.NewServer()
trpcS.Register(&ExampleRPCRcvr{})

// Use the muxed listeners for your servers.
go grpcS.Serve(grpcL)
go httpS.Serve(httpL)
go trpcS.Accept(trpcL)

// Start serving!
m.Serve()


	// // We first match the connection against HTTP2 fields. If matched, the
	// // connection will be sent through the "grpcl" listener.
	// grpcl := m.Match(cmux.HTTP2HeaderFieldPrefix("content-type", "application/grpc"))

	// //Otherwise, we match it againts a websocket upgrade request.
	// wsl := m.Match(cmux.HTTP1HeaderField("Upgrade", "websocket"))

	// // Otherwise, we match it againts HTTP1 methods. If matched,
	// // it is sent through the "httpl" listener.
	// httpl := m.Match(cmux.HTTP1Fast())
	// // If not matched by HTTP, we assume it is an RPC connection.
	// rpcl := m.Match(cmux.Any())

	// // Then we used the muxed listeners.
	// go serveGRPC(grpcl)
	// go serveWS(wsl)
	// go serveHTTP(httpl)
	// go serveRPC(rpcl)

	// if err := m.Serve(); !strings.Contains(err.Error(), "use of closed network connection") {
	// 		panic(err)
	// }

	// Start serving!
// m.Serve()

}

func trimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}
