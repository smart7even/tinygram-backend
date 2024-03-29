package repository

import (
	"database/sql"

	"github.com/smart7even/golang-do/internal/domain"
)

type PGReminderRepo struct {
	db *sql.DB
}

func NewPGReminderRepo(db *sql.DB) *PGReminderRepo {
	return &PGReminderRepo{
		db: db,
	}
}

func (r *PGReminderRepo) Create(reminder domain.Reminder) (int, error) {
	lastInsertedId := 0
	err := r.db.QueryRow("INSERT INTO reminders (user_id, name, description, remind_at, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id", reminder.UserId, reminder.Name, reminder.Description, reminder.RemindAt, reminder.CreatedAt).Scan(&lastInsertedId)
	return lastInsertedId, err
}

func (r *PGReminderRepo) ReadAll(userId string) ([]domain.Reminder, error) {
	rows, err := r.db.Query("SELECT id, name, description, remind_at, created_at, user_id FROM reminders WHERE user_id = $1", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reminders := []domain.Reminder{}
	for rows.Next() {
		var reminder domain.Reminder
		if err := rows.Scan(&reminder.Id, &reminder.Name, &reminder.Description, &reminder.RemindAt, &reminder.CreatedAt, &reminder.UserId); err != nil {
			return nil, err
		}
		reminders = append(reminders, reminder)
	}

	return reminders, nil
}

func (r *PGReminderRepo) Read(id int, userId string) (domain.Reminder, error) {
	var reminder domain.Reminder
	err := r.db.QueryRow("SELECT id, name, description, remind_at, created_at FROM reminders WHERE id = $1 AND user_id = $2", id, userId).Scan(&reminder.Id, &reminder.Name, &reminder.Description, &reminder.RemindAt, &reminder.CreatedAt)
	return reminder, err
}

func (r *PGReminderRepo) Update(reminder domain.Reminder) error {
	_, err := r.db.Exec("UPDATE reminders SET name = $1, description = $2, remind_at = $3 WHERE id = $4 AND user_id = $5", reminder.Name, reminder.Description, reminder.RemindAt, reminder.Id, reminder.UserId)
	return err
}

func (r *PGReminderRepo) Delete(id int, userId string) error {
	_, err := r.db.Exec("DELETE FROM reminders WHERE id = $1 AND user_id = $2", id, userId)
	return err
}

func (r *PGReminderRepo) GetClosestReminders() ([]domain.Reminder, error) {
	rows, err := r.db.Query("SELECT id, user_id, name, description, remind_at, created_at FROM reminders WHERE remind_at > NOW()")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reminders := []domain.Reminder{}
	for rows.Next() {
		var reminder domain.Reminder
		if err := rows.Scan(&reminder.Id, &reminder.UserId, &reminder.Name, &reminder.Description, &reminder.RemindAt, &reminder.CreatedAt); err != nil {
			return nil, err
		}
		reminders = append(reminders, reminder)
	}

	return reminders, nil
}

func (r *PGReminderRepo) CreateReminderSent(reminderSent domain.ReminderSent) error {
	// write all fields to the database: id, user_id, reminder_id, device_id, sent_at
	_, err := r.db.Exec("INSERT INTO reminder_sent (user_id, reminder_id, device_id, sent_at) VALUES ($1, $2, $3, $4)", reminderSent.UserId, reminderSent.ReminderId, reminderSent.DeviceId, reminderSent.SentAt)
	return err
}

func (r *PGReminderRepo) ReadReminderSent(reminderId int, userId string, deviceId int) (domain.ReminderSent, error) {
	var reminderSent domain.ReminderSent
	err := r.db.QueryRow("SELECT id, user_id, reminder_id, device_id, sent_at FROM reminder_sent WHERE reminder_id = $1 AND user_id = $2 AND device_id = $3", reminderId, userId, deviceId).Scan(&reminderSent.Id, &reminderSent.UserId, &reminderSent.ReminderId, &reminderSent.DeviceId, &reminderSent.SentAt)
	return reminderSent, err
}
