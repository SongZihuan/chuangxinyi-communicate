package utils

import "database/sql"

func GetUserName(phone string, email sql.NullString, username sql.NullString, nickname sql.NullString) string {
	if nickname.Valid && len(nickname.String) != 0 {
		return nickname.String
	}

	if username.Valid && len(username.String) != 0 {
		return username.String
	}

	if email.Valid && len(email.String) != 0 {
		return email.String
	}

	return phone
}
