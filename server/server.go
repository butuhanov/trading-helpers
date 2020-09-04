package server

import (
	"context"
	"net"

	log "github.com/sirupsen/logrus"

	"github.com/butuhanov/trading-helpers/news"
	"google.golang.org/grpc"

	ps "github.com/butuhanov/trading-helpers/proto"
)

type NewsServiceServer struct {
}

func (s *NewsServiceServer) GetNews(ctx context.Context,
	req *ps.MessageParams) (*ps.LastNews, error) {

	var err error
	response := new(ps.LastNews)

	a := "./news/example_data/sources.txt"
	b := "./news/example_data/keywords.txt"
	res, err := news.CheckNews(a, b)
	log.Println(res, err)

	return response, err
}

func StartServer() {
	log.Info("starting server...")

	server := grpc.NewServer()

	instance := new(NewsServiceServer)

	ps.RegisterNewsServiceServer(server, instance)

	log.Debug("creating grpc listener...")

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Unable to create grpc listener:", err)
	}

	log.Debug("start server serve listener...")

	if err = server.Serve(listener); err != nil {
		log.Fatal("Unable to start server:", err)
	}

	log.Info("server started...")
}
