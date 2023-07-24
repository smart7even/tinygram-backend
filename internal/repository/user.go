package repository

import (
	"context"
	"database/sql"
	"fmt"

	firebase "firebase.google.com/go"
	"github.com/smart7even/golang-do/internal/domain"
	"github.com/smart7even/golang-do/internal/service"
)

type PGUserRepo struct {
	db          *sql.DB
	firebaseApp firebase.App
}

func NewPGUserRepo(db *sql.DB, firebase firebase.App) service.UserRepo {
	return &PGUserRepo{
		db:          db,
		firebaseApp: firebase,
	}
}

func (r *PGUserRepo) Create(token string) error {
	auth, err := r.firebaseApp.Auth(context.Background())

	if err != nil {
		return fmt.Errorf("error while creating user: %v", err)
	}

	authToken, err := auth.VerifyIDToken(context.Background(), token)

	if err != nil {
		return fmt.Errorf("error while verifying token: %v", err)
	}

	firebaseUser, getUserErr := auth.GetUser(context.Background(), authToken.UID)

	if getUserErr != nil {
		return fmt.Errorf("error while getting user info from Firebase: %v", getUserErr)
	}

	_, createUserErr := r.db.Exec("INSERT INTO users(id, name) VALUES ($1, $2)", firebaseUser.UID, firebaseUser.DisplayName)

	return createUserErr
}

func (r *PGUserRepo) ReadAll() ([]domain.User, error) {
	rows, err := r.db.Query("SELECT id, name FROM users")

	if err != nil {
		fmt.Printf("Error while requesting users: %v", err)
		return nil, err
	}

	defer rows.Close()

	var users = make([]domain.User, 0)

	for rows.Next() {
		var user domain.User
		rows.Scan(&user.Id, &user.Name)
		users = append(users, user)
	}

	return users, nil
}

func (r *PGUserRepo) Read(id string) (domain.User, error) {
	row := r.db.QueryRow("SELECT id, name FROM users WHERE id = $1", id)

	var user domain.User

	err := row.Scan(&user.Id, &user.Name)

	if err != nil {
		return user, fmt.Errorf("error while getting user: %v", err)
	}

	return user, nil
}

func (r *PGUserRepo) ReadByToken(token string) (*domain.User, error) {
	auth, err := r.firebaseApp.Auth(context.Background())

	if err != nil {
		return nil, fmt.Errorf("error while creating user: %v", err)
	}

	authToken, err := auth.VerifyIDToken(context.Background(), token)

	if err != nil {
		return nil, fmt.Errorf("error while verifying token: %v", err)
	}

	firebaseUser, getUserErr := auth.GetUser(context.Background(), authToken.UID)

	if getUserErr != nil {
		return nil, fmt.Errorf("error while getting user info from Firebase: %v", getUserErr)
	}

	return &domain.User{Id: firebaseUser.UID, Name: firebaseUser.DisplayName}, nil
}

func (r *PGUserRepo) Update(user domain.User) error {

	res, err := r.db.Exec("UPDATE users SET name = $1 WHERE id = $2", user.Name, user.Id)

	if err != nil {
		fmt.Printf("Error while editing user: %v", err)
		return err
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		fmt.Printf("Error while getting affected rows: %v", err)
		return err
	}

	if rowsAffected == 1 {
		return nil
	} else {
		return &service.UserDoesNotExist{UserId: user.Id}
	}
}

func (r *PGUserRepo) Delete(token string) error {
	auth, err := r.firebaseApp.Auth(context.Background())

	if err != nil {
		return fmt.Errorf("error while creating user: %v", err)
	}

	authToken, err := auth.VerifyIDToken(context.Background(), token)

	if err != nil {
		return fmt.Errorf("error while verifying token: %v", err)
	}

	res, err := r.db.Exec("DELETE FROM users WHERE id = $1", authToken.UID)

	if err != nil {
		return fmt.Errorf("error while deleting user: %v", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return fmt.Errorf("error while getting affected rows: %v", err)
	}

	if rowsAffected == 1 {
		return nil
	} else {
		return &service.UserDoesNotExist{UserId: authToken.UID}
	}
}
