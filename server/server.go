package server

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"io"
	"net"
	"net/http"

	"github.com/butuhanov/trading-helpers/news"
	ps "github.com/butuhanov/trading-helpers/proto"
	"github.com/golang/protobuf/jsonpb"
	_struct "github.com/golang/protobuf/ptypes/struct"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
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

	resStruct := news.Results{}


	for _,v:=range(res){
		// log.Debug(i, v)

		resStruct.NewsItem = append(resStruct.NewsItem, v)

	}



data :=&_struct.Struct{Fields: make(map[string]*_struct.Value)}

cc, _ := json.Marshal(resStruct)
		var bb bytes.Buffer
		bb.Write(cc)
		if err := (&jsonpb.Unmarshaler{}).Unmarshal(&bb, data); err != nil {
			log.Fatal(err)
	}

		response.News = append(response.News, data)

	log.Debug(len(resStruct.NewsItem))

	return response, err
}


func StartServer( grpcPort, httpPort string) {
		log.Info("starting server...")

		log.Info("starting grpc server...")
		server := grpc.NewServer()

		instance := new(NewsServiceServer)

		ps.RegisterNewsServiceServer(server, instance)
		listener, err := net.Listen("tcp", "localhost:"+grpcPort)

		go server.Serve(listener)

		log.Info("starting http server...")
		var (
			// command-line options:
			// gRPC server endpoint
			grpcServerEndpoint = flag.String("grpc-server-endpoint",  "localhost:"+httpPort, "gRPC server endpoint")
		)
		log.Info("endpoint:"+*grpcServerEndpoint)

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
		log.Fatal(http.ListenAndServe("localhost:"+httpPort, mux))

		// check example curl -X POST http://127.0.0.1:23456/v1/news/last -d '{"sources":"news/example_data/sources.txt", "keywords":"news/example_data/keywords.txt"}'

}

// http.HandleFunc("/health-check", HealthCheckHandler)
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// Очень простой хендлер проверки состояния.
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	// В будущем мы хотим сообщать сообщать о состоянии
	// базы данных или кеша (например Redis) выполняя
	// простой PING и отдавать все это в запросе
	io.WriteString(w, `{"alive": true}`)
}


// Стандартная обработка ошибок
func checkError(err error) {
	if err != nil {
		log.Warn("При выполнении операции произошла ошибка:", err)
	}
}
