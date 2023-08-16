package database

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"api-boilerplate/models/domain"
)

// User represents the user database table.
type User struct {
	ID        int       `db:"id"`
	Username  string    `db:"username"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
}

type UserService struct {
	db *DB
}

func NewUserService(db *DB) *UserService {
	return &UserService{
		db: db,
	}
}

// CreateUser creates user in database table users with username and password
func (s *UserService) CreateUser(ctx context.Context, username, hashedPassword string) (int, error) {
	userID, err := createUser(ctx, s.db.db, username, hashedPassword)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return userID, nil
}

func createUser(ctx context.Context, tx Queryable, username, hashedPassword string) (int, error) {
	var userID int
	rows, err := tx.NamedQuery(`INSERT INTO Users (username, password, created_at) VALUES (:username, :password, :time) RETURNING id`,
		map[string]interface{}{
			"username": username,
			"password": hashedPassword,
			"time":     time.Now(),
		})
	if err != nil {
		log.Println(err)
		return 0, err
	}

	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&userID)
		if err != nil {
			log.Println(err)
			return 0, err
		}
	}

	if userID == 0 {
		return 0, errors.New("user not created")
	}

	return userID, nil
}

func transformUser(databaseUser User) domain.User {
	return domain.User{
		ID:        databaseUser.ID,
		Username:  databaseUser.Username,
		Password:  databaseUser.Password,
		CreatedAt: databaseUser.CreatedAt,
	}
}

// GetUser returns a User object for a given userID.
func (s *UserService) GetUser(ctx context.Context, userID int) (domain.User, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.User{}, err
	}
	defer tx.Rollback()

	user, err := getUser(ctx, tx, userID)
	if err != nil {
		log.Println(err)
		return domain.User{}, err
	}

	// ... more things here

	err = tx.Commit()
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

// getUser is the real function that gets the user from the database.
func getUser(ctx context.Context, tx Queryable, userID int) (domain.User, error) {
	user := User{}
	err := tx.GetContext(ctx, &user, "SELECT * FROM Users WHERE id = $1", userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, errors.New("no user exists with that id")
		}
	}

	return transformUser(user), nil
}

// GetUserFromUsername returns a User for a given username as string.
func (s *UserService) GetUserFromUsername(ctx context.Context, username string) (domain.User, bool) {
	user, err := getUserFromUsername(ctx, s.db.db, username)
	if err != nil {
		log.Println(err)
		return domain.User{}, false
	}

	return user, true
}

func getUserFromUsername(ctx context.Context, tx Queryable, username string) (domain.User, error) {
	user := User{}
	err := tx.GetContext(ctx, &user, "SELECT * FROM Users WHERE username = $1", username)
	if err != nil {
		log.Println(err)
		if err == sql.ErrNoRows {
			return domain.User{}, err
		}
	}

	return transformUser(user), nil
}

// DeleteUser deletes an entry from the table. Needs to have a valid userID.
func (s *UserService) DeleteUser(ctx context.Context, userID int) error {
	return deleteUser(ctx, s.db.db, userID)
}

func deleteUser(ctx context.Context, tx Queryable, userID int) error {
	_, err := tx.NamedExecContext(ctx, "DELETE Users where id=:id",
		map[string]interface{}{
			"id": userID,
		})
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
