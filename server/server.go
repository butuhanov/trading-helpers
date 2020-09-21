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

	a := req.Sources
	b := req.Keywords

	res, err := news.CheckNews(a, b)

	// stringByte := "\x00" + strings.Join(res, "\x20\x00") // x20 = space and x00 = null

	// response.News = []byte(stringByte)

	response.News = []byte(res)

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
