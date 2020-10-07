package news

import (
	"bufio"
	"encoding/json"

	"crypto/md5"
	"encoding/hex"

	"regexp"
	"sort"

	"os"
	"strings"
	"sync"

	"path/filepath"

	log "github.com/sirupsen/logrus"
)

// Глобальный массив хэшей новостей которые уже были обработаны
var knownNews []string

// максимальный размер буфера известных новостей
const maxNewsLength = 200

// максимальная глубина поиска
const maxDepth = 20

// таймаут запроса
const httpGetTimeout = 5

	// Возвращаем JSON в формате
	// Ключевое слово
	// Дата
	// Источник
	// Где найдено
	// Заголовок
	// Описание
	// Ссылка


// Новости
type News struct {
	SourceTitle string
	Title       string
	Description string
	Link        string
	Date        string
	Hash        string
	Error       string
}

	// Результаты
	type Result struct {
		Keyword     string `json:"keyword"`
		Date        string `json:"date"`
		Source      string `json:"source"`
		Place       string `json:"place"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Link        string `json:"link"`
		Hash        string `json:"hash"`
	}

	type Results struct {
    NewsItem []Result `json:"news_item"`
}

func init() {
	// Log as JSON instead of the default ASCII formatter.
	// log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	// log.SetLevel(log.ErrorLevel)
	// log.SetLevel(log.InfoLevel)
	log.SetLevel(log.DebugLevel)

	// log.SetReportCaller(true)

	// log.SetFormatter(&log.TextFormatter{
	// 	DisableColors: false,
	// 	FullTimestamp: false,
	// })

	// log.Trace("Something very low level.")
	// log.Debug("Useful debugging information.")
	// log.Info("Something noteworthy happened!")
	// log.Warn("You should probably take a look at this.")
	// log.Error("Something failed but I'm not quitting.")

}

func parseSource() {
	// TODO: function to parse source
}

func checkKeyWord(data []string, keyword string) []Result {



	// TODO: function to check keyword in the source
	// parseSource()

	log.Debug("ищем:", keyword, " ====================================================")

	// log.Info(data)

	var result = make([]byte, 0)

	var resultStruct []Result


	for _, el := range data {
		var m News
		err := json.Unmarshal([]byte(el), &m)
		checkError(err)

		// log.Info(m)

		// log.Printf("len: %d, cap: %d arr:%v\n", len(knownNews), cap(knownNews), knownNews)

		// Поиск ключевого слова в заголовке

		if m.Error != "" {
			log.Info("При получении данных получена ошибка в источнике ", m.Link, " ошибка:", m.Error)

			if strings.HasPrefix(m.Error, "context deadline exceeded") {
				log.Debug("ошибка по таймауту, возможно надо повторить попытку")

		}


		} else {
			if strings.Contains(strings.ToLower(m.Title), keyword) {
				log.Debug("Найдено ", keyword, " в заголовке")

				r := Result{keyword, m.Date, m.SourceTitle, "заголовок", m.Title, RemoveHtmlTags(m.Description), m.Link, m.Hash}
				b, err := json.Marshal(r)
				checkError(err)

				resultStruct = append(resultStruct, r)

				result = append(result, b...)
				result = append(result, ',')

			} else {
				if strings.Contains(strings.ToLower(m.Description), keyword) {
					log.Debug("Найдено ", keyword, " в описании")

					r := Result{keyword, m.Date, m.SourceTitle, "описание", m.Title, RemoveHtmlTags(m.Description), m.Link, m.Hash}
					b, err := json.Marshal(r)
					checkError(err)

					resultStruct = append(resultStruct, r)

					result = append(result, b...)
					result = append(result, ',')

				}
			}
		}

		// Поиск ключевого слова в описании

	}

	// log.Debug(resultStruct)

	return resultStruct

	// return strconv.Itoa(len(data))
}

// CheckNews возвращает вхождения ключевых слов в новостных источниках в виде массива
// Входные параметры - массивы источников и ключевых слов
func CheckNews(sourceFile, keywordFile string) ([]Result, error) {


	log.Debug("Проверяем ", sourceFile, " и ", keywordFile)

	var resultStruct []Result

	// resultStruct := make(Results, 0)

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("текущая директория", dir)

	// var result = make([]byte, 0)

	var sources, keywords []string

	sourcesFromFile, err := readDataFromFile(sourceFile)
	if err != nil {
		return nil, err
	}
	keywordsFromFile, err := readDataFromFile(keywordFile)
	if err != nil {
		return nil, err
	}
	sources = append(sources, sourcesFromFile...)
	keywords = append(keywords, keywordsFromFile...)

	// dataCh := make(chan []string)  // канал для результатов запроса
	dataCh := make(chan []string, len(sources)) // буферизованный
	log.Debug("вместимость канала:", cap(dataCh))

	wg := new(sync.WaitGroup)
	wg.Add(cap(dataCh))
	for _, source := range sources {

		go readRSS(source, wg, dataCh) // перебираем все источники
	}

	wg.Wait()
	log.Debug("все запросы выполнены")

	var data []string
	var dataHTML []string

	for i := 0; i < cap(dataCh); i++ {
		data = append(data, <-dataCh...)
	}

	close(dataCh)

	log.Debug("получено записей:", len(data))
	// log.Printf("данные:%v", data)


	for _, el := range data { // перебираем все ключевые слова

		var m News
		err := json.Unmarshal([]byte(el), &m)
		checkError(err)

		// Если ошибка с XML пробуем распарсить HTML
		if strings.HasPrefix(m.Error, "XML syntax error") {
			log.Debug("ошибка с XML, пробуем распарсить HTML")

			res, err := ReadHTML(m.Link)
			if err ==nil {
				dataHTML = append(dataHTML, res...)
			} else {
				log.Warn("ошибка при парсинге HTML страницы ", err.Error())
			}

		}

	// resultStruct = append(resultStruct, checkKeyWord(data, keyword)...)

}


	for _, keyword := range keywords { // перебираем все ключевые слова для RSS потоков

		resultStruct = append(resultStruct, checkKeyWord(data, keyword)...)

	}


	for _, keyword := range keywords { // перебираем все ключевые слова для html страниц

		resultStruct = append(resultStruct, checkKeyWord(dataHTML, keyword)...)

	}

	return resultStruct, nil

}

// Стандартная обработка ошибок
func checkError(err error) {
	if err != nil {
		log.Warn("При выполнении операции произошла ошибка:", err)
	}
}

func readDataFromFile(source string) ([]string, error) {

	var result []string

	file, err := os.Open(source)
	if err != nil {
		return nil, err
		// log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text()[0] != 35 && scanner.Text()[0] != 47 {
			result = append(result, scanner.Text())
		}

	}

	if err := scanner.Err(); err != nil {
		return nil, err
		// log.Fatal(err)
	}

	return result, nil

}



func RemoveHtmlTags(in string) string {
	// regex to match html tag
	const pattern = `(<\/?[a-zA-A]+?[^>]*\/?>)*`
	r := regexp.MustCompile(pattern)
	groups := r.FindAllString(in, -1)
	// should replace long string first
	sort.Slice(groups, func(i, j int) bool {
		return len(groups[i]) > len(groups[j])
	})
	for _, group := range groups {
		if strings.TrimSpace(group) != "" {
			in = strings.ReplaceAll(in, group, "")
		}
	}
	return in
}

func updateKnownNews(hash string){
	if len(knownNews) < maxNewsLength {
		knownNews = append(knownNews, hash)
	} else {
		knownNews = append(knownNews[maxNewsLength:], knownNews[1:]...)
		knownNews = append(knownNews, hash)
	}
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
