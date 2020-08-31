package news

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
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

type News struct {
	Title       string
	Description string
	Link        string
	Hash        string
}

func readRSS(source string) ([]string, error) {

	log.Println("checking", source)

	response, err := http.Get(source)

	checkError(err)

	defer response.Body.Close()

	XMLdata, err := ioutil.ReadAll(response.Body)

	checkError(err)

	rss := new(Rss)

	buffer := bytes.NewBuffer(XMLdata)

	decoded := xml.NewDecoder(buffer)

	err = decoded.Decode(rss)

	checkError(err)

	// fmt.Printf("Title : %s\n", rss.Channel.Title)
	// fmt.Printf("Description : %s\n", rss.Channel.Description)
	// fmt.Printf("Link : %s\n", rss.Channel.Link)

	total := len(rss.Channel.Items)

	if total > maxDepth {
		total = maxDepth
	}

	// fmt.Printf("Total items : %v\n", total)

	var result = make([]string, 0)

	for i := 0; i < total; i++ {

		title := rss.Channel.Items[i].Title
		description := rss.Channel.Items[i].Description
		link := rss.Channel.Items[i].Link
		hash := getMD5Hash(title + description + link)

		data := &News{
			Title:       title,
			Description: description,
			Link:        link,
			Hash:        hash,
		}

		json, err := json.Marshal(data)

		checkError(err)

		result = append(result, string(json))

	}

	return result, nil

}

func getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

// findElement takes a slice and looks for an element in it. If found it will
// return it's key, otherwise it will return -1 and a bool of false.
func findElement(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			log.Printf("нашли на позиции %v\n", i)
			return i, true
		}
	}
	return -1, false
}
