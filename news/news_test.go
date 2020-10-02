package news_test

import (
	"encoding/json"
	"fmt"

	"github.com/butuhanov/trading-helpers/news"
)

func ExampleCheckNews() {
	a := "sources.txt"
	b := "keywords.txt"
	res, err := news.CheckNews(a, b)
	result, err := json.Marshal(res)

	fmt.Println(string(result), err)
	// Output:null <nil>
}
