package model

import "errors"

var (
	qeuryCreateError = errors.New("Ошибка генерации запроса, возможно в сруктуре отсутвуют поля с тегом 'db'.")
)
