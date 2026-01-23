package storage

import "database/sql"

type User struct {
	Username string `json:"username"`
	Email    string `json:"email,omitempty"`
}

type Storage struct {
	DB *sql.DB
}

func (s *Storage) GetUsers() ([]User, error) {
	rows, err := s.DB.Query(`SELECT username, COALESCE(email, '') AS email
	FROM users`)
	if err != nil {
		return []User{}, err
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var u User
		err = rows.Scan(&u.Username, &u.Email)
		if err != nil {
			return []User{}, err
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return []User{}, err
	}

	return users, nil
}
