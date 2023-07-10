package controller

import (
	"crypto/rand"
	"log"
	"math/big"
	"net/http"
	"os"

	"web-application/src/model"

	"github.com/gorilla/sessions"
)

// Требуется реализация

const authCookieName = "auth"

var store = sessions.NewCookieStore([]byte(os.Getenv("session_key")))

// Проверка на наличие авторизации.
//
// Возвращает UserInfo авторизированного пользователя.
func sessionCheckSecret(w *http.ResponseWriter, r *http.Request) model.UserInfo {
	var user model.UserInfo
	var key interface{}
	session, err := store.Get(r, authCookieName)
	if err != nil {
		log.Println("Не удалось получить ключ с клиента: ", err)
		goto loginOut
	}
	key = session.Values[authCookieName]
	if key == nil {
		// Если неудалось получить ключ
		return model.UserInfo{Id: 0}
	}
	err = user.GetDataFromDB(model.QuryWhereSessionForUserInfo(key.(string)))
	if err != nil {
		// Если не получилось идентифицировать пользователя по ключу сессии
		goto loginOut
	}
	return user
loginOut:
	sessionLogOut(*w, r)
	return model.UserInfo{Id: 0}
}

// Авторизация. Запись ключа в куки и в базу данных
func sessionLogin(w *http.ResponseWriter, r *http.Request, user *model.UserInfo, request *model.ApiData) error {
	session, _ := store.Get(r, authCookieName)

	// Реализовать генерацию ключа
	const lenghtKey = 32
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-+=[]()#@|<>"
	buffKey := make([]byte, lenghtKey)
newGenerate:
	for i := 0; i < lenghtKey; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return err
		}
		buffKey[i] = letters[num.Int64()]
	}
	key := string(buffKey)
	// Проверка наличия ключа в базе данных
	if model.RequestWhereSessionKeys(key) {
		// Если ключ найден, повторяем генерацию
		goto newGenerate
	}
	var sessionData model.SessionsFullInfo
	// Запись ключа в кукисы
	session.Values[authCookieName] = key
	// Запись данных в структуру, для передачи её в БД
	lastDate := request.DateTimeOfReceipt.Format("2006-01-02")
	sessionData.UserID = user.Id
	sessionData.Name = request.SessionName
	sessionData.SecretKey = key
	sessionData.CreationDate = lastDate
	sessionData.LastDate = lastDate
	sessionData.LastTime = request.DateTimeOfReceipt.Format("15:04:05")
	// Оправка кукисов на клиент
	err := session.Save(r, *w)
	if err != nil {
		return err
	}
	// Запись данных о сессии в базу данных
	err = sessionData.PushToDataBase()
	if err != nil {
		return err
	}
	return nil
}

// Стирание куки авторизации на клиенте и из БД.
func sessionLogOut(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, authCookieName)

	// Отчиска записи о сессии из базы данных
	go model.RemoveSessionOnDataBase(session.Values[authCookieName])

	// Отменить аутентификацию пользователей
	session.Values[authCookieName] = nil
	session.Save(r, w)
}
