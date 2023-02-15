package database

import (
	"database/sql"
	"errors"
	"log"
	"time"
)

// User represents the user database table.
type User struct {
	ID        int       `db:"id"`
	Username  string    `db:"username"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
}

// CreateUser creates user in database table users with username and password
func CreateUser(username, hashedPassword string) (User, error) {
	var userID int
	rows, err := db.NamedQuery(`INSERT INTO Users (username, password, created_at) VALUES (:username, :password, :time) RETURNING id`,
		map[string]interface{}{
			"username": username,
			"password": hashedPassword,
			"time":     time.Now(),
		})
	if err != nil {
		log.Println(err)
		return User{}, err
	}

	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&userID)
		if err != nil {
			log.Println(err)
			return User{}, err
		}
	}

	if userID == 0 {
		return User{}, errors.New("user not created")
	}

	user, err := GetUser(userID)
	if err != nil {
		log.Println(err)
		return User{}, err
	}

	return user, nil
}

// GetUser returns a User object for a given userID.
func GetUser(userID int) (User, error) {
	user := User{}
	err := db.Get(&user, "SELECT * FROM Users WHERE id = $1", userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, errors.New("no user exists with that id")
		}
	}

	return user, nil
}

// GetUserFromUsername returns a User for a given username as string.
func GetUserFromUsername(username string) (User, bool) {
	user := User{}
	err := db.Get(&user, "SELECT * FROM Users WHERE username = $1", username)
	if err != nil {
		log.Println(err)
		if err == sql.ErrNoRows {
			return user, false
		}
	}

	return user, true
}

// DeleteUser deletes an entry from the table. Needs to have a valid userID.
func DeleteUser(userID int) error {
	_, err := db.NamedExec("DELETE Users where id=:id",
		map[string]interface{}{
			"id": userID,
		})
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
