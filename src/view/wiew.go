package view

import (
	"html/template"
	"log"
	"net/http"
	"web-application/src/model"
)

// Хранит части HTML кода подставляемые в “page.html” и
// список ссылок на JavaScript скрипты которые должны быть добавлены к странице.
type mainPageContent struct {
	Name        string
	Menu        template.HTML
	Content     template.HTML
	JavaScripts []string
}

// Отрисовка основной страницы, страницы авторизации в приложении
func MainPage(w *http.ResponseWriter) {
	template, err := template.ParseFiles("src/templates/main.html")
	if err != nil {
		http.Error(*w, "Ошибка отрисовки страницы.", http.StatusInternalServerError)
	}
	template.Execute(*w, nil)
}

// Отрисовка страницы ошибки
func ErrorPage(w *http.ResponseWriter, message string, statusCode int) {
	http.Error(*w, message, statusCode)
}

// Отрисовка страницы регистрации
func Registration(w *http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		templ, err := template.ParseFiles("src/templates/pages/registration.html")
		if err != nil {
			log.Println("Не удалось открыть HTML страницу с регистрацией: ", err)
			http.Error(*w, "Ошибка отрисовки страницы.", http.StatusInternalServerError)
		}
		templ.Execute(*w, nil)
	}
}

// Отрисовка страницы профиля пользователя
func Profile(w *http.ResponseWriter, urlPath *string, authorizedUser *model.UserInfo) {
	var user model.UserInfo
	err := user.GetDataFromDB(model.QuryUrlPathProfileWhere(urlPath))
	if err != nil {
		// если пользователь не найден в базе данных
		ErrorPage(w, "Пользователь не найден", http.StatusNotFound)
		return
	}
	var content mainPageContent
	// Название страницы в метаданных
	content.Name = "Профиль пользователя"
	// Получаем HTML код главного меню
	content.Menu, err = menuHTMLcontent(authorizedUser, "@")
	if err != nil {
		log.Println("Не удалось отрисовать главное меню: ", err)
		ErrorPage(w, "Ошибка страницы", http.StatusInternalServerError)
		return
	}
	// Получаем HTML код основного содержимого страницы
	content.Content, err = profileHTMLcontent(urlPath)
	if err != nil {
		log.Println("Не удалось отрисовать основное содержимое страницы профиля: ", err)
		ErrorPage(w, "Ошибка страницы", http.StatusInternalServerError)
		return
	}

	templ, err := template.ParseFiles("src/templates/page.html")
	if err != nil {
		log.Fatal("Не удалось открыть основной шаблон страницы: ", err)
	}
	// Рендерим страницу и отправляем её клиенту
	err = templ.Execute(*w, content)
	if err != nil {
		// если не удалось отрисовать страницу по шаблону
		ErrorPage(w, "Ошибка страницы", http.StatusInternalServerError)
		return
	}
}

// Отрисовка страницы со списком заданий
func Tasks(w *http.ResponseWriter, authorizedUser *model.UserInfo) {
	var content mainPageContent
	var err error
	// Название страницы в метаданных
	content.Name = "Список заданий"
	// Получаем HTML код главного меню
	content.Menu, err = menuHTMLcontent(authorizedUser, "tasks")
	if err != nil {
		log.Println("Не удалось отрисовать главное меню: ", err)
		ErrorPage(w, "Ошибка страницы", http.StatusInternalServerError)
		return
	}
	// Получаем HTML код основного содержимого страницы
	content.Content, err = tasksHTMLcontent(authorizedUser)
	if err != nil {
		log.Println("Не удалось отрисовать основное содержимое страницы с заданиями: ", err)
		ErrorPage(w, "Ошибка страницы", http.StatusInternalServerError)
		return
	}

	templ, err := template.ParseFiles("src/templates/page.html")
	if err != nil {
		log.Fatal("Не удалось открыть основной шаблон страницы: ", err)
	}
	// Рендерим страницу и отправляем её клиенту
	err = templ.Execute(*w, content)
	if err != nil {
		// если не удалось отрисовать страницу по шаблону
		ErrorPage(w, "Ошибка страницы", http.StatusInternalServerError)
		return
	}
}

// Отрисовка страницы с заданием
func Task(w *http.ResponseWriter, path string, authorizedUser *model.UserInfo) {
	var content mainPageContent
	var err error
	// Название страницы в метаданных
	content.Name = "Список заданий"
	// Получаем HTML код главного меню
	content.Menu, err = menuHTMLcontent(authorizedUser, "tasks")
	if err != nil {
		log.Println("Не удалось отрисовать главное меню: ", err)
		ErrorPage(w, "Ошибка страницы", http.StatusInternalServerError)
		return
	}
	// Получаем HTML код основного содержимого страницы
	content.Content, err = getHTMLsegment("src/templates/tasks/"+path+".html", nil)
	if err != nil {
		log.Println("Не удалось отрисовать основное содержимое страницы задания: ", err)
		ErrorPage(w, "Ошибка страницы", http.StatusInternalServerError)
		return
	}

	templ, err := template.ParseFiles("src/templates/page.html")
	if err != nil {
		log.Fatal("Не удалось открыть основной шаблон страницы: ", err)
	}
	// Рендерим страницу и отправляем её клиенту
	err = templ.Execute(*w, content)
	if err != nil {
		// если не удалось отрисовать страницу по шаблону
		ErrorPage(w, "Ошибка страницы", http.StatusInternalServerError)
		return
	}
}

// Отрисовка страницы с результатами
func Results(w *http.ResponseWriter, authorizedUser *model.UserInfo) {
	var content mainPageContent
	var err error
	// Название страницы в метаданных
	content.Name = "Список заданий"
	// Получаем HTML код главного меню
	content.Menu, err = menuHTMLcontent(authorizedUser, "results")
	if err != nil {
		log.Println("Не удалось отрисовать главное меню: ", err)
		ErrorPage(w, "Ошибка страницы", http.StatusInternalServerError)
		return
	}
	// Получаем HTML код основного содержимого страницы
	content.Content, err = resultsHTMLcontent(authorizedUser)
	if err != nil {
		log.Println("Не удалось отрисовать основное содержимое страницы с результатами: ", err)
		ErrorPage(w, "Ошибка страницы", http.StatusInternalServerError)
		return
	}

	templ, err := template.ParseFiles("src/templates/page.html")
	if err != nil {
		log.Fatal("Не удалось открыть основной шаблон страницы: ", err)
	}
	// Рендерим страницу и отправляем её клиенту
	err = templ.Execute(*w, content)
	if err != nil {
		// если не удалось отрисовать страницу по шаблону
		ErrorPage(w, "Ошибка страницы", http.StatusInternalServerError)
		return
	}
}

// Отрисовка страницы с литературой
func Literature(w *http.ResponseWriter, authorizedUser *model.UserInfo) {
	var content mainPageContent
	var err error
	// Название страницы в метаданных
	content.Name = "Список заданий"
	// Получаем HTML код главного меню
	content.Menu, err = menuHTMLcontent(authorizedUser, "literature")
	if err != nil {
		log.Println("Не удалось отрисовать главное меню: ", err)
		ErrorPage(w, "Ошибка страницы", http.StatusInternalServerError)
		return
	}
	// Получаем HTML код основного содержимого страницы
	content.Content, err = literatureHTMLcontent(authorizedUser)
	if err != nil {
		log.Println("Не удалось отрисовать основное содержимое страницы литературы: ", err)
		ErrorPage(w, "Ошибка страницы", http.StatusInternalServerError)
		return
	}

	templ, err := template.ParseFiles("src/templates/page.html")
	if err != nil {
		log.Fatal("Не удалось открыть основной шаблон страницы: ", err)
	}
	// Рендерим страницу и отправляем её клиенту
	err = templ.Execute(*w, content)
	if err != nil {
		// если не удалось отрисовать страницу по шаблону
		ErrorPage(w, "Ошибка страницы", http.StatusInternalServerError)
		return
	}
}
