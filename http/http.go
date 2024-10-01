package http

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "httpServer/docs"
	"math/rand"
	"net/http"
	"time"
)

type Storage interface {
	GetTaskResult(taskID string) (*map[string]string, error)
	GetTaskStatus(taskID string) (*string, error)
	CreateTask(string) error
	DoTask(taskID string, result string) error
}

type Server struct {
	storage Storage
}

type TaskStatus struct {
	Status string `json:"status"`
}

type TaskResponse struct {
	TaskID string `json:"task_id"`
}

func newServer(storage Storage) *Server {
	return &Server{storage: storage}
}

func generateRandomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	result := make([]rune, n)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

// @Summary Create a task
// @Description Create a task and get a taskID
// @Tags tasks
// @Accept json
// @Produce json
// @Success 201 {object} TaskResponse "Task successfully created"
// @Router /task [post]
func (s *Server) createTaskHandler(w http.ResponseWriter, r *http.Request) {

	taskId := uuid.New().String()

	err := s.storage.CreateTask(taskId)
	if err != nil {
		http.Error(w, "Can't create a task", http.StatusInternalServerError)
		return
	}

	go func() {
		time.Sleep(20 * time.Second)

		result := generateRandomString(10)

		s.storage.DoTask(taskId, result)
	}()

	response := TaskResponse{TaskID: taskId}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(response)

}

// @Summary Get a status of task
// @Description Get a status of task by id
// @Tags tasks
// @Param taskID path string true "taskID"
// @Produce json
// @Success 200 {object} TaskStatus "Status of task"
// @Failure 404 {string} string "Task not found"
// @Router /status/{taskID} [get]
func (s *Server) statusTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "taskID")

	status, err := s.storage.GetTaskStatus(taskID)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	response := TaskStatus{Status: *status}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// @Summary Get a result of task
// @Description Get a result of task by id
// @Tags tasks
// @Param taskID path string true "taskID"
// @Produce json
// @Success 200 {object} map[string]string "result of task"
// @Failure 404 {string} string "Can't find a task"
// @Router /result/{taskID} [get]
func (s *Server) resultTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "taskID")

	result, err := s.storage.GetTaskResult(taskID)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(*result)
}

// @Summary run server
// @Description create and run http server
// @Tags server
// @Param addr query string true "Address for run server"
// @Router /server [get]
func CreateAndRunServer(storage Storage, addr string) error {
	server := newServer(storage)

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // Разрешить все источники (можно ограничить до определённых доменов)
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	}))

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Get("/status/{taskID}", server.statusTaskHandler)
	r.Get("/result/{taskID}", server.resultTaskHandler)
	r.Post("/task", server.createTaskHandler)

	httpServer := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	return httpServer.ListenAndServe()
}
