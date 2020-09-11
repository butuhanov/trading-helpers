package main

import (
	// "context"
	// "time"
	"fmt"
	// "github.com/butuhanov/trading-helpers/server"
	"github.com/butuhanov/trading-helpers/news"
	// log "github.com/sirupsen/logrus"
	// "google.golang.org/grpc"

	// ps "github.com/butuhanov/trading-helpers/proto"
)

func main() {

	// go server.StartServer()

	// time.Sleep(1000 * time.Millisecond)

	// log.Info("starting client...")

	// conn, _ := grpc.Dial("127.0.0.1:8080", grpc.WithInsecure())
	// client := ps.NewNewsServiceClient(conn)

	// resp, err := client.GetNews(context.Background(),
	// 	&ps.MessageParams{})

	// if err != nil {
	// 	log.Fatal("could not get answer: ", err)
	// }

	// log.Println("Output:\n", string((resp.News)), err)

	a := "./sources.txt"
	b := "./keywords.txt"

	res, _ := news.CheckNews(a, b)

	fmt.Println(string(res))
}
