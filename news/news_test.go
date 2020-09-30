package news_test

import (
	"fmt"

	"github.com/butuhanov/trading-helpers/news"
)

func ExampleCheckNews() {
	a := "sources.txt"
	b := "keywords.txt"
	res, err := news.CheckNews(a, b)
	fmt.Println(string(res), err)
	// Output:open sources.txt: no such file or directory
}
