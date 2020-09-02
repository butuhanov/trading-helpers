package news

import (
	"bufio"
	"encoding/json"
	"html"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

// Глобальный массив хэшей новостей которые уже были обработаны
var knownNews []string

// максимальный размер буфера известных новостей
const maxNewsLength = 200

// максимальная глубина поиска
const maxDepth = 20

// таймаут запроса
const httpGetTimeout = 3

func parseSource() {
	// TODO: function to parse source
}

func checkKeyWord(data []string, keyword string) string {
	// TODO: function to check keyword in the source
	// parseSource()

	// log.Printf("ищем:%v", keyword)

	for _, el := range data {
		var m News
		err := json.Unmarshal([]byte(el), &m)
		checkError(err)

		// log.Printf("len: %d, cap: %d arr:%v\n", len(knownNews), cap(knownNews), knownNews)

		// Поиск ключевого слова в заголовке

		if m.Error != "" {
			log.Println("При получении данных получена ошибка", m.Error, "пропускаем...")
		} else {
			if strings.Contains(strings.ToLower(m.Title), keyword) {
				log.Println(m.SourceTitle, "\n", m.Date, "\nНайдено", keyword, "в заголовке:", html.UnescapeString(m.Title), m.Link)
			}
			if strings.Contains(strings.ToLower(m.Description), keyword) {
				log.Println(m.SourceTitle, "\n", m.Date, "\nНайдено", keyword, "в описании:", m.Title, "-", html.UnescapeString(m.Description), m.Link)
			}
		}

		// Поиск ключевого слова в описании

	}

	// log.Println(data[1])

	return strconv.Itoa(len(data))
}

// CheckNews возвращает вхождения ключевых слов в новостных источниках в виде массива
// Входные параметры - массивы источников и ключевых слов
func CheckNews(sourceFile, keywordFile string) ([]string, error) {

	var result = make([]string, 0)

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
	log.Println("вместимость канала", cap(dataCh))

	wg := new(sync.WaitGroup)
	wg.Add(cap(dataCh))
	for _, source := range sources {

		go readRSS(source, wg, dataCh) // перебираем все источники
	}

	wg.Wait()
	log.Println("все запросы выполнены")

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

	log.Printf("получено записей:%v", len(data))
	// log.Printf("данные:%v", data)

	for _, keyword := range keywords { // перебираем все ключевые слова
		result = append(result, checkKeyWord(data, keyword))

	}

	return result, nil

}

// Стандартная обработка ошибок
func checkError(err error) {
	if err != nil {
		log.Println("При выполнении операции произошла ошибка:", err)
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
