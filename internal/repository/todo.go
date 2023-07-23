package repository

import (
	"database/sql"
	"fmt"

	"github.com/smart7even/golang-do/internal/domain"
	"github.com/smart7even/golang-do/internal/service"
)

type PGTodoRepo struct {
	db *sql.DB
}

func NewPGTodoRepo(db *sql.DB) service.TodoRepo {
	return &PGTodoRepo{
		db: db,
	}
}

func (r *PGTodoRepo) Create(todo domain.Todo) error {
	_, err := r.db.Exec("INSERT INTO todos(name, complete) VALUES ($1, $2)", todo.Name, todo.Complete)

	return err
}

func (r *PGTodoRepo) ReadAll() ([]domain.Todo, error) {
	rows, err := r.db.Query("SELECT id, name, complete FROM todos")

	if err != nil {
		fmt.Printf("Error while requesting todos: %v", err)
		return nil, err
	}

	defer rows.Close()

	var todos []domain.Todo = make([]domain.Todo, 0)

	for rows.Next() {
		var todo domain.Todo
		rows.Scan(&todo.Id, &todo.Name, &todo.Complete)
		todos = append(todos, todo)
	}

	return todos, nil
}

func (r *PGTodoRepo) Update(todo domain.Todo) error {

	res, err := r.db.Exec("UPDATE todos SET name = $1, complete = $2 WHERE id = $3", todo.Name, todo.Complete, todo.Id)

	if err != nil {
		fmt.Printf("Error while editing todo: %v", err)
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
		return &service.TodoDoesNotExist{TodoId: todo.Id}
	}
}

func (r *PGTodoRepo) Delete(todoId int64) error {
	res, err := r.db.Exec("DELETE FROM todos WHERE id = $1", todoId)

	if err != nil {
		return fmt.Errorf("error while deleting todo: %v", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return fmt.Errorf("error while getting affected rows: %v", err)
	}

	if rowsAffected == 1 {
		return nil
	} else {
		return &service.TodoDoesNotExist{TodoId: todoId}
	}
}
