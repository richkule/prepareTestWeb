// Нынче не работает, раньше успешно парсил вопросы
package dataParse

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"unicode"
)

type question struct {
	name       string   // Имя вопроса
	answers    []string // Ответы
	rightAnsId int      // Индекс правильного ответа на вопрос
}
type topic struct {
	name      string     // Имя темы с вопросами
	questions []question // Вопросы, содержащиеся в теме
}

const (
	// Адрес сайта, с которого происходит парсинг вопросов
	rootHref      = "https://softmakerkz.blogspot.com/search/label/Ответы%20на%20вопросы%20по%201С%208.3"
	hrefSelector  = ".post-title a"    // Селектор на искомую ссылку
	topicSelector = "h3"               // Селектор на назавание темы с вопросами
	bodySelector  = "div[dir = 'ltr']" // Селектор на содержимое поста
	lenAnsId      = 3                  // Длина номера ответа "n. "
	minLenContent = 5                  // Минимальная длина тега, для его глубокой проверки(каждого из элементов тега)
)

// Получение обработанного содержимого передаваемого url
func getParsedHTTP(url string) *goquery.Document {
	// Получение HTTP запроса
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Обработка HTML документа
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return doc
}

// Функция, собирающая адреса с корневого сайта
func grabHrefs() []string {
	hrefArr := make([]string, 0, 14) // Массив возвращаемых адресов
	doc := getParsedHTTP(rootHref)
	doc.Find(hrefSelector).Each(func(i int, s *goquery.Selection) {
		// Сбор ссылок из каждого найденного элемента
		href, isExists := s.Attr("href")
		if isExists {
			hrefArr = append(hrefArr, href)
		}
	})
	return hrefArr
}

// Функция, удаляющая спецсимволы и пробелы слева от строки
func lStrip(text string) string {
	var stripInd int
	for ind, val := range text {
		if !unicode.IsSpace(rune(val)) {
			stripInd = ind
			break
		}
	}
	return text[stripInd:]
}

// Функция, собирающая вопросы для темы данной ссылки
func grabInfo(href string) topic {
	doc := getParsedHTTP(href)
	topicName := lStrip(doc.Find(topicSelector).Text()) // Название темы
	postBody := doc.Find(bodySelector)                  // Тело поста(в котором содержатся вопросы)
	var qArr []question                                 // Массив вопросов
	var q question                                      // Структура вопроса
	flag := false                                       // Флаг, что ведется запись верного ответа
	//Функция, парсящая HTML в массив вопросов
	parseFun := func(i int, s2 *goquery.Selection) {
		// Если найденная строка - название вопроса
		if s2.Nodes[0].Data == "h2" {
			if q.name != "" {
				qArr = append(qArr, q)
			}
			q = question{name: lStrip(s2.Text())}
			return
		}
		// Если предыдущая строка с ответом оборвалась(ответ является верным)
		if flag {
			q.answers[len(q.answers)-1] = q.answers[len(q.answers)-1] + s2.Text()
			flag = false
			return
		}
		// Добавление ответа в вопрос
		temp := lStrip(s2.Text())
		if temp != "" {
			if unicode.IsDigit(rune(temp[0])) {
				if len(temp) <= lenAnsId {
					flag = true
					q.answers = append(q.answers, temp)
					q.rightAnsId = int(temp[0]) - '0' - 1
				} else {
					q.answers = append(q.answers, temp)
				}
			}
		}
	}
	postBody.Contents().Each(func(i int, s *goquery.Selection) {
		// Не глубокая(не каждого элемента тега) проверка
		if s.Contents().Length() < minLenContent {
			parseFun(i, s)
		} else { // Глубокая проверка, в случае если тег содержит ответы или вопросы, а не является их содержимым
			s.Contents().Each(parseFun)
		}
	})
	qArr = append(qArr, q)
	tSt := topic{topicName, qArr}
	return tSt
}

func main() {
	a := grabHrefs()
	tArr := make([]topic, 14)
	ch := make(chan bool)
	for ind, val := range a {
		ind, val := ind, val
		go func() {
			tArr[ind] = grabInfo(val)
			ch <- true
			return
		}()
	}
	// Ожидание завершение всех go рутин
	for i := 0; i < len(a); i++ {
		<-ch
	}
	for _, elem := range tArr {
		println("Название темы " + elem.name)
		for _, elem2 := range elem.questions {
			println("Вопрос " + elem2.name)
		}
	}
}
