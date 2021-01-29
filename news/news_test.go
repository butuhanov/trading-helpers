package news_test

import (
	"encoding/json"
	"fmt"

	"github.com/butuhanov/trading-helpers/news"
)

func ExampleCheckNews() {
	a := "/tmp/sources.txt"
	b := "/tmp/keywords.txt"
	res, err := news.CheckNews(a, b)
	result, err := json.Marshal(res)

	fmt.Println(string(result), err)
	// Output:null <nil>
}


func ExampleReadHTML() {

	res, err := news.ReadHTML("http://yandex.ru")

	result, err := json.Marshal(res)

	fmt.Println(string(result), err)
	// Output:[] <nil>
}
