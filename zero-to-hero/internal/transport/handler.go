package transport

import (
	"encoding/json"
	"log"
	"net/http"
	"zero-to-hero/internal/storage"
)

type UserStorage interface {
	GetUsers() ([]storage.User, error)
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
