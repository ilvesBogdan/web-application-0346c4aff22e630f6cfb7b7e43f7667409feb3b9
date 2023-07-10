package controller

import (
	"log"
	"net/http"
	"strings"
	"web-application/src/view"
)

// Выполняет действия в зависимости от url запроса.
func HandleRequest() {

	// Адреса страниц.
	http.HandleFunc("/", mainPathHandler)
	http.HandleFunc("/api/", listenApiRequests)
	http.HandleFunc("/logout/", loginOut)
	http.HandleFunc("/favicon.ico", faviconHandler)
	http.HandleFunc("/registration/", registrationHandler)

	http.HandleFunc("/tasks/", tasksHandler)
	http.HandleFunc("/results/", resultsHandler)
	http.HandleFunc("/literature/", literatureHandler)

	// Присвоение url ссылок файлам из папки "src/static"
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("src/static"))))

	log.Println("Сервер запущен.")

	// Порт работы web сервера.
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("Ошибка запуска сервера: ", err)
	}
}

// **********************************************************
// Тут можно посмотреть все статус коды из пакета http
// https://go.dev/src/net/http/status.go
// **********************************************************

// Обработка динамических адресов страниц
func mainPathHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		view.ErrorPage(&w, "Неверный метод запроса", http.StatusMethodNotAllowed)
		return
	}
	userSession := sessionCheckSecret(&w, r)
	if r.URL.Path == "/" {
		if 1 > userSession.Id {
			view.MainPage(&w)
		} else {
			http.Redirect(w, r, "/@"+userSession.Login, http.StatusSeeOther)
		}
		return
	}
	// Первый символ после слеша
	switch r.URL.Path[1] {
	case 64:
		// Если после слеша идет символ "@"
		if 1 > userSession.Id {
			// Если пользователь не авторизован, перенапрявляет
			// на главную страницу
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		if 4 > len(r.URL.Path) {
			// если адрес сильно короткий
			view.ErrorPage(&w, "Пользователь не найден", http.StatusNotFound)
			return
		}
		view.Profile(&w, &r.URL.Path, &userSession)
		return
	}
	view.ErrorPage(&w, "Cтраница не найдена", http.StatusNotFound)
}

func registrationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		view.ErrorPage(&w, "Неверный метод запроса",
			http.StatusMethodNotAllowed)
		return
	}
	userSession := sessionCheckSecret(&w, r)
	if 1 < userSession.Id {
		// Если пользователь авторизован
		view.ErrorPage(&w, "Страница недоступна авторизованным пользователям.",
			http.StatusForbidden)
		return
	}
	view.Registration(&w, r)
}

// Обработка запроса к странице с заданиями
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		view.ErrorPage(&w, "Неверный метод запроса", http.StatusMethodNotAllowed)
		return
	}
	userSession := sessionCheckSecret(&w, r)
	if 1 > userSession.Id {
		// Если пользователь не авторизован, перенапрявляет
		// на главную страницу
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	path := strings.Split(r.URL.Path, "/")
	switch len(path) {
	case 3:
		view.Tasks(&w, &userSession)
	case 4:
		view.Task(&w, path[2], &userSession)
	default:
		view.ErrorPage(&w, "Страница не найдена", http.StatusNotFound)
	}
}

// Обработка запроса к странице с результатами
func resultsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		view.ErrorPage(&w, "Неверный метод запроса", http.StatusMethodNotAllowed)
		return
	}
	userSession := sessionCheckSecret(&w, r)
	if 1 > userSession.Id {
		// Если пользователь не авторизован, перенапрявляет
		// на главную страницу
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	view.Results(&w, &userSession)
}

// Обработка запроса к странице с литературой
func literatureHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		view.ErrorPage(&w, "Неверный метод запроса", http.StatusMethodNotAllowed)
		return
	}
	userSession := sessionCheckSecret(&w, r)
	if 1 > userSession.Id {
		// Если пользователь не авторизован, перенапрявляет
		// на главную страницу
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	view.Literature(&w, &userSession)
}

// Путь к иконке вкладок страниц.
func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "src/static/images/icon.ico")
}

// Выход из профиля
func loginOut(w http.ResponseWriter, r *http.Request) {
	sessionLogOut(w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
