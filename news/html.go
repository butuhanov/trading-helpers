package news

import (
	"net"
	"net/http"
	"time"

	"github.com/gocolly/colly/v2"

	log "github.com/sirupsen/logrus"
)


func ReadHTML() {
	log.Debug("парсим html")

	// Instantiate default collector
	// ollector manages the network communication and responsible for the execution of the attached callbacks while a collector job is running.
	// Create a collector with default settings:
//c1 := colly.NewCollector()

	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("yandex.ru", "www.yandex.ru"),
		// change User-Agent and url revisit options
		colly.UserAgent("xy"),
	colly.AllowURLRevisit(),
	)

	// Colly uses Golang’s default http client as networking layer. HTTP options can be tweaked by changing the default HTTP roundtripper.
	c.WithTransport(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	})

	// You can attach different type of callback functions to a Collector to control a collecting job or retrieve information.
// Before making a request print "Visiting ..."
c.OnRequest(func(r *colly.Request) {
	log.Println("Visiting", r.URL.String())
})
c.OnError(func(_ *colly.Response, err error) {
	log.Println("Something went wrong:", err)
})

c.OnResponse(func(r *colly.Response) {
	log.Println("Visited", r.Request.URL)
})

c.OnHTML("tr td:nth-of-type(1)", func(e *colly.HTMLElement) {
	log.Println("First column of a table row:", e.Text)
})


c.OnHTML("li.list__item", func(e *colly.HTMLElement) {
	 log.Println("Found news:", e.Text)
	 log.Println("link:", e.ChildAttr("a.home-link","href"))

})

c.OnXML("//h1", func(e *colly.XMLElement) {
	log.Println(e.Text)
})

c.OnScraped(func(r *colly.Response) {
	log.Println("Finished", r.Request.URL)
})

// Full list of collector attributes https://godoc.org/github.com/gocolly/colly#Collector

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		// link := e.Attr("href")
		// Print link
		// log.Printf("Link found: %q -> %s\n", e.Text, link)
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
	//	c.Visit(e.Request.AbsoluteURL(link))
	})



	// Start scraping on https://hackerspaces.org
	c.Visit("https://yandex.ru/")
}