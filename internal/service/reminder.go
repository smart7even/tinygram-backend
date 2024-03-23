package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/smart7even/golang-do/internal/domain"
)

type ReminderRepo interface {
	Create(reminder domain.Reminder) (int, error)
	ReadAll(userId string) ([]domain.Reminder, error)
	Read(id int, userId string) (domain.Reminder, error)
	Update(reminder domain.Reminder) error
	Delete(id int, userId string) error
	GetClosestReminders() ([]domain.Reminder, error)
	CreateReminderSent(reminderSent domain.ReminderSent) error
	ReadReminderSent(reminderId int, userId string, deviceId int) (domain.ReminderSent, error)
}

func NewReminderService(reminderRepo ReminderRepo) *ReminderService {
	return &ReminderService{
		reminderRepo: reminderRepo,
	}
}

type ReminderService struct {
	reminderRepo ReminderRepo
}

func (s *ReminderService) Create(reminder domain.Reminder) (*domain.Reminder, error) {
	id, err := s.reminderRepo.Create(reminder)

	if err != nil {
		return nil, err
	}

	reminder.Id = id

	return &reminder, nil
}

func (s *ReminderService) ReadAll(userId string) ([]domain.Reminder, error) {
	return s.reminderRepo.ReadAll(userId)
}

func (s *ReminderService) Read(id int, userId string) (domain.Reminder, error) {
	return s.reminderRepo.Read(id, userId)
}

func (s *ReminderService) Update(reminder domain.Reminder) error {
	return s.reminderRepo.Update(reminder)
}

func (s *ReminderService) Delete(id int, userId string) error {
	return s.reminderRepo.Delete(id, userId)
}

func (s *ReminderService) GetClosestReminders() ([]domain.Reminder, error) {
	return s.reminderRepo.GetClosestReminders()
}

func (s *ReminderService) CreateReminderSent(reminderSent domain.ReminderSent) error {
	return s.reminderRepo.CreateReminderSent(reminderSent)
}

func (s *ReminderService) ReadReminderSent(reminderId int, userId string, deviceId int) (domain.ReminderSent, error) {
	return s.reminderRepo.ReadReminderSent(reminderId, userId, deviceId)
}

func StartReminderChecker(s *Services, firebaseApp *firebase.App) {
	fmt.Println("Starting reminder checker")

	client, err := firebaseApp.Messaging(context.Background())

	if err != nil {
		fmt.Printf("Error while getting Firebase Cloud Messaging client: %v", err)
		return
	}

	ticker := time.NewTicker(10 * time.Second) // Check every 10 seconds
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				checkReminders(s, client)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func checkReminders(s *Services, client *messaging.Client) {
	reminders, err := s.Reminder.GetClosestReminders()
	if err != nil {
		return
	}

	for _, reminder := range reminders {
		// Send reminder

		// Print reminder to console
		fmt.Printf("Reminder: %v\n", reminder)

		timeUntilReminder := time.Until(reminder.RemindAt)

		// Print time until reminder to console
		fmt.Printf("Time until reminder: %v\n", timeUntilReminder)

		// If reminder is 1 minute away, send push notification
		if timeUntilReminder < time.Minute {
			// Get devices for user
			devices, err := s.Device.ReadAll(reminder.UserId)

			if err != nil {
				fmt.Printf("Error while getting devices for user: %v", err)
				return
			}

			// Send push notification to each device
			for _, device := range devices {
				// Check if reminder was already sent to this device
				reminderSent, err := s.Reminder.ReadReminderSent(reminder.Id, reminder.UserId, device.Id)

				if err != nil && err != sql.ErrNoRows {
					fmt.Printf("Error while reading reminder sent: %v", err)
					continue
				}

				if reminderSent.Id != 0 {
					fmt.Printf("Reminder already sent to device: %v\n", device)
					continue
				}

				// Send push notification
				fmt.Printf("Sending push notification to device: %v\n", device)

				message := &messaging.Message{
					Data: map[string]string{
						"reminder_id": fmt.Sprintf("%v", reminder.Id),
					},
					Notification: &messaging.Notification{
						Title: reminder.Name,
						Body:  reminder.Description,
					},
					Token: device.DeviceToken,
				}

				_, err = client.Send(context.Background(), message)

				if err != nil {
					fmt.Printf("Error while sending push notification: %v", err)
					// TODO: if token is invalid, set device as inactive
					continue
				}

				// Create reminder sent
				reminderSent = domain.ReminderSent{
					ReminderId: reminder.Id,
					UserId:     reminder.UserId,
					DeviceId:   device.Id,
					SentAt:     time.Now(),
				}

				err = s.Reminder.CreateReminderSent(reminderSent)

				if err != nil {
					fmt.Printf("Error while creating reminder sent: %v", err)
					continue
				}

				fmt.Printf("Reminder sent: %v\n", reminderSent)
			}
		}
	}
}
