package main

import (
	"context"

	// "time"

	// "io"
	"net"
	"net/http"

	"flag"
	// "fmt"
	"strings"

	// "github.com/butuhanov/trading-helpers/server"

	log "github.com/sirupsen/logrus"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"google.golang.org/grpc"

	"github.com/butuhanov/trading-helpers/news"
	ps "github.com/butuhanov/trading-helpers/proto"
)
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



log.Info("starting grpc server...")

	server := grpc.NewServer()

	instance := new(NewsServiceServer)

	ps.RegisterNewsServiceServer(server, instance)
	listener, err := net.Listen("tcp", ":8080")
	go server.Serve(listener)


	log.Info("starting http server...")
var (
  // command-line options:
  // gRPC server endpoint
  grpcServerEndpoint = flag.String("grpc-server-endpoint",  "localhost:8080", "gRPC server endpoint")
)

ctx := context.Background()
  ctx, cancel := context.WithCancel(ctx)
	defer cancel()


	// Register gRPC server endpoint
  // Note: Make sure the gRPC server is running properly and accessible
  mux := runtime.NewServeMux()
  opts := []grpc.DialOption{grpc.WithInsecure()}
	// err := gw.RegisterYourServiceHandlerFromEndpoint(ctx, mux,  *grpcServerEndpoint, opts)
	err = ps.RegisterNewsServiceHandlerFromEndpoint(ctx, mux,  *grpcServerEndpoint, opts)
  if err != nil {
    log.Fatal(err)
  }


  // Start HTTP server (and proxy calls to gRPC server endpoint)
	log.Fatal(http.ListenAndServe("127.0.0.1:23456", mux))


}

func trimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}
