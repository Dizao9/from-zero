package transport

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"zero-to-hero/internal/storage"
)

type UserStorage interface {
	GetUsers() ([]storage.User, error)
	CreateUser(user storage.User) (int, error)
	DeleteUser(id int) error
	UpdateUser(id int, u storage.User) error
}
type Handler struct {
	Store UserStorage
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	resp, err := h.Store.GetUsers()
	if err != nil {
		log.Printf("failed to get users data")
		http.Error(w, "Error with getting data abot users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("failed to encode response data")
		http.Error(w, "Failed to encode data", http.StatusInternalServerError)
		return
	}
}

// @Summary Создать пользователя
// @Description Создает нового юзера и возвращает его id
// @Accept json
// @Produce json
// @Param input body storage.User true "Данные пользователя"
// @Success 201 {object} map[string]int
// @Failure 400 {string} string "Bad Request"
// @Router /users [post]
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u storage.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Problem with decoder for request body", http.StatusBadRequest)
		return
	}

	if u.Username == "" {
		http.Error(w, "Field required", http.StatusBadRequest)
		return
	}

	if u.Email == "" {
		http.Error(w, "Field required", http.StatusBadRequest)
		return
	}

	id, err := h.Store.CreateUser(u)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"id":%d}`, id) //можно создать структуру и записать через неё, но не слишком ли для одного поля?
}

func (h *Handler) HandleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetUsers(w, r)
	case http.MethodPost:
		h.CreateUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// @Summary Удалить пользователя
// @Description Удаляет пользователя
// @Param path
// @Success 204
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal server error"
// @Failure 404 {string} string "User wasn't found to delete"
// @Router /users/{id} [delete]
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid format user id", http.StatusBadRequest)
		return
	}

	err = h.Store.DeleteUser(id)
	if err != nil {
		if err == storage.UserNotFound {
			http.Error(w, "user wasn't found to delete", http.StatusNotFound)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// @Summary Обновить поля пользователя
// @Description Обновляет существующие поля пользователя, необходимо полная структура юзера
// @Accept json
// @Produce json
// @Param path
// @Success 200
// @Failure 400 [string] string "Bad Request"
// @Failure 404 [string] string "User not found"
// @Failure 500 [string] string "Internal server error"
// @Router /users/{id} [put]
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user id format", http.StatusBadRequest)
		return
	}
	var u storage.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Problem with decoder", http.StatusBadRequest)
		return
	}

	if u.Username == "" {
		http.Error(w, "field required", http.StatusBadRequest)
		return
	}

	if u.Email == "" {
		http.Error(w, "field required", http.StatusBadRequest)
		return
	}

	err = h.Store.UpdateUser(id, u)
	if err != nil {
		if err == storage.UserNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		log.Printf("[UPDATE USER] internal server error: %v\n", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
