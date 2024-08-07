package dbmodel

import (
	api "asset-go/src/api"
	"log"
)

type User struct {
	Uid     int64
	Email   string
	Name    string
	Gender  string
	Picture string
	LineId  *string
}

func GetUserByEmail(email string) (*User, error) {

	db := GetDb()
	log.Println("Query user with mail", email)
	rows, err := db.Query("SELECT * FROM user WHERE Email=?", email)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.Uid, &user.Email, &user.Name, &user.Gender, &user.Picture, &user.LineId)
		if err != nil {
			log.Println("Failed to find user with mail", email, err)
			return nil, err
		} else {
			log.Println("Find user", user)
			return &user, nil
		}
	}
	log.Println("User not found.")
	return nil, nil
}

func InsertUser(user api.UserInfo) (*User, error) {

	var u *User

	u, err = GetUserByEmail(user.Email)
	if u != nil {
		log.Println("User exist", u)
		return u, nil
	}
	log.Println("Ready to create a new user", user)

	db := GetDb()
	stmt, err := db.Prepare("INSERT INTO user SET Email=?, Name=?, Gender=?, Picture=?")
	if err != nil {
		return nil, err
	}

	_, queryError := stmt.Exec(
		user.Email, user.Name, user.Gender, user.Picture)
	if queryError != nil {
		log.Println("Failed to insert", queryError)
		return nil, err
	}
	return GetUserByEmail(user.Email)
}
