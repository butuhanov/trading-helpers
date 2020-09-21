package main

import (
	"context"
	"time"

	"flag"
	// "fmt"
	"strings"

	"github.com/butuhanov/trading-helpers/server"
	// "github.com/butuhanov/trading-helpers/news"
	ps "github.com/butuhanov/trading-helpers/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {

	var sources = flag.String("sources", "sources.txt", "The 'sources' option value")
	var keywords = flag.String("keywords", "keywords.txt", "The 'keywords' option value")

	// parse flag's options
	flag.Parse()

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

	log.Println("Output:\n", string((resp.News)), err)

	// res, err := news.CheckNews(*sources, *keywords)
	// if err != nil {
	// 	log.Warn("При выполнении операции произошла ошибка:", err)
	// 	fmt.Println("{error:", err, "}")
	// }

	// fmt.Printf("%v", trimSuffix(string(res), ","))

	// fmt.Println(string(res))
}

func trimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}
