package news_test

import (
	"fmt"

	"github.com/butuhanov/trading-helpers/news"
)

func ExampleCheckNews() {
	a := make([]string, 2)
	b := make([]string, 2)

	a[0] = "https://www.interfax.ru/rss.asp"
	a[1] = "https://www.interfax.ru/rss.asp"
	b[0] = "сша"
	b[1] = "газпром"
	// b[1] = "test4"
	fmt.Println(news.CheckNews(a, b))
	// Output: [20 20 20 20]
}
