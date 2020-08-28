package news

import (
	"log"
	"strconv"
)

func parseSource() {
	// TODO: function to parse source
}

func checkKeyWord(source, keyword string) string {
	// TODO: function to check keyword in the source
	// parseSource()

	_, titles := readRSS(source)

	log.Printf("lenght:%v", len(titles))

	log.Println(titles)

	return strconv.Itoa(len(titles))
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
