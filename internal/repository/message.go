package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/smart7even/golang-do/internal/domain"
	"github.com/smart7even/golang-do/internal/service"
)

type MySqlMessageRepo struct {
	db *sql.DB
}

func NewMySQLMessageRepo(db *sql.DB) service.MessageRepo {
	return &MySqlMessageRepo{
		db: db,
	}
}

func (r *MySqlMessageRepo) Create(message domain.Message, userId string) error {
	row := r.db.QueryRow("SELECT user_id, chat_id FROM chat_user WHERE user_id = ? AND chat_id = ?", userId, message.ChatId)
	if row.Err() == nil {
		row.Scan()
		message.Id = uuid.New().String()
		message.SentAt = time.Now()
		_, err := r.db.Exec("INSERT INTO messages (id, chat_id, user_id, text, sent_at) VALUES (?, ?, ?, ?, ?)", message.Id, message.ChatId, message.UserId, message.Text, message.SentAt)
		return err
	} else {
		return errors.New("you are not a member of this chat")
	}
}

func (r *MySqlMessageRepo) ReadAll(chatId string, userId string) ([]domain.Message, error) {
	row := r.db.QueryRow("SELECT user_id, chat_id FROM chat_user WHERE user_id = ? AND chat_id = ?", userId, chatId)

	if row.Err() == nil {
		rows, err := r.db.Query("SELECT id, chat_id, user_id, text, sent_at FROM messages WHERE chat_id = ? ORDER BY sent_at", chatId)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		messages := []domain.Message{}
		for rows.Next() {
			var message domain.Message
			if err := rows.Scan(&message.Id, &message.ChatId, &message.UserId, &message.Text, &message.SentAt); err != nil {
				return nil, err
			}
			messages = append(messages, message)
		}

		return messages, nil
	} else {
		return []domain.Message{}, errors.New("you are not a member of this chat")
	}
}

func (r *MySqlMessageRepo) Update(message domain.Message, userId string) error {
	row := r.db.QueryRow("SELECT user_id, chat_id FROM chat_user WHERE user_id = ? AND chat_id = ?", userId, message.ChatId)

	if row.Err() == nil {
		_, err := r.db.Exec("UPDATE messages SET text = ? WHERE id = ?", message.Text, message.Id)
		return err
	} else {
		return errors.New("you are not a member of this chat")
	}
}

func (r *MySqlMessageRepo) Delete(id string, userId string) error {
	row := r.db.QueryRow("SELECT user_id, chat_id FROM chat_user WHERE user_id = ? AND chat_id = ?", userId, id)

	if row.Err() == nil {
		_, err := r.db.Exec("DELETE FROM messages WHERE id = ?", id)
		return err
	} else {
		return errors.New("you are not a member of this chat")
	}
}
