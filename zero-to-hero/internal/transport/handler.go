package transport

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"zero-to-hero/internal/storage"
)

type UserStorage interface {
	GetUsers() ([]storage.User, error)
	CreateUser(user storage.User) (int, error)
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
