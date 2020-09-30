package prepareTestWeb

import (
	"github.com/gorilla/sessions"
	"html/template"
	"net/http"
	"regexp"
)

const (
	LogPage        = `static/html/login.html`      // Ссылка на страницу регистрации
	IndexPage      = `static/html/index.html`      // Ссылка на главную страницу
	HeaderPath     = `static/html/header.html`     // Ссылка на шапку для страниц
	TestPage       = `static/html/test.html`       // Ссылка на тестовую страницу
	CreateTestPage = `static/html/createTest.html` // Ссылка на страницу создания теста
	EditTestPage   = `static/html/editTest.html`   // Ссылка на страницу редактирования теста
	EditTopicPage  = `static/html/editTopic.html`  // Ссылка на страницу редактирования темы
	SessionName    = "Session"                     // Название сессии в которой хранятся id пользователя
	CookRowName    = "SessId"                      // Название поля сесии в котором хранятся куки авторизированных пользователей
	GuestUserId    = 1                             // Идентефикатор не уникального пользователя
)

var (
	RegPass      = regexp.MustCompile(`[a-zA-Z0-9-_.@#$%]{8,20}$`) // Регулярное выражение для пароля
	RegLog       = regexp.MustCompile(`^[a-zA-Z0-9-_.]{5,20}$`)    // Регулярное выражение для логина
	RegTestEdit  = regexp.MustCompile(`^/edit/test/(\d+)/?$`)      // Шаблон пути для редактирования теста
	RegTopicEdit = regexp.MustCompile(`^/edit/topic/(\d+)/?$`)     // Шаблон пути для редактирования темы
	RegQuesEdit  = regexp.MustCompile(`^/edit/question/(\d+)/?$`)  // Шаблон пути для редактирования вопроса
	Store        = sessions.NewCookieStore([]byte("TestProject"))
	Templates    = template.Must(template.ParseFiles(LogPage, IndexPage, HeaderPath, TestPage, CreateTestPage, EditTestPage, EditTopicPage))
)

// Тип ручек для работы с идентефицированными пользователями и возвратом ошибки
type HandlerIdFunc func(http.ResponseWriter, *http.Request, *SessUs) error
