package prepareTestWeb

import "html/template"

type (
	SessId       string        // Тип реализующий id сессии
	UsId         int           // Тип реализующий id пользователя в бд
	TopicId      uint          // Идентефикатор темы
	TopicName    string        // Имя тема
	TopicDesc    string        // Описание темы
	TestId       uint          // Уникальный идентефикатор теста
	TestName     string        // Имя теста
	TestDesc     string        // Описание теста
	AuthorName   string        // Имя автора
	TestRate     float32       // Рейтинг теста
	QuestionId   uint          // Уникальный идентефикатор вопроса
	QuestionType uint8         // Тип теста 0 - С письменным ответом, 1 - С выбором ответа, 2 - С выбором нескольких ответов, 3 - С выбором соответствия
	Body         template.HTML // Хранит html код
)

type SessUs struct {
	SessId
	UsId
}

type Test struct {
	TestId
	TestName
	TestDesc
	AuthorName
	TestRate
}

type Topic struct {
	TopicId
	TopicName
	TopicDesc
}
type Question struct {
	QuestionId
	Body
	TopicId
	QuestionType
}
