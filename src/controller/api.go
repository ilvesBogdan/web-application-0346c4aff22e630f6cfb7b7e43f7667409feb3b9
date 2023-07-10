package controller

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"web-application/src/model"
)

// Обработка данных полученных методом POST.
//
// Формирование и оправка ответа обратно в браузер.
func listenApiRequests(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}
	if r.Header.Get("content-type") == "application/json" {
		request := model.GetStructApiData(r)
		var statusCode int
		var response string
		switch request.TypeData {
		case "login":
			response, statusCode = loginIn(&w, r, &request)
			if 299 < statusCode {
				http.Error(w, response, statusCode)
				return
			}
			w.Write([]byte(response))
			return
		case "registration":
			response, statusCode = registration(&w, r, &request)
			if 299 < statusCode {
				http.Error(w, response, statusCode)
				return
			}
			w.Write([]byte(response))
			return
		}
	} else {
		log.Println("Content-Type не является application/json")
	}
}

func loginIn(w *http.ResponseWriter, r *http.Request, request *model.ApiData) (string, int) {
	var authData model.AuthData
	var user model.UserInfo
	var err error
	authData, err = model.GetStructAuthDataOfJsonString(r, &request.Data)
	if err != nil {
		// Ошибка парсинга данных
		return `{"error":"object.data parse json"}`, http.StatusInternalServerError
	}
	err = user.GetDataFromDB(authData.QeuryLogInWhere())
	if err != nil {
		// Пользователь не найден в базе данных
		return `{"error":"incorrect data"}`, http.StatusOK
	}
	err = sessionLogin(w, r, &user, request)
	if err != nil {
		log.Println("Ошибка авторизации: ", err)
	}
	// Оправляем JSON строку, где нет ошибок, и адрес страницы
	// на которую нужно перенаправить, после успешной авторизации
	return `{"error":"","redirect":"@` + user.Login + `"}`, http.StatusOK
}

func registration(w *http.ResponseWriter, r *http.Request, request *model.ApiData) (string, int) {
	user, avatar, err := model.GetStructUserInfoOfJsonString(&request.Data)
	if err != nil {
		return fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError
	}
	// Очищаем строку, для экономии памяти.
	request.Data = ""
	user.LoginToLower()

	// Конвертируем строку base64 в объект изображения и получаем формат
	if "none" != avatar {
		go func(login string, base64img string) {
			file, err := os.Create(
				fmt.Sprintf("src/static/images/profile/%v.base64", login))
			if err != nil {
				log.Println(
					"Не удалось создать временный файл для аватарки: ", err)
				return
			}
			file.WriteString(base64img)
			file.Close()
			base64img = ""
			cmd := exec.Command(
				"python",
				"src\\writeProfileImage.py",
				fmt.Sprintf("%v.base64", login))
			raw, err := cmd.Output()
			if err != nil {
				log.Println(
					"Не удалось получить AvatarId из Python скрипта: ", err)
				return
			}
			avatarId, err := strconv.Atoi(strings.TrimSpace(string(raw[:])))
			if err != nil {
				log.Panicln("При сохранение аватарки произошла ошибка: ", string(raw[:]))
				return
			}
			model.SetAvatarIdWhereLogin(avatarId, login)
		}(user.Login, avatar)
	}

	// Очищаем строку, для экономии.
	avatar = ""

	user.RegistrationDate = request.DateTimeOfReceipt
	err = user.PushToDataBase()
	if err != nil {
		log.Println(err)
		return fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError
	}
	// Получаем id который присвоила записи база данных
	user.Id, err = model.GetIdByLogin(user.Login)
	if err != nil {
		log.Println(err)
		return `{"error":"Не удалось получить id записи"}`, http.StatusInternalServerError
	}
	// Открываем новую сессию для пользователя.
	err = sessionLogin(w, r, &user, request)
	if err != nil {
		log.Println(err)
		return `{"error":"Авторизация не удалась"}`, http.StatusInternalServerError
	}
	// Отправляем в браузер результат работы ввиде json строки
	return `{"error":"","redirect":"/@` + user.Login + `"}`, http.StatusOK
}
