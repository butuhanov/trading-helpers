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
	"time"
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
	SourceTitle string
	Title       string
	Description string
	Link        string
	Hash        string
}

func readRSS(source string) ([]string, error) {

	// log.Println("проверяем", source)

	c := &http.Client{
		Timeout: httpGetTimeout * time.Second,
	}

	response, err := c.Get(source)

	if err != nil {
		return nil, err
		// log.Fatal(err)
	}

	defer response.Body.Close()

	XMLdata, err := ioutil.ReadAll(response.Body)

	checkError(err)

	rss := new(Rss)

	buffer := bytes.NewBuffer(XMLdata)

	decoded := xml.NewDecoder(buffer)

	err = decoded.Decode(rss)

	if err != nil {
		log.Println("ошибка при декодировании:", err)
		return nil, err
	}

	// fmt.Printf("Title : %s\n", rss.Channel.Title)
	// fmt.Printf("Description : %s\n", rss.Channel.Description)
	// fmt.Printf("Link : %s\n", rss.Channel.Link)

	sourceTitle := rss.Channel.Title

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

		// Если новость уже просмотрена, то переходим к следующей
		_, ok := findElement(knownNews, hash)
		// log.Printf("ищем элемент %v в %v, результат %v\n", m.Hash, knownNews, ok)
		if ok {
			// log.Println("новость", title, "уже проверяли.. пропускаем..")
			continue
		}

		data := &News{
			SourceTitle: sourceTitle,
			Title:       title,
			Description: description,
			Link:        link,
			Hash:        hash,
		}

		json, err := json.Marshal(data)

		checkError(err)

		if len(knownNews) < maxNewsLength {
			knownNews = append(knownNews, hash)
		} else {
			knownNews = append(knownNews[maxNewsLength:], knownNews[1:]...)
			knownNews = append(knownNews, hash)
		}

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
			// log.Printf("нашли на позиции %v\n", i)
			return i, true
		}
	}
	return -1, false
}
