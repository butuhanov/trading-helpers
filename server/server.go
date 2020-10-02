package server

import (
	"context"
	"net"

	"net/http"

	"flag"

	"strconv"
	"strings"

	"bytes"

	"encoding/json"

	"github.com/golang/protobuf/jsonpb"

	// "os"
	// "os/signal"
	// "time"

	log "github.com/sirupsen/logrus"

	"github.com/butuhanov/trading-helpers/news"
	"google.golang.org/grpc"

	ps "github.com/butuhanov/trading-helpers/proto"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	_struct "github.com/golang/protobuf/ptypes/struct"
)

type NewsServiceServer struct {
}


	// Результаты
	type Result struct {
		Keyword     string `json:"keyword"`
		Date        string `json:"date"`
		Source      string `json:"source"`
		Place       string `json:"place"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Link        string `json:"link"`
	}


	type Results struct {
    Results []Result `json:"results"`
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


	result, err := json.Marshal(res)

	checkError(err)

	//  log.Debug(string(result))

	// for i, v := range result {

	// 	log.Debug(i, v)

	// }

	// response.News, err = strconv.Unquote(string(result))

	resStruct := []Result{}
	// resStruct := Results{}

	jsonErr := json.Unmarshal(result, &resStruct)

	log.Debug(len(resStruct))

	// for i,v:=range(resStruct){
	// 	log.Debug(i, v)
	// }


	if jsonErr != nil {
			log.Fatal(jsonErr)
	}


	data :=&_struct.Struct{Fields: make(map[string]*_struct.Value)}


	for i,v:=range(resStruct){
		log.Debug(i, v)
		cc, _ := json.Marshal(v)
		var bb bytes.Buffer
		bb.Write(cc)
		if err := (&jsonpb.Unmarshaler{}).Unmarshal(&bb, data); err != nil {
			log.Fatal(err)
	}

		response.News = append(response.News, data)
	}


	// checkError(err)

	// log.Debug(response)


	return response, err
}


func StartServer() {
	log.Info("starting server...")

	log.Info("starting grpc server...")

	server := grpc.NewServer()

	instance := new(NewsServiceServer)

	ps.RegisterNewsServiceServer(server, instance)
	listener, err := net.Listen("tcp", "localhost:8080")
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
	log.Fatal(http.ListenAndServe("localhost:23456", mux))

	// check example curl -X POST http://127.0.0.1:23456/v1/news/last -d '{"sources":"news/example_data/sources.txt", "keywords":"news/example_data/keywords.txt"}'

	log.Info("server started...")
}


// Стандартная обработка ошибок
func checkError(err error) {
	if err != nil {
		log.Warn("При выполнении операции произошла ошибка:", err)
	}
}


func convert( b []byte ) string {
	s := make([]string,len(b))
	for i := range b {
			s[i] = strconv.Itoa(int(b[i]))
	}
	return strings.Join(s,",")
}