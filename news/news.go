package news

import (
	"bufio"
	"encoding/json"

	// "log"
	"os"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

// Глобальный массив хэшей новостей которые уже были обработаны
var knownNews []string

// максимальный размер буфера известных новостей
const maxNewsLength = 200

// максимальная глубина поиска
const maxDepth = 20

// таймаут запроса
const httpGetTimeout = 4

func init() {
	// Log as JSON instead of the default ASCII formatter.
	// log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
	// log.SetLevel(log.DebugLevel)

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

func checkKeyWord(data []string, keyword string) []byte {

	// Возвращаем JSON в формате
	// Ключевое слово
	// Дата
	// Источник
	// Где найдено
	// Заголовок
	// Описание
	// Ссылка

	type Result struct {
		Keyword     string
		Date        string
		Source      string
		Place       string
		Title       string
		Description string
		Link        string
	}

	// TODO: function to check keyword in the source
	// parseSource()

	log.Debug("ищем:", keyword, " ====================================================")

	// log.Info(data)

	var result = make([]byte, 0)

	for _, el := range data {
		var m News
		err := json.Unmarshal([]byte(el), &m)
		checkError(err)

		// log.Info(m)

		// log.Printf("len: %d, cap: %d arr:%v\n", len(knownNews), cap(knownNews), knownNews)

		// Поиск ключевого слова в заголовке

		if m.Error != "" {
			log.Info("При получении данных получена ошибка, пропускаем ", m.Link)
		} else {
			if strings.Contains(strings.ToLower(m.Title), keyword) {
				log.Debug("Найдено ", keyword, " в заголовке")

				r := Result{keyword, m.Date, m.SourceTitle, "заголовок", m.Title, m.Description, m.Link}
				b, err := json.Marshal(r)
				checkError(err)

				result = append(result, b...)

			} else {
				if strings.Contains(strings.ToLower(m.Description), keyword) {
					log.Debug("Найдено ", keyword, " в описании")

					r := Result{keyword, m.Date, m.SourceTitle, "описание", m.Title, m.Description, m.Link}
					b, err := json.Marshal(r)
					checkError(err)

					result = append(result, b...)

				}
			}
		}

		// Поиск ключевого слова в описании

	}

	return result

	// return strconv.Itoa(len(data))
}

// CheckNews возвращает вхождения ключевых слов в новостных источниках в виде массива
// Входные параметры - массивы источников и ключевых слов
func CheckNews(sourceFile, keywordFile string) ([]byte, error) {

	var result = make([]byte, 0)

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

	for i := 0; i < cap(dataCh); i++ {
		data = append(data, <-dataCh...)
	}

	close(dataCh)
	// data, err := readRSS(source)

	// if err != nil {
	// 	log.Println("ошибка при парсинге:", err, "пропускаем источник")
	// 	continue
	// }
	// log.Printf("получено записей:%v", len(data))

	log.Debug("получено записей:", len(data))
	// log.Printf("данные:%v", data)

	for _, keyword := range keywords { // перебираем все ключевые слова
		result = append(result, checkKeyWord(data, keyword)...)

	}

	return result, nil

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
