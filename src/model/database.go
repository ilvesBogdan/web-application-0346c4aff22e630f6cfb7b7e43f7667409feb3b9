package model

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var database *sqlx.DB

// Имена всех таблицы в базе данных.
var tables tablesNames

// Создание подключения к базе данных PostgreSQL.
func ConnectDataBase(login string, password string, dataBaseName string) {
	var err error
	database, err = sqlx.Connect(
		"postgres",
		fmt.Sprintf("sslmode=disable user=%v password=%v dbname=%v", login, password, dataBaseName))
	if err != nil {
		log.Fatal("Не удалось подключиться к базе данных: ", err)
	}
	err = database.Ping()
	if err != nil {
		log.Fatal("Нет пинга до базы данных: ", err)
	}
	log.Print("Подключено к базе данных. Номер подключения: ", database.Stats().Idle)
	// Заполнение массива Tables названиями таблиц
	rows, err := database.Query(`
		SELECT table_name FROM information_schema.tables
		WHERE table_schema NOT IN ('information_schema','pg_catalog')`)
	if err != nil {
		log.Fatal("Не удалось выполнить запрос о таблицах в базе данных: ", err)
	}
	// Имена таблиц
	var requestTablesNames []string
	for rows.Next() {
		var tableName string
		err = rows.Scan(&tableName)
		requestTablesNames = append(requestTablesNames, tableName)
	}
	// Проверка на наличие ошибок в цикле выше.
	err = rows.Err()
	if err != nil {
		log.Fatal("Ошибка при парсинге запроса о таблицах в базе данных: ", err)
	}
	// Заполнение глобальной структуры с именами таблиц
	fillInTheTableStructure(requestTablesNames)
}

// Закрывает соединение с базой данных.
func CloseDataBase() {
	err := database.Close()
	if err != nil {
		log.Fatal("Не удалось закрыть соединение с базой данных: ", err)
	}
}

// Проверяет наличие ключа в таблице Sessions.
//
// Возвращает true, если найден такой же ключ.
func RequestWhereSessionKeys(key string) bool {
	var count int8
	err := database.Get(&count,
		fmt.Sprintf("SELECT COUNT(Secret_key) FROM Sessions WHERE Secret_key='%v'", key))
	if err != nil {
		log.Fatal(
			fmt.Sprintf(
				"Не получилось выполнить запрос на проверку наличия ключа '%v' в sessions: ",
				key),
			err)
	}
	return count != 0
}

// Записывает avatar_id в базу данных пользователю по его логину.
func SetAvatarIdWhereLogin(avatarId int, login string) {
	_, err := database.Exec(fmt.Sprintf(
		"UPDATE %v SET avatar_id='%v' WHERE login='%v'",
		tables.Users, avatarId, login,
	))
	if err != nil {
		log.Print("Не удалось записать 'avatar_id' в базу данных: ", err)
	}
}

// Заполнение структуры данными из базы данных.
//
// По условию WHERE, условие передавать как аргумент.
//
// Прмер: u.GetDataFromDB(model.UserInfo{Id: 5})
func (user *UserInfo) GetDataFromDB(where string) error {
	err := database.Get(user,
		fmt.Sprintf("SELECT %v.* FROM %v %v", tables.Users, tables.Users, where))
	if err != nil {
		log.Print("Пользователь не найден в базе данных: ", err)
		return err
	}
	user.trimSpace()
	return nil
}

// Получение id пользователя, по его логину.
func GetIdByLogin(login string) (int, error) {
	var id int
	err := database.Get(&id, fmt.Sprintf(
		"SELECT id FROM %v WHERE login='%v'", tables.Users, login))
	if err != nil {
		return 0, err
	}
	return id, nil
}

// Записывает структуру в базу данных.
func (sessionStruct *SessionsFullInfo) PushToDataBase() error {
	query, err := qeuryInsert(sessionStruct, tables.Sessions)
	if err != nil {
		return err
	}
	_, err = database.Exec(query)
	if err != nil {
		return err
	}
	return err
}

// Записывает структуру в базу данных.
func (userInfoStruct *UserInfo) PushToDataBase() error {
	query, err := qeuryInsert(userInfoStruct, tables.Users, "id")
	if err != nil {
		return err
	}
	_, err = database.Exec(query)
	if err != nil {
		return err
	}
	return err
}

// Записывает структуру в базу данных.
func (messageStruct *Message) PushToDataBase() error {
	query, err := qeuryInsert(messageStruct, tables.Message, "text")
	if err != nil {
		return err
	}
	_, err = database.Exec(query)
	if err != nil {
		return err
	}
	return err
}

// Получить структуру сообщения по условию из базы данных.
func GetMessageFromDB(where string) (Message, error) {
	var message Message
	err := database.Get(&message, fmt.Sprintf(
		// Только за один этот SQL запрос Богдану нужно поставить отлично по ГПО
		`SELECT id, sender, recipient, date, time, is_read,
		(
			WITH RECURSIVE temp(id, text, next)
			AS (SELECT t1.id, t1.text, t1.next, CAST (t1.text AS varchar(255)) AS fullText
			FROM %v t1 WHERE t1.id = (SELECT text_id FROM %v %v)
			UNION
			SELECT t2.id, t2.text, t2.next, CAST (temp.fullText || t2.text as varchar(255))
			FROM %v t2 JOIN temp ON (temp.next = t2.id))
			SELECT MAX(fulltext) FROM temp
		) as text
	FROM %v %v`,
		tables.messageText, tables.Message, where,
		tables.messageText, tables.Message, where))
	if err != nil {
		log.Print(fmt.Sprintf("Ошибка получения личного сообщения по условию “%v” из базы данных: ",
			where), err)
		return Message{}, err
	}
	return message, nil
}

// Стирание данных о сессии из базы данных по передаваемому ключу.
func RemoveSessionOnDataBase(key interface{}) {
	switch key.(type) {
	case nil:
		return
	case string:
		_, err := database.Exec(
			fmt.Sprintf("DELETE FROM %v WHERE secret_key = '%v'", tables.Sessions, key.(string)))
		if err != nil {
			log.Print("Пользователь не найден в базе данных: ", err)
			return
		}
	default:
		log.Println("В функцию RemoveSessionOnDataBase попал ключь с неизвестным типом данных")
	}
}

// Получить задания по id копетенции
func GetFullTasksByComp() (map[string][]Task, []string) {
	result := make(map[string][]Task)
	var keys []string
	var competences []Competence
	var tasks []Task
	err := database.Select(&competences, fmt.Sprintf("SELECT * FROM %v", tables.Competence))
	if err != nil {
		log.Print("Ошибка получения компетенций из базы данных: ", err)
		return result, keys
	}
	err = database.Select(&tasks, fmt.Sprintf("SELECT * FROM %v", tables.Tasks))
	if err != nil {
		log.Print("Ошибка получения заданий из базы данных: ", err)
		return result, keys
	}
	for _, task := range tasks {
		for _, competence := range competences {
			if task.CompetenceId == competence.Id {
				result[competence.Name] = append(result[competence.Name], task)
			}
		}
	}
	for _, competence := range competences {
		keys = append(keys, competence.Name)
	}
	return result, keys
}
