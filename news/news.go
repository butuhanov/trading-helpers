package news

import (
	"encoding/json"
	"log"
	"strconv"
)

// Глобальный массив хэшей новостей которые уже были обработаны
var knownNews []string

// максимальный размер буфера известных новостей
const maxNewsLength = 5

// максимальная глубина поиска
const maxDepth = 20

func parseSource() {
	// TODO: function to parse source
}

func checkKeyWord(source, keyword string) string {
	// TODO: function to check keyword in the source
	// parseSource()

	data, err := readRSS(source)

	checkError(err)

	log.Printf("lenght:%v", len(data))

	for _, el := range data {
		var m News
		err := json.Unmarshal([]byte(el), &m)
		checkError(err)

		// Если новость уже просмотрена, то переходим к следующей
		_, ok := findElement(knownNews, m.Hash)
		log.Printf("ищем элемент %v в %v, результат %v\n", m.Hash, knownNews, ok)
		if ok {
			continue
		}

		if len(knownNews) < maxNewsLength {
			knownNews = append(knownNews, m.Hash)
		} else {
			knownNews = append(knownNews[maxNewsLength:], knownNews[1:]...)
			knownNews = append(knownNews, m.Hash)
		}

		log.Printf("len: %d, cap: %d arr:%v\n", len(knownNews), cap(knownNews), knownNews)

		log.Println(m.Title)
	}

	// log.Println(data[1])

	return strconv.Itoa(len(data))
}

// CheckNews возвращает вхождения ключевых слов в новостных источниках в виде массива
// Входные параметры - массивы источников и ключевых слов
func CheckNews(sources, keywords []string) []string {

	var result = make([]string, 0)

	for _, source := range sources { // перебираем все источники

		for _, keyword := range keywords { // перебираем все ключевые слова

			log.Println("checking:", keywords, source)

			result = append(result, checkKeyWord(source, keyword))

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
