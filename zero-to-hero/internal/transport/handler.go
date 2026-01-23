package transport

import (
	"encoding/json"
	"log"
	"net/http"

	"zero-to-hero/internal/storage"
)

type Handler struct {
	Store *storage.Storage
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	res, err := h.Store.GetUsers()
	if err != nil {
		log.Printf("[Get Users]error with getting users data: %v", err)
		http.Error(w, "Failed to get users data", http.StatusInternalServerError) //По хорошему, валидация ошибок на отсутсвие юзеров и действительно ошибки
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("Encode error: %v", err) //как мы обсуждали, поздно слать ошибку, но что тогда делать?
	}
}
