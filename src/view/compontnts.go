package view

import (
	"bytes"
	"html/template"
	"log"
	"web-application/src/model"
)

// Возвращает строку с содержимым HTML шаблона с подставленными значениями
func getHTMLsegment(pathToHtmlTemplate string, content any) (template.HTML, error) {
	tmpl, err := template.ParseFiles(pathToHtmlTemplate)
	if err != nil {
		log.Fatal("Не удалось открыть шаблон профиля, по пути: “", pathToHtmlTemplate,
			"” Ошибка: ", err)
	}
	var htmlBytes bytes.Buffer
	err = tmpl.Execute(&htmlBytes, content)
	if err != nil {
		// если не удалось отрисовать страницу по шаблону
		return "", err
	}
	return template.HTML(htmlBytes.String()), nil
}

// Возвращает HTML код компонента главного меню страницы
//
// Первый аргумент (authorizedUser) – ссылка на структуру с данными авторизированного пользователя.
//
// Второй аргумент (button_name) – название из HTML шаблона кнопки меню, для её подсветки.
func menuHTMLcontent(authorizedUser *model.UserInfo, button_name string) (template.HTML, error) {
	content := make(map[string]any)

	// Тут заносим данные которые нужно передать в HTML шаблон
	content["user_login"] = authorizedUser.Login
	content["user_char"] = authorizedUser.FullName[:2]
	content["user_avatar"] = int(authorizedUser.AvatarId)
	content["button_highlighting"] = button_name
	// -------------------------------------------------------

	html, err := getHTMLsegment("src/templates/partsOfPages/mainMenu.html", content)
	if err != nil {
		return "Ошибка страницы", err
	}
	return html, nil
}

// Возвращает HTML код части страницы ответственной за личный кабинет пользователя
func profileHTMLcontent(urlPath *string) (template.HTML, error) {
	var html template.HTML
	var user model.UserInfo
	err := user.GetDataFromDB(model.QuryUrlPathProfileWhere(urlPath))
	if err != nil {
		// если пользователь не найден в базе данных
		return "Пользователь не найден", err
	}
	content := make(map[string]any)

	// Тут заносим данные которые нужно передать в HTML шаблон
	content["fio"] = user.FullName
	content["avatar"] = int(user.AvatarId)
	// -------------------------------------------------------

	html, err = getHTMLsegment("src/templates/partsOfPages/profile.html", content)
	if err != nil {
		return "Ошибка страницы", err
	}
	return html, nil
}

// Возвращает HTML код компонента страницы с заданиями
func tasksHTMLcontent(authorizedUser *model.UserInfo) (template.HTML, error) {
	content := make(map[string]any)

	// Тут заносим данные которые нужно передать в HTML шаблон
	content["data"], content["data_keys"] = model.GetFullTasksByComp()
	// -------------------------------------------------------

	html, err := getHTMLsegment("src/templates/partsOfPages/tasks.html", content)
	if err != nil {
		return "Ошибка страницы", err
	}
	return html, nil
}

// Возвращает HTML код компонента страницы с результатами
func resultsHTMLcontent(authorizedUser *model.UserInfo) (template.HTML, error) {
	content := make(map[string]any)

	// Тут заносим данные которые нужно передать в HTML шаблон
	content["test"] = "test"
	// -------------------------------------------------------

	html, err := getHTMLsegment("src/templates/partsOfPages/results.html", content)
	if err != nil {
		return "Ошибка страницы", err
	}
	return html, nil
}

// Возвращает HTML код компонента страницы с литературой
func literatureHTMLcontent(authorizedUser *model.UserInfo) (template.HTML, error) {
	content := make(map[string]any)

	// Тут заносим данные которые нужно передать в HTML шаблон
	content["test"] = "test"
	// -------------------------------------------------------

	html, err := getHTMLsegment("src/templates/partsOfPages/literature.html", content)
	if err != nil {
		return "Ошибка страницы", err
	}
	return html, nil
}
