package news_test

import (
	"fmt"

	"github.com/butuhanov/trading-helpers/news"
)

func ExampleCheckNews() {
	a := "./example_data/sources.txt"
	b := "./example_data/keywords.txt"
	res, err := news.CheckNews(a, b)
	fmt.Println(res, err)
	// Output:[] open ./example_data/sources.txt: no such file or directory
}
