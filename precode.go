package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Ниже напишите обработчики для каждого эндпоинта
// ...

//Обработчик для получения всех задач
func getTasks(res http.ResponseWriter, req *http.Request) {

	tasks, err := json.Marshal(tasks)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(tasks)
}

//Обработчик для отправки задачи на сервер
func postTask(res http.ResponseWriter, req *http.Request) {
	var task Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(buf.Bytes(), &task)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	tasks[task.ID] = task

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
}

//Обработчик для получения задачи по ID
func getTask(res http.ResponseWriter, req *http.Request) {		
	id := chi.URLParam(req, "id")

	if _, ok := tasks[id]; !ok {
		http.Error(res, "task is not found", http.StatusBadRequest)
		return
	}

	task, err := json.Marshal(tasks[id])
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
	
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(task)
}

//Обработчик удаления задачи по ID
func deleteTask(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")

	if _, ok := tasks[id]; !ok {
		http.Error(res, "task is not found", http.StatusBadRequest)
		return
	}

	delete(tasks, id)

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()
	// здесь регистрируйте ваши обработчики
	// ...
	r.Get("/tasks", getTasks)
	r.Post("/tasks", postTask)
	r.Get("/tasks/{id}", getTask)
	r.Delete("/tasks/{id}", deleteTask)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
