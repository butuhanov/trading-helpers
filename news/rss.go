package news

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"

	// "log"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type Rss struct {
	Channel Channel `xml:"channel"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Date        string `xml:"pubDate"`
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
	Date        string
	Hash        string
	Error       string
}

func readRSS(source string, wg *sync.WaitGroup, dataCh chan []string) error {

	defer wg.Done()

	// log.Println("проверяем", source)

	var result = make([]string, 0)

	c := &http.Client{
		Timeout: httpGetTimeout * time.Second,
	}

	response, err := c.Get(source)

	if err != nil {
		generateWarn("Ошибка при получении ответа", source, err.Error(), dataCh)
		return err
	}

	defer response.Body.Close()

	XMLdata, err := ioutil.ReadAll(response.Body)

	if err != nil {
		generateWarn("Ошибка при чтении ответа", source, err.Error(), dataCh)
		return err
	}

	rss := new(Rss)

	buffer := bytes.NewBuffer(XMLdata)

	decoded := xml.NewDecoder(buffer)

	err = decoded.Decode(rss)

	if err != nil {
		generateWarn("Ошибка при декодировании", source, err.Error(), dataCh)
		return err
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

	for i := 0; i < total; i++ {

		title := rss.Channel.Items[i].Title
		description := rss.Channel.Items[i].Description
		link := rss.Channel.Items[i].Link
		date := rss.Channel.Items[i].Date
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
			Date:        date,
			Hash:        hash,
		}

		json, err := json.Marshal(data)

		if err != nil {
			generateWarn("Ошибка при формировании JSON", source, err.Error(), dataCh)
			return err
		}

		if len(knownNews) < maxNewsLength {
			knownNews = append(knownNews, hash)
		} else {
			knownNews = append(knownNews[maxNewsLength:], knownNews[1:]...)
			knownNews = append(knownNews, hash)
		}

		result = append(result, string(json))

	}
	dataCh <- result
	return nil

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

func generateWarn(text, source, err string, dataCh chan []string) {
	data := &News{
		Link:  source,
		Error: err,
	}
	var result = make([]string, 0)

// Showing file, function name, and line number with logrus:
	if pc, file, line, ok := runtime.Caller(1); ok {
		file = file[strings.LastIndex(file, "/")+1:]
		funcName := runtime.FuncForPC(pc).Name()
		log.WithFields(
			log.Fields{
				"msg":    err,
				"source": source,
				"code":    fmt.Sprintf("%s:%s:%d", file, funcName, line),
			}).Warn(text)
	}

	json, err1 := json.Marshal(data)
	checkError(err1)

	result = append(result, string(json))
	dataCh <- result

}
