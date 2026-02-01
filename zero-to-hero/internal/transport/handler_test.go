package transport

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"zero-to-hero/internal/storage"
)

type MockStorage struct {
	MockGetUsersFunc func() ([]storage.User, error)
}

func (m *MockStorage) GetUsers() ([]storage.User, error) {
	if m.MockGetUsersFunc != nil {
		return m.MockGetUsersFunc()
	}

	return nil, nil
}

func (m *MockStorage) CreateUser(u storage.User) (int, error) {
	return 0, nil
}

func (m *MockStorage) DeleteUser(id int) error {
	return nil
}

func (m *MockStorage) UpdateUser(id int, u storage.User) error {
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
