package main

import (
	"context"
	"time"

	"github.com/butuhanov/trading-helpers/server"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	// "context"
	// "time"
	ps "github.com/butuhanov/trading-helpers/proto"
)

func main() {

	go server.StartServer()

	time.Sleep(1000 * time.Millisecond)

	log.Info("starting client...")

	conn, _ := grpc.Dial("127.0.0.1:8080", grpc.WithInsecure())
	client := ps.NewNewsServiceClient(conn)

	resp, err := client.GetNews(context.Background(),
		&ps.MessageParams{})

	if err != nil {
		log.Fatal("could not get answer: ", err)
	}

	log.Println("Output:\n", string([]byte(resp.News)), err)

	// a := "./news/example_data/sources.txt"
	// b := "./news/example_data/keywords.txt"

	// res, err := news.CheckNews(a, b)

	// log.Println(string(res), err)
}
