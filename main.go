package main

import (
	"os"

	"web-application/src/controller"
	"web-application/src/model"
)

func main() {
	// Инициализация подключения к базе данных
	model.ConnectDataBase(
		os.Getenv("db_user"),
		os.Getenv("db_password"),
		os.Getenv("main_db_name"),
	)
	// При завершении работы приложения соединение с базой данных будет закрыто
	defer model.CloseDataBase()
	controller.HandleRequest()
}
