package news

func parseSource() {
	// TODO: function to parse source
}

func checkKeyWord(source, keyword string) string {
	// TODO: function to check keyword in the source
	parseSource(source)

	return
}

// CheckNews возвращает вхождения ключевых слов в новостных источниках в виде массива
// Входные параметры - массивы источников и ключевых слов
func CheckNews(sources, keywords []string) []string {

	var result = make([]string)

	for _, source := range sources { // перебираем все источники

		for _, keyword := range keywords { // перебираем все ключевые слоав

			append(result, checkKeyWord(source, keyword))

		}

	}

	return result

}
