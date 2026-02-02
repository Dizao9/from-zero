package transport

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"zero-to-hero/internal/storage"
)

type MockStorage struct {
	MockGetUsersFunc    func() ([]storage.User, error)
	MockCreateUsersFunc func(u storage.User) (int, error)
	MockUpdateUsersFunc func(id int, u storage.User) error
	MockDeleteUsersFunc func(id int) error
}

func (m *MockStorage) GetUsers() ([]storage.User, error) {
	if m.MockGetUsersFunc != nil {
		return m.MockGetUsersFunc()
	}

	return nil, nil
}

func (m *MockStorage) CreateUser(u storage.User) (int, error) {
	if m.MockCreateUsersFunc != nil {
		return m.MockCreateUsersFunc(u)
	}
	return 0, nil
}

func (m *MockStorage) DeleteUser(id int) error {
	if m.MockDeleteUsersFunc != nil {
		return m.MockDeleteUsersFunc(id)
	}
	return nil
}

func (m *MockStorage) UpdateUser(id int, u storage.User) error {
	if m.MockUpdateUsersFunc != nil {
		return m.MockUpdateUsersFunc(id, u)
	}
	return nil
}

func TestHandler_GetUsers_TableDriven(t *testing.T) {
	test := []struct {
		name           string
		mockReturnData []storage.User
		mockReturnErr  error
		expectedCode   int
	}{
		{
			name:           "Success Case",
			mockReturnData: []storage.User{{Username: "Max"}},
			mockReturnErr:  nil,
			expectedCode:   http.StatusOK,
		},
		{
			name:           "DB Error",
			mockReturnData: nil,
			mockReturnErr:  errors.New("datal db error"),
			expectedCode:   http.StatusInternalServerError,
		},
		{
			name:           "Empty case",
			mockReturnData: []storage.User{},
			mockReturnErr:  nil,
			expectedCode:   http.StatusOK,
		},
	}

	for _, tc := range test {
		t.Run(tc.name, func(t *testing.T) {
			mockStore := &MockStorage{
				MockGetUsersFunc: func() ([]storage.User, error) {
					return tc.mockReturnData, tc.mockReturnErr
				},
			}

			h := &Handler{Store: mockStore}

			req := httptest.NewRequest(http.MethodGet, "/users", nil)
			w := httptest.NewRecorder()

			h.GetUsers(w, req)

			if w.Code != tc.expectedCode {
				t.Errorf("want %d, got %d", tc.expectedCode, w.Code)
			}
		})
	}
}

func TestHandler_CreateUser_TableDriven(t *testing.T) {
	test := []struct {
		name          string
		mockReturnID  int
		mockReturnErr error
		inputUser     storage.User
		expectedCode  int
	}{
		{
			name:          "Success",
			mockReturnID:  1,
			mockReturnErr: nil,
			inputUser:     storage.User{Username: "test", Email: "test@gmail.com"},
			expectedCode:  http.StatusCreated,
		},
		{
			name:          "Storage Error",
			mockReturnID:  0,
			mockReturnErr: errors.New("db down"),
			inputUser:     storage.User{Username: "test", Email: "test@gmail.com"},
			expectedCode:  http.StatusInternalServerError,
		},
		{
			name:          "Fields Required (Validation error)",
			mockReturnID:  0,
			mockReturnErr: nil,
			inputUser:     storage.User{Username: "", Email: "s"},
			expectedCode:  http.StatusBadRequest,
		},
	}

	for _, tc := range test {
		t.Run(tc.name, func(t *testing.T) {
			mockStore := &MockStorage{
				MockCreateUsersFunc: func(u storage.User) (int, error) {
					return tc.mockReturnID, tc.mockReturnErr
				},
			}

			h := &Handler{Store: mockStore}

			bodyBytes, _ := json.Marshal(tc.inputUser)
			bodyReader := bytes.NewReader(bodyBytes)

			req := httptest.NewRequest(http.MethodPost, "/users", bodyReader)

			w := httptest.NewRecorder()

			h.CreateUser(w, req)

			if w.Code != tc.expectedCode {
				t.Errorf("want status %d, got %d", tc.expectedCode, w.Code)
			}
		})
	}
}

func TestHandler_UpdateUser_TableDriven(t *testing.T) {
	tests := []struct {
		name         string
		inputData    storage.User
		inputID      string
		returnErr    error
		expectedCode int
	}{
		{
			name:         "Success",
			inputData:    storage.User{Username: "testUpdate", Email: "test@gmail.com"},
			inputID:      "1",
			returnErr:    nil,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Field required (Validation)",
			inputData:    storage.User{Username: "", Email: ""},
			inputID:      "1",
			returnErr:    nil,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "DB failed",
			inputData:    storage.User{Username: "testUpdate", Email: "test@gmail.com"},
			inputID:      "1",
			returnErr:    errors.New("db failed"),
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "User Not Found",
			inputData:    storage.User{Username: "testUpdate", Email: "test@gmail.com"},
			inputID:      "1",
			returnErr:    storage.UserNotFound,
			expectedCode: http.StatusNoContent,
		},
		{
			name:         "Invalid id (validation)",
			inputData:    storage.User{Username: "testUpdate", Email: "test@gmail.com"},
			inputID:      "aaa",
			returnErr:    nil,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockStore := &MockStorage{
				MockUpdateUsersFunc: func(id int, u storage.User) error {
					return tc.returnErr
				},
			}

			h := &Handler{Store: mockStore}
			bodyBytes, _ := json.Marshal(tc.inputData)
			bodyReader := bytes.NewReader(bodyBytes)

			req := httptest.NewRequest(http.MethodPut, "/users/{id}", bodyReader)
			req.SetPathValue("id", tc.inputID)
			w := httptest.NewRecorder()
			h.UpdateUser(w, req)

			if w.Code != tc.expectedCode {
				t.Errorf("expected %d, got %d", tc.expectedCode, w.Code)
			}
		})
	}
}

func TestHandler_DeleteUser_TableDriven(t *testing.T) {
	tests := []struct {
		name         string
		inputID      string
		returnErr    error
		expectedCode int
	}{
		{
			name:         "Success",
			inputID:      "1",
			returnErr:    nil,
			expectedCode: http.StatusNoContent,
		},
		{
			name:         "Invalid user id",
			inputID:      "a",
			returnErr:    nil,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "User not found",
			inputID:      "1",
			returnErr:    storage.UserNotFound,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "DB failed",
			inputID:      "1",
			returnErr:    errors.New("db failed"),
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockStore := &MockStorage{
				MockDeleteUsersFunc: func(id int) error {
					return tc.returnErr
				},
			}

			h := &Handler{Store: mockStore}

			req := httptest.NewRequest(http.MethodDelete, "/users/{id}", nil)
			req.SetPathValue("id", tc.inputID)

			w := httptest.NewRecorder()

			h.DeleteUser(w, req)

			if w.Code != tc.expectedCode {
				t.Errorf("expected %d, got %d", tc.expectedCode, w.Code)
			}
		})
	}
}
