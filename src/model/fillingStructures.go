package model

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	attr "github.com/ssrathi/go-attr"
)

// Заполняет структуру названиями таблиц из базы данных.
//
// А так же проверяет на наличие таблицы из структуры в безе данных.
func fillInTheTableStructure(tablesNames []string) {
	structList, err := attr.Tags(tables, "table")
	// Если не получилось добыть из структуры теги по атрибуту “db”
	if err != nil {
		log.Fatal("Ошибка в функции “fillInTheTableStructure”, при получении полей структуры с тегом “table”: ", err)
	}
	notOk := true
	// Перебор всех свойств структуры
	for structKey, key := range structList {
		for _, i := range tablesNames {
			if i == key {
				err = attr.SetValue(&tables, structKey, i)
				if err != nil {
					log.Fatal("Ошибка в функции “fillInTheTableStructure”, при записи данных: ", err)
				}
				notOk = false
				break
			}
		}
		if notOk {
			if 0 < len(key) {
				log.Fatal(fmt.Sprintf("Таблица “%v” не найдена в базе данных.", key))
			}
			log.Fatal(fmt.Sprintf("В структуре таблиц, у поля “%v”, отсутствует тег “table”.", structKey))
		}
		notOk = true
	}
}

// Заполняет структуру model.ApiData данными полученными из POST запроса.
func GetStructApiData(r *http.Request) ApiData {
	var request ApiData
	// Текущее время
	dateTime := time.Now()
	rawBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Ошибка при чтении тела запроса: ", err)
		panic(err)
	}
	err = json.Unmarshal(rawBody, &request)
	if err != nil {
		log.Println("Ошибка парсинга json из тела запроса и записи данных в структуру ApiData: ", err)
		panic(err)
	}
	request.SecretKey = "SecretKey"
	request.SessionName = "SessionName"
	// Ниже не время и дата, которые будут записаны, а это шаблоны
	request.DateTimeOfReceipt = dateTime
	return request
}

// Возвращает заполненую структуру model.AuthData на основе данных из json строки.
//
// Не заполняет поле "SecretKey"
func GetStructAuthDataOfJsonString(r *http.Request, jsonString *string) (AuthData, error) {
	var request AuthData
	// Время входа в приложение
	request.Date = time.Now()
	raw := []byte(*jsonString)
	err := json.Unmarshal(raw, &request)
	if err != nil {
		log.Println("Ошибка парсинга json и записи данных в структуру AuthData: ", err)
		return AuthData{}, err
	}
	// Имя сессии
	request.ClientName = fmt.Sprintf(`ip: "%v" user agent: "%v"`, r.RemoteAddr, r.UserAgent())
	return request, nil
}

// Возвращает заполненую структуру model.UserInfo на основе данных из json строки.
func GetStructUserInfoOfJsonString(jsonString *string) (UserInfo, string, error) {
	var user UserInfo
	var base64Image string
	raw := []byte(*jsonString)
	err := json.Unmarshal(raw, &user)
	if err != nil {
		log.Println("Ошибка парсинга json и записи данных в структуру UserInfo: ", err)
		return UserInfo{}, "", err
	}
	// Получение строки типа base64 с переданной аватаркой.
	base64Image, err = func(raw []byte) (string, error) {
		data := make(map[string]interface{})
		err := json.Unmarshal(raw, &data)
		if err != nil {
			return "", err
		}
		return data["avatar_base64"].(string), nil
	}(raw)
	if err != nil {
		log.Println("Ошибка парсинга аватарки из json: ", err)
		return UserInfo{}, "", err
	}
	return user, base64Image, nil
}
