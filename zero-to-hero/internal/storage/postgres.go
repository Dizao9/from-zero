package storage

import "database/sql"

type Storage struct {
	DB *sql.DB
}

type User struct {
	Username string `json:"username"`
	Email    string `json:"email,omitempty"`
}

func (s *Storage) GetUsers() ([]User, error) {
	rows, err := s.DB.Query("SELECT username, COALESCE(email, '') FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.Username, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
