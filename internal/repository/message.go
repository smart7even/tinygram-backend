package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/smart7even/golang-do/internal/domain"
	"github.com/smart7even/golang-do/internal/service"
)

type PGMessageRepo struct {
	db *sql.DB
}

func NewPGMessageRepo(db *sql.DB) service.MessageRepo {
	return &PGMessageRepo{
		db: db,
	}
}

func (r *PGMessageRepo) Create(message domain.Message, userId string) error {
	row := r.db.QueryRow("SELECT user_id, chat_id FROM chat_user WHERE user_id = $1 AND chat_id = $2", userId, message.ChatId)
	if row.Err() == nil {
		row.Scan()
		message.Id = uuid.New().String()
		message.SentAt = time.Now()
		_, err := r.db.Exec("INSERT INTO messages (id, chat_id, user_id, text, sent_at) VALUES ($1, $2, $3, $4, $5)", message.Id, message.ChatId, message.UserId, message.Text, message.SentAt)
		return err
	} else {
		return errors.New("you are not a member of this chat")
	}
}

func (r *PGMessageRepo) ReadAll(chatId string, userId string) ([]domain.Message, error) {
	row := r.db.QueryRow("SELECT user_id, chat_id FROM chat_user WHERE user_id = $1 AND chat_id = $2", userId, chatId)

	if row.Err() == nil {
		rows, err := r.db.Query("SELECT messages.id, chat_id, user_id, users.name, text, sent_at FROM messages INNER JOIN users ON users.id = user_id WHERE chat_id = $1 ORDER BY sent_at;", chatId)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		messages := []domain.Message{}
		for rows.Next() {
			var message domain.Message
			if err := rows.Scan(&message.Id, &message.ChatId, &message.UserId, &message.UserName, &message.Text, &message.SentAt); err != nil {
				return nil, err
			}
			messages = append(messages, message)
		}

		return messages, nil
	} else {
		return []domain.Message{}, errors.New("you are not a member of this chat")
	}
}

func (r *PGMessageRepo) Update(message domain.Message, userId string) error {
	row := r.db.QueryRow("SELECT user_id, chat_id FROM chat_user WHERE user_id = $1 AND chat_id = $2", userId, message.ChatId)

	if row.Err() == nil {
		_, err := r.db.Exec("UPDATE messages SET text = $1 WHERE id = $2", message.Text, message.Id)
		return err
	} else {
		return errors.New("you are not a member of this chat")
	}
}

func (r *PGMessageRepo) Delete(id string, userId string) error {
	row := r.db.QueryRow("SELECT user_id, chat_id FROM chat_user WHERE user_id = $1 AND chat_id = $2", userId, id)

	if row.Err() == nil {
		_, err := r.db.Exec("DELETE FROM messages WHERE id = $1", id)
		return err
	} else {
		return errors.New("you are not a member of this chat")
	}
}
