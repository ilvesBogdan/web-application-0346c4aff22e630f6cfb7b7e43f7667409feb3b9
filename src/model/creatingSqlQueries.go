package model

import (
	"fmt"
	"log"
	"strconv"
	"time"

	attr "github.com/ssrathi/go-attr"
)

// Перебирает поля передаваемой структуры, и передает в передаваемую функцию тег поля “db”,
// и значение приведенное к string.
//
// Передаваемая функция исполнится столько раз сколько полей в структуре с тегом “db”.
func queryGenerator(dataStruct interface{}, process func(key string, value string)) error {
	structList, err := attr.Tags(dataStruct, "db")
	// Если не получилось добыть из структуры теги по атрибуту “db”
	if err != nil {
		log.Fatal("Ошибка в функции “queryGenerator”, при получении полей структуры с тегом “db”: ", err)
		return err
	}
	var value string
	var rawValue interface{}
	// Перебор всех свойств структуры
	for structKey, key := range structList {
		rawValue, err = attr.GetValue(dataStruct, structKey)
		if err != nil {
			log.Fatal("Ошибка в функции “queryGenerator”, при получении значения поля: ", err)
			return err
		}
		switch typeValue := rawValue.(type) {
		case int:
			value = strconv.Itoa(typeValue)
		case int8:
			value = strconv.Itoa(int(typeValue))
		case int16:
			value = strconv.Itoa(int(typeValue))
		case int32:
			value = strconv.Itoa(int(typeValue))
		case int64:
			value = strconv.FormatInt(typeValue, 10)
		case string:
			value = typeValue
		case time.Time:
			value = typeValue.Format("2006-01-02")
		default:
			value = fmt.Sprintf("%v", rawValue)
			log.Println(fmt.Sprintf(
				"В функции “queryGenerator” не не нашлось преобразования к string, значение было преобразовано как %v.",
				value))
		}
		// выполнение переданной функции
		if 0 < len(value) {
			process(key, value)
			value = ""
		}
	}
	return nil
}

// Создает строку с SQL командой INSERT, на добавление данных в таблицу.
//
// Первым аргументом принимаются заполненные структуры с тэгом “db”.
//
// А вторым название таблицы, куда должны быть вставлены данные.
//
// Вывод: "INSERT INTO NameTable(arg1, arg2) VALUES(value1, value2)"
func qeuryInsert(dataStruct interface{}, table string, excludedColumns ...string) (string, error) {
	columns, values := "", ""
	var function func(key string, value string)
	if 0 < len(excludedColumns) {
		function = func(key string, value string) {
			for _, ec := range excludedColumns {
				if ec != key {
					columns += key + ", "
					values += "'" + value + "', "
				}
			}
		}
	} else {
		function = func(key string, value string) {
			columns += key + ", "
			values += "'" + value + "', "
		}
	}
	err := queryGenerator(dataStruct, function)
	if err != nil {
		log.Fatal("Ошибка генератора в функции “qeuryInsert”")
		return "", err
	}
	if len(columns) == 0 || len(values) == 0 {
		return "", qeuryCreateError
	}
	// Обрезка двух последних символов в конце каждой строки - “, ”
	columns, values = columns[:len(columns)-2], values[:len(values)-2]
	return fmt.Sprintf("INSERT INTO %v(%v) VALUES(%v)", table, columns, values), nil
}

// Генерирует условие для поиска пользователя в базе данных, при попытке его аутентификации.
//
// Вывод: "WHERE login='test_login' AND password='qwerty123'"
func qeuryLogInWhere(login *string, password *string) string {
	return fmt.Sprintf("WHERE login='%v' AND password='%v'", *login, *password)
}

// Генерирует условие для поиска пользователя в базе данных, при открытии его профиля.
//
// Вывод: "WHERE login='test_login'"
func QuryUrlPathProfileWhere(urlPath *string) string {
	login := (*urlPath)[2:]
	return fmt.Sprintf("WHERE login='%v'", login)
}

func QuryWhereSessionForUserInfo(sessionKey string) string {
	return fmt.Sprintf(
		"JOIN %v ON %v.id = %v.user_id WHERE %v.secret_key = '%v'",
		tables.Sessions, tables.Users, tables.Sessions, tables.Sessions, sessionKey)
}
