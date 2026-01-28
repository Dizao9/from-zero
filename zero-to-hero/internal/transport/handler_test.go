package transport

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"zero-to-hero/internal/storage"
)

type MockStorage struct{}

func (m *MockStorage) GetUsers() ([]storage.User, error) {
	return []storage.User{
		{Username: "test_user", Email: "text@example.com"},
	}, nil
}

func TestHandler_GetUsers(t *testing.T) {
	mockStore := &MockStorage{}
	handler := &Handler{
		Store: mockStore,
	}

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	handler.GetUsers(w, req)

	resp := w.Result()

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 status, god %d", resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type = application json")
	}

	var users []storage.User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		t.Error("decoder was falled")
	}

	if len(users) != 1 {
		t.Errorf("expected 1 user, got %d", len(users))
	}

	expectedUser := "test_user"
	if users[0].Username != expectedUser {
		t.Errorf("expected username %s, got %s", expectedUser, users[0].Username)
	}
}
