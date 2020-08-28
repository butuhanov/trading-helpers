package news

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Rss struct {
	Channel Channel `xml:"channel"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

func readRSS(source string) (error, []string) {

	log.Println("checking", source)

	response, err := http.Get(source)

	if err != nil {
		log.Println(err)
		fmt.Errorf("Cannot response %v", source)
		return err, nil
	}

	defer response.Body.Close()

	XMLdata, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rss := new(Rss)

	buffer := bytes.NewBuffer(XMLdata)

	decoded := xml.NewDecoder(buffer)

	err = decoded.Decode(rss)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Title : %s\n", rss.Channel.Title)
	fmt.Printf("Description : %s\n", rss.Channel.Description)
	fmt.Printf("Link : %s\n", rss.Channel.Link)

	total := len(rss.Channel.Items)

	fmt.Printf("Total items : %v\n", total)

	var result = make([]string, 0)

	for i := 0; i < total; i++ {
		result = append(result, rss.Channel.Items[i].Title)

		// fmt.Printf("[%d] item title : %s\n", i, rss.Channel.Items[i].Title)
		// fmt.Printf("[%d] item description : %s\n", i, rss.Channel.Items[i].Description)
		// fmt.Printf("[%d] item link : %s\n\n", i, rss.Channel.Items[i].Link)
	}

	return nil, result

}
