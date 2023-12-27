package utils

import "database/sql"

func GetUserName(defaultName string, username sql.NullString, nickname sql.NullString) string {
	if nickname.Valid && len(nickname.String) != 0 {
		return nickname.String
	}

	if username.Valid && len(username.String) != 0 {
		return username.String
	}

	if len(defaultName) > 6 {
		return defaultName[0:6]
	}

	return defaultName
}

func GetUserNameByPhone(phone string, username sql.NullString, nickname sql.NullString) string {
	if nickname.Valid && len(nickname.String) != 0 {
		return nickname.String
	}

	if username.Valid && len(username.String) != 0 {
		return username.String
	}

	return phone
}
