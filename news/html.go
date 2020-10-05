package news

import (
	"net"
	"net/http"
	"time"
	"regexp"

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
// MaxDepth is 2, so only the links on the scraped page
		// and links on those pages are visited
		colly.MaxDepth(2),
		colly.Async(true),
		// Visit only root url and urls which start with "e" or "h" on httpbin.org
		colly.URLFilters(
			regexp.MustCompile("http://httpbin\\.org/(|e.+)$"),
			regexp.MustCompile("http://httpbin\\.org/h.+"),
		),
	)

		// Limit the maximum parallelism to 2
	// This is necessary if the goroutines are dynamically
	// created to control the limit of simultaneous requests.
	//
	// Parallelism can be controlled also by spawning fixed
	// number of go routines.
// Limit the number of threads started by colly to two
	// when visiting links which domains' matches "*httpbin.*" glob
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*httpbin.*",
		Parallelism: 2,
		RandomDelay: 5 * time.Second,
	})

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
	log.Debug("Visiting ", r.URL.String())
})

// Set HTML callback
	// Won't be called if error occurs
	// c.OnHTML("*", func(e *colly.HTMLElement) { // all page
	// 	log.Println(e)
	// })

	// Before making a request put the URL with
	// the key of "url" into the context of the request
	c.OnRequest(func(r *colly.Request) {
		r.Ctx.Put("url", r.URL.String())
	})

	// After making a request get "url" from
	// the context of the request
	c.OnResponse(func(r *colly.Response) {
		log.Debug(r.Ctx.Get("url"))
	})

	// Set error handler
	c.OnError(func(r *colly.Response, err error) {
		log.Warn("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})


c.OnResponse(func(r *colly.Response) {
	log.Debug("Visited ", r.Request.URL)
})

c.OnHTML("tr td:nth-of-type(1)", func(e *colly.HTMLElement) {
	log.Debug("First column of a table row:", e.Text)
})


c.OnHTML("li.list__item", func(e *colly.HTMLElement) {
	 log.Debug("Found news:", e.Text)
	 log.Debug("link:", e.ChildAttr("a.home-link","href"))

})

c.OnXML("//h1", func(e *colly.XMLElement) {
	log.Debug(e.Text)
})

c.OnScraped(func(r *colly.Response) {
	log.Debug("Finished ", r.Request.URL)
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