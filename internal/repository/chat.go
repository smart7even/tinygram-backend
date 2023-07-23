package repository

import (
	"database/sql"

	"github.com/smart7even/golang-do/internal/domain"
	"github.com/smart7even/golang-do/internal/service"
)

type PGChatRepo struct {
	db *sql.DB
}

func NewPGChatRepo(db *sql.DB) service.ChatRepo {
	return &PGChatRepo{
		db: db,
	}
}

func (r *PGChatRepo) Create(chat domain.Chat) error {
	_, err := r.db.Exec("INSERT INTO chats (id, name) VALUES ($1, $2)", chat.Id, chat.Name)
	return err
}

func (r *PGChatRepo) ReadAll() ([]domain.Chat, error) {
	rows, err := r.db.Query("SELECT id, name FROM chats")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	chats := []domain.Chat{}
	for rows.Next() {
		var chat domain.Chat
		if err := rows.Scan(&chat.Id, &chat.Name); err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}

	return chats, nil
}

func (r *PGChatRepo) Update(chat domain.Chat) error {
	_, err := r.db.Exec("UPDATE chats SET name = $1 WHERE id = $2", chat.Name, chat.Id)
	return err
}

func (r *PGChatRepo) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM chats WHERE id = $1", id)
	return err
}

func (r *PGChatRepo) Join(chatId string, userId string) error {
	_, err := r.db.Exec("INSERT INTO chat_user (chat_id, user_id) VALUES ($1, $2)", chatId, userId)
	return err
}
