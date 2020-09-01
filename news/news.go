package news

import (
	"encoding/json"
	"html"
	"log"
	"strconv"
	"strings"
)

// Глобальный массив хэшей новостей которые уже были обработаны
var knownNews []string

// максимальный размер буфера известных новостей
const maxNewsLength = 200

// максимальная глубина поиска
const maxDepth = 20

func parseSource() {
	// TODO: function to parse source
}

func checkKeyWord(data []string, keyword string) string {
	// TODO: function to check keyword in the source
	// parseSource()

	log.Printf("ищем:%v", keyword)

	for _, el := range data {
		var m News
		err := json.Unmarshal([]byte(el), &m)
		checkError(err)

		// log.Printf("len: %d, cap: %d arr:%v\n", len(knownNews), cap(knownNews), knownNews)

		// Поиск ключевого слова в заголовке
		// log.Println(m.Title)
		if strings.Contains(strings.ToLower(m.Title), keyword) {
			log.Println("Найдено в заголовке:", html.UnescapeString(m.Title))
		}
		if strings.Contains(strings.ToLower(m.Description), keyword) {
			log.Println("Найдено в описании:", m.Title, "-", html.UnescapeString(m.Description))
		}

		// Поиск ключевого слова в описании

	}

	// log.Println(data[1])

	return strconv.Itoa(len(data))
}

// CheckNews возвращает вхождения ключевых слов в новостных источниках в виде массива
// Входные параметры - массивы источников и ключевых слов
func CheckNews(sources, keywords []string) []string {

	var result = make([]string, 0)

	for _, source := range sources { // перебираем все источники

		data, err := readRSS(source)
		if err != nil {
			log.Println("ошибка при парсинге, пропускаем источник")
			continue
		}
		log.Printf("получено записей:%v", len(data))

		for _, keyword := range keywords { // перебираем все ключевые слова
			result = append(result, checkKeyWord(data, keyword))

		}

	}

	return result

}

// Стандартная обработка ошибок
func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
