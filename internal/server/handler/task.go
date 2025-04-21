package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"test-workmate/internal/service"

	"github.com/gorilla/mux"
)

type TaskHandler struct {
	svc *service.TaskService
}

func NewTaskHandler(svc *service.TaskService) *TaskHandler {
	return &TaskHandler{svc: svc}
}

// @Summary      Создать задачу
// @Description  Создает новую задачу и возвращает её ID
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]string
// @Failure      500  {string}  string
// @Router       /tasks [post]
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	id, err := h.svc.CreateTask(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]string{"task_id": id}
	json.NewEncoder(w).Encode(resp)
}

// @Summary      Получить задачу
// @Description  Возвращает задачу по ID
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID задачи"
// @Success      200  {object}  domain.Task
// @Failure      404  {string}  string
// @Router       /tasks/{id} [get]
func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	task, err := h.svc.GetTask(r.Context(), id)
	if err != nil {
		http.Error(w, "failed to get task: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if task == nil {
		http.Error(w, fmt.Sprintf("task with id '%s' not found", id), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(task)
}
