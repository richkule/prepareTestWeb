package prepareTestWeb

import (
	"html/template"
)

type (
	//Header    template.HTML Содержит шапку страниццы, не используется как отдельный тип потому, что нарушает логику
	// работы функции обработки шаблонов
	RegWrong  string     // Текст об ошибке регистрации передаваемый на страницу
	AutoWrong string     // Текст об ошибке авторизации, передеваемый на страницу
	Tests     []Test     // Срез содержащий тесты
	Topics    []Topic    // Срез содержащий темы тестов
	Questions []Question // Срез содержащий вопросы
	UserName  string     // Имя пользователя
)

// Для заполнения шаблона шапки
type DataHeader struct {
	UserName
}

// Для заполнения шаблона редактирования теста
type DataTest struct {
	Header template.HTML
	Tests
}

// Для заполнения шаблона редактирования теста
type DataEditTest struct {
	Header template.HTML
	Topics
	TestId
}

// Для заполнения шаблона редкатроивания темы
type DataEditTopic struct {
	Header template.HTML
	Questions
	TopicId
}

// Для заполнения шаблона главной страницы
type DataIndex struct {
	Header template.HTML
}

// Для заполнения шаблона авторизации
type DataLogin struct {
	AutoWrong
	RegWrong
}

// Для заполнения шаблона создания теста
type DataCreateTest struct {
	Header template.HTML
	Tests
}
