package moex

import (
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func moex() {

	dateFrom := os.Getenv("DATE_FROM")

	// Получаем открытые позиции

	var records [][]string

	var ticker = os.Getenv("TICKER")
	// GAZR
	// SBRF
	records, err := getOpenPositions(ticker, dateFrom)

	if err != nil {
		log.Fatal(err)
	}

	// Сохраняем в файловой системе
	writeToCSV(os.Getenv("OUTPUT_FILE"), records)

	// Получаем исторические данные
	records, err = getHistoryData(ticker, dateFrom)
	if err != nil {
		log.Fatal(err)
	}

	writeToCSV(os.Getenv("OUTPUT_FILE"), records)

	// Получаем агрегированые данные
	records, err = getAggregatesData(ticker, dateFrom)
	if err != nil {
		log.Fatal(err)
	}

	writeToCSV(os.Getenv("OUTPUT_FILE"), records)

}

func getHistoryData(ticker, dateFrom string) ([][]string, error) {

	var filename string

	// Создать папку с данными, если не существует
	if _, err := os.Stat(os.Getenv("DATA_DIR")); os.IsNotExist(err) {
		os.Mkdir(os.Getenv("DATA_DIR"), 0600)
	}

	if _, err := os.Stat(os.Getenv("DATA_DIR") + ticker); os.IsNotExist(err) {
		os.Mkdir(os.Getenv("DATA_DIR")+ticker, 0600)
	}

	var moexURL = os.Getenv("MOEX_HISTORY_URL") + ticker + dateFrom

	var items []string
	var records [][]string

	var body []byte

	type requestData struct {
		History struct {
			Columns []string        `json:"columns"`
			Data    [][]interface{} `json:"data"`
		} `json:"history"`
	}

	t, err := time.Parse("02-01-2006", dateFrom)

	if err != nil {
		log.Fatal(err)
	}

	dateString := moexURL + t.Format("2006-01-02")

	filename = os.Getenv("DATA_DIR") + ticker + "/history_" + t.Format("2006-01-02") + "-" + time.Now().Format("20060102") + ".json"

	//если файл уже существует
	log.Println("Ищем файл " + filename)
	if _, err := os.Stat(filename); err == nil {
		body, err = ioutil.ReadFile(filename)
		if err != nil {
			log.Print("Ошибка при загрузке данных из файла " + filename)
			log.Fatal(err)
		}
	} else {

		// Get the JSON response from the URL.
		response, err := http.Get(dateString)
		if err != nil {
			log.Print("Ошибка при получении ответа JSON ")
			log.Fatal(err)
		}
		defer response.Body.Close()
		// Read the body of the response into []byte.
		body, err = ioutil.ReadAll(response.Body)
		if err != nil {
			log.Print("Ошибка при получении тела ответа JSON ")
			log.Fatal(err)
		}

	}
	var rd requestData
	// Unmarshal the JSON data into the variable.
	if err := json.Unmarshal(body, &rd); err != nil {
		log.Print("Ошибка при разборе ответа JSON ")
		log.Fatal(err)
	}

	if len(rd.History.Data) > 0 {
		// TODO
		// fmt.Println(rd.History.Columns)
		// if rd.History.Columns ==  [BOARDID TRADEDATE SHORTNAME SECID NUMTRADES VALUE ADMITTEDQUOTE MP2VALTRD MARKETPRICE3TRADESVALUE ADMITTEDVALUE WAVAL] {
		// 	fmt.Println("формат данных совпадает")
		// }

		for _, record := range rd.History.Data {
			// fmt.Println(rd.History.Data[idx])

			items = []string{}

			items = append(items, record[1].(string))
			items = append(items, floatToString(record[4].(float64)))
			items = append(items, floatToString(record[5].(float64)))
			items = append(items, floatToString(record[6].(float64)))
			items = append(items, floatToString(record[7].(float64)))
			items = append(items, floatToString(record[8].(float64)))
			items = append(items, floatToString(record[11].(float64)))
			items = append(items, floatToString(record[12].(float64)))

			records = append(records, items)
		}

	}

	// Marshal the data.
	outputData, err := json.Marshal(rd)
	if err != nil {
		log.Print("Ошибка при упорядочивании данных ")
		log.Fatal(err)
	}
	// Save the marshalled data to a file.
	if err := ioutil.WriteFile(filename, outputData, 0644); err != nil {
		log.Print("Ошибка при сохранении данных ")
		log.Fatal(err)
	}

	return records, nil
}

// getOpenPositions получает открытые позиции
// Пример возвращаемых данных
// [[03-06-2019 123560 137283 70118 56395 387356] [04-06-2019 130427 114516 34502 50413 329858] [05-06-2019 132312 114258 33007 51061 330638]]
func getOpenPositions(ticker, dateFrom string) ([][]string, error) {

	var items []string
	var records [][]string
	var fl, s string

	currentDate := time.Now()
	dateString := "?d=" + dateFrom

	t, err := time.Parse("02-01-2006", dateFrom)

	if err != nil {
		log.Fatal(err)
	}

	for i := 0; currentDate.Format("20060102") != t.Format("20060102"); i++ {

		items = []string{}

		dateString = "?d=" + t.Format("20060102")

		fl, s, err = getOpenPositionsFromCSV(dateString, ticker)
		if err != nil {
			log.Println("При получении данных возникла ошибка:", err)
		}
		if s != "0" {
			items = append(items, t.Format("02-01-2006"))
			items = append(items, fl) // Лонги
			items = append(items, s)  // Совокупный объем открытых позиций

			records = append(records, items)
		}

		t = t.AddDate(0, 0, 1) // Subtract 1 Day

	}
	return records, nil
}

// getOpenPositionsFromCSV загружает открытые позиции из файла csv
func getOpenPositionsFromCSV(dateString, ticker string) (fl, s string, e error) {

	var data [][]string

	url := os.Getenv("MOEX_OP_URL") + dateString

	// Создать папку с данными, если не существует
	if _, err := os.Stat(os.Getenv("DATA_DIR")); os.IsNotExist(err) {
		os.Mkdir(os.Getenv("DATA_DIR"), 0600)
	}

	if _, err := os.Stat(os.Getenv("DATA_DIR") + ticker); os.IsNotExist(err) {
		os.Mkdir(os.Getenv("DATA_DIR")+ticker, 0600)
	}

	// Удаляем спецсимволы и whitespace-символы (пробелы, переносы и тп)
	// пунктуация (! [! - /: - @ [- `{- ~])
	// пробел (≡ [\ t \ n \ v \ f \ r])
	// См. https://code.tutsplus.com/ru/tutorials/regular-expressions-with-go-part-1--cms-30403

	// var re = regexp.MustCompile(`[[:punct:]]|[[:space:]]`)
	// filename := re.ReplaceAllString(dateString, "") + ".csv"

	// Без использования регулярки
	filename := "positions_" + strings.Trim(dateString, "d=")

	//если файл уже существует
	log.Println("Ищем файл " + os.Getenv("DATA_DIR") + ticker + "/" + filename)

	if _, err := os.Stat(os.Getenv("DATA_DIR") + ticker + "/" + filename); err == nil {
		log.Println("Нашли. Загружаем данные из файла ")

		data, err = readCSVFromFile(os.Getenv("DATA_DIR") + ticker + "/" + filename)

		if err != nil {
			log.Println("При загузке данных из файла произошла ошибка ", err)
			return "", "", err
		}
	} else {
		// Файл не существует
		// Скачиваем и парсим
		log.Println("Файл не существует. Пытаемся скачать и прочитать... ")
		data, err = readCSVFromURL(url)
		if err != nil {
			return "", "", err
		}

		err = writeToCSV(os.Getenv("DATA_DIR")+ticker+"/"+filename, data)
		if err != nil {
			log.Println("Ошибка при сохранении csv:", err)
			return "", "", err
		}

		data, err = readCSVFromFile(os.Getenv("DATA_DIR") + ticker + "/" + filename)

		if err != nil {
			log.Println("При загузке данных из файла произошла ошибка ", err)
			return "", "", err
		}

	}

	var FsummG, UsummG int

	for _, row := range data {

		if row[1] == ticker {
			if row[3] == "F" {

				if row[4] == "1,0000" {
					fl = row[8]
				} else {
					i7, _ := strconv.Atoi(row[7])
					i8, _ := strconv.Atoi(row[8])
					UsummG = i7 + i8
				}
			}

		}
	}

	s = strconv.Itoa(FsummG + UsummG)

	return fl, s, nil

}

func getAggregatesData(ticker, dateFrom string) ([][]string, error) {

	var moexURL = os.Getenv("MOEX_AG_URL") + ticker
	var body []byte

	// Создать папку с данными, если не существует
	if _, err := os.Stat(os.Getenv("DATA_DIR")); os.IsNotExist(err) {
		os.Mkdir(os.Getenv("DATA_DIR"), 0600)
	}

	if _, err := os.Stat(os.Getenv("DATA_DIR") + ticker); os.IsNotExist(err) {
		os.Mkdir(os.Getenv("DATA_DIR")+ticker, 0600)
	}

	type requestData []struct {
		Charsetinfo struct {
			Name string `json:"name"`
		} `json:"charsetinfo,omitempty"`
		Aggregates []struct {
			MarketName  string  `json:"market_name"`  // Рынок
			MarketTitle string  `json:"market_title"` // Рынок название
			Engine      string  `json:"engine"`
			Tradedate   string  `json:"tradedate"`
			SecID       string  `json:"secid"`
			Value       float64 `json:"value"`     // Объём сделок, руб
			Volume      int     `json:"volume"`    // Объём сделок, шт
			Numtrades   int     `json:"numtrades"` // Количество сделок
		} `json:"aggregates,omitempty"`
	}

	var items []string
	var records [][]string

	var filename string

	currentDate := time.Now()
	dateString := "&date=" + dateFrom

	t, err := time.Parse("02-01-2006", dateFrom)

	if err != nil {
		log.Fatal(err)
	}

	for i := 0; currentDate.Format("20060102") != t.Format("20060102"); i++ {

		items = []string{}

		dateString = "&date=" + t.Format("20060102")

		// Без использования регулярки
		filename = os.Getenv("DATA_DIR") + ticker + "/aggregate_" + strings.Trim(dateString, "&date=") + ".json"

		//если файл уже существует
		log.Println("Ищем файл " + filename)
		if _, err := os.Stat(filename); err == nil {
			log.Println("Нашли. Загружаем данные из файла " + filename)
			body, err = ioutil.ReadFile(filename)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Println("Файла нет. Скачиваем данные с адреса " + moexURL + "&date=" + dateString)
			// Get the JSON response from the URL.
			response, err := http.Get(moexURL + dateString)
			if err != nil {
				log.Fatal(err)
			}
			defer response.Body.Close()
			// Read the body of the response into []byte.
			body, err = ioutil.ReadAll(response.Body)
			if err != nil {
				log.Fatal(err)
			}
		}

		var rd requestData
		// Unmarshal the JSON data into the variable.
		if err := json.Unmarshal(body, &rd); err != nil {
			log.Fatal(err)
		}

		if len(rd[1].Aggregates) > 0 {

			shares, ndm, otc, repo, moexboard := 0.0, 0.0, 0.0, 0.0, 0.0
			for i := 0; i < len(rd[1].Aggregates); i++ {

				switch rd[1].Aggregates[i].MarketName {
				case "shares": // Рынок акций
					shares = rd[1].Aggregates[i].Value
				case "ndm": // Режим переговорных сделок
					ndm = rd[1].Aggregates[i].Value
				case "otc": // Внебиржевые сделки ОТС
					otc = rd[1].Aggregates[i].Value
				case "repo": // Рынок сделок РЕПО
					repo = rd[1].Aggregates[i].Value
				case "moexboard": // MOEX Board
					moexboard = rd[1].Aggregates[i].Value
				}

			}
			items = append(items, t.Format("02-01-2006"))
			items = append(items, floatToString(shares))
			items = append(items, floatToString(ndm))
			items = append(items, floatToString(otc))
			items = append(items, floatToString(repo))
			items = append(items, floatToString(moexboard))

			records = append(records, items)

			//	fmt.Println(rd[1].Aggregates[0].MarketTitle)

			// dateString = strings.Trim(dateString, "&date=")
			// fmt.Print(dateString)

			// fmt.Print(" Основной рынок: " + floatToString(sharesValue))
			// fmt.Println("; Внебиржевые сделки: " + floatToString(otcValue))

			// fmt.Println("Доля " + "Внебиржевые сделки/Основной рынок" + ": " + floatToString((otcValue / sharesValue * 100)) + " %")

		} else {
			// dateString = strings.Trim(dateString, "&date=")
			// fmt.Println(dateString + ": Для этой даты нет данных ")
		}

		// Marshal the data.
		outputData, err := json.Marshal(rd)
		if err != nil {
			log.Fatal(err)
		}
		// Save the marshalled data to a file.
		if err := ioutil.WriteFile(filename, outputData, 0644); err != nil {
			log.Fatal(err)
		}

		t = t.AddDate(0, 0, 1) // Subtract 1 Day

	}
	return records, nil

}

// Функции для работы с CSV
// readCSVFromURL читает данные CSV с URL
func readCSVFromURL(url string) ([][]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	reader := csv.NewReader(resp.Body)
	reader.Comma = ';'
	data, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return data, nil
}

// readCSVFromFile читает данные CSV из файла
func readCSVFromFile(filename string) ([][]string, error) {

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close() // this needs to be after the err check

	data, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}
	return data, nil
}

// writeToCSV сохраняет данные CSV в файл
func writeToCSV(csvfilename string, records [][]string) error {
	csvfile, err := os.Create(csvfilename)
	if err != nil {
		log.Println("Error:", err)
		return err
	}
	defer csvfile.Close()

	writer := csv.NewWriter(csvfile)
	for _, record := range records {
		err := writer.Write(record)
		if err != nil {
			log.Println("Error:", err)
			return err
		}
	}
	writer.Flush()
	return nil
}

// Прочие функции
// floatToString конвертирует float number в строку
func floatToString(input float64) string {
	return strconv.FormatFloat(input, 'f', 2, 64)
}
