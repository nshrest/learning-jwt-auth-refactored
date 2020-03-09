package userRepository

import (
	"database/sql"
	"learning-jwt-auth-refactored/models"
	"log"
)

type UserRepository struct{}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (u UserRepository) Signup(db *sql.DB, user models.User) models.User {
	stmt := "insert into users (username, password) values($1, $2) RETURNING id;"
	err := db.QueryRow(stmt, user.Username, user.Password).Scan(&user.ID)
	logFatal(err)
	user.Password = ""
	return user
}

func (u UserRepository) Login(db *sql.DB, user models.User) (models.User, error) {
	row := db.QueryRow("select * from users where username=$1", user.Username)
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return user, err
	}

	return user, nil
}
