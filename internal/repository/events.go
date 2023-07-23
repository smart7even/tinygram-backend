package repository

import (
	"database/sql"
	"encoding/json"

	"github.com/smart7even/golang-do/internal/domain"
)

type MySqlEventRepo struct {
	db *sql.DB
}

func NewMySqlEventRepo(db *sql.DB) *MySqlEventRepo {
	return &MySqlEventRepo{
		db: db,
	}
}

func (r *MySqlEventRepo) Create(event domain.Event) error {
	payload, jsonErr := json.Marshal(event.Payload)

	if jsonErr != nil {
		return jsonErr
	}

	_, err := r.db.Exec("INSERT INTO events (id, name, description, created_at, payload) VALUES ($1, $2, $3, $4, $5)", event.Id, event.Name, event.Description, event.CreatedAt, payload)
	return err
}

func (r *MySqlEventRepo) ReadAll() ([]domain.Event, error) {
	rows, err := r.db.Query("SELECT id, name, description, created_at, payload FROM events")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []domain.Event{}
	for rows.Next() {
		var event domain.Event
		if err := rows.Scan(&event.Id, &event.Name, &event.Description, &event.CreatedAt, &event.Payload); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}
