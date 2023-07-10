package model

import (
	"strings"
	"time"
)

// Таблицы находящиеся в базе данных.
//
// При создании новой таблицы, нужно добавлять её в эту структуру.
type tablesNames struct {
	// Структура поля:
	// Название таблицы для backend’a | всегда string | `table:"название таблицы в базе данных"`
	Users       string `table:"users"`
	Sessions    string `table:"sessions"`
	Competence  string `table:"competence"`
	Tasks       string `table:"tasks"`
	Message     string `table:"message"`
	messageText string `table:"message_text"`
}

// Вся информация о пользователе хранимая в базе данных.
type UserInfo struct {
	Id               int       `db:"id"`
	Login            string    `db:"login"       json:"login"`
	Password         string    `db:"password"    json:"password"`
	FullName         string    `db:"full_name"   json:"full_name"`
	Email            string    `db:"email"       json:"email"`
	Faculty          string    `db:"faculty"     json:"faculty"`
	Group            string    `db:"study_group" json:"group"`
	AvatarId         int       `db:"avatar_id"`
	RegistrationDate time.Time `db:"registration_date"`
}

// Убирает лишние пробелы в конце строк.
func (u *UserInfo) trimSpace() {
	u.Login = strings.TrimSpace(u.Login)
	u.Password = strings.TrimSpace(u.Password)
	u.FullName = strings.TrimSpace(u.FullName)
	u.Email = strings.TrimSpace(u.Email)
	u.Faculty = strings.TrimSpace(u.Faculty)
	u.Group = strings.TrimSpace(u.Group)
}

// Привести логин к нижнему регистру, и удалить из него пробелы.
func (u *UserInfo) LoginToLower() {
	u.Login = strings.ToLower(strings.ReplaceAll(u.Login, " ", ""))
}

// Данные полученные через api методом POST из браузера.
type ApiData struct {
	TypeData          string `json:"type"`
	SessionName       string
	SecretKey         string
	DateTimeOfReceipt time.Time
	Data              string `json:"data"`
}

// Данные полученные из браузера при попытке авторизации.
type AuthData struct {
	Login      string `json:"login"`
	Password   string `json:"password"`
	ClientName string
	SecretKey  string
	Date       time.Time
}

// Генерирует условие для поиска пользователя в базе данных, при попытке его аутентификации.
//
// Вывод: "WHERE login='test_login' AND password='qwerty123'"
func (auth *AuthData) QeuryLogInWhere() string {
	auth.Login = strings.ToLower(strings.ReplaceAll(auth.Login, " ", ""))
	// Добавить хжширование пароля
	return qeuryLogInWhere(&auth.Login, &auth.Password)
}

// Информация о сессии хранимая в базе данных.
type SessionsFullInfo struct {
	UserID       int    `db:"user_id"`
	Name         string `db:"name"`
	SecretKey    string `db:"secret_key"`
	CreationDate string `db:"creation_date"`
	LastDate     string `db:"last_date"`
	LastTime     string `db:"last_time"`
}

func (s *SessionsFullInfo) trimSpace() {
	s.Name = strings.TrimSpace(s.Name)
	s.SecretKey = strings.TrimSpace(s.SecretKey)
}

// Компетенция
type Competence struct {
	Id        int    `db:"id"`
	Name      string `db:"name"`
	Chort_nam string `db:"chort_name"`
}

// Информация о задании
type Task struct {
	Id           int    `db:"id"`
	Name         string `db:"name"`
	Description  string `db:"description"`
	CompetenceId int    `db:"competence_id"`
	HtmlNum      int    `db:"html_num"`
	Mark         int    `db:"mark"`
}

// Информация о личном сообщении
type Message struct {
	Id        int    `db:"id"`
	Sender    int    `db:"sender"`
	Recipient int    `db:"recipient"`
	Date      string `db:"date"`
	Time      string `db:"time"`
	Is_read   bool   `db:"is_read"`
	Text      string `db:"text"`
}
