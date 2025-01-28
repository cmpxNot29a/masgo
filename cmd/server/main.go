package main

import (
	"fmt"
	"net/http"

	"github.com/cmpxNot29a/masgo/internal/handlers"
	"github.com/cmpxNot29a/masgo/internal/storage"
)

// Константы для путей и сообщений об ошибках.
const (
	serverAddress = "[::]:8080" // Адрес и порт сервера.
	pathUpdate    = "/update/"  // Путь для обновления метрик.

	msgStartingServer   = "Starting server on "     // Сообщение о запуске сервера.
	errorStartingServer = "Error starting server: " // Сообщение об ошибке при запуске сервера.
)

func main() {
	// Создаем хранилище метрик.
	storage := storage.NewMemStorage()

	// Создаем новый HTTP-мультиплексор.
	mux := http.NewServeMux()
	// Регистрируем обработчик `UpdateHandler` для пути `/update/`.
	mux.HandleFunc(pathUpdate, handlers.UpdateHandler(storage))

	// Выводим сообщение о запуске сервера.
	fmt.Println(msgStartingServer, serverAddress)
	// Запускаем HTTP-сервер, слушающий на указанном адресе и порту.
	err := http.ListenAndServe(serverAddress, mux)
	if err != nil {
		// В случае ошибки выводим сообщение об ошибке.
		fmt.Println(errorStartingServer, err)
	}
}
