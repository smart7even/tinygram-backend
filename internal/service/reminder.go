package service

import (
	"fmt"
	"time"

	"github.com/smart7even/golang-do/internal/domain"
)

type ReminderRepo interface {
	Create(reminder domain.Reminder) error
	ReadAll(userId string) ([]domain.Reminder, error)
	Read(id int, userId string) (domain.Reminder, error)
	Update(reminder domain.Reminder) error
	Delete(id int, userId string) error
	GetClosestReminders() ([]domain.Reminder, error)
}

func NewReminderService(reminderRepo ReminderRepo) *ReminderService {
	return &ReminderService{
		reminderRepo: reminderRepo,
	}
}

type ReminderService struct {
	reminderRepo ReminderRepo
}

func (s *ReminderService) Create(reminder domain.Reminder) error {
	return s.reminderRepo.Create(reminder)
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

func StartReminderChecker(reminderService *ReminderService) {
	fmt.Println("Starting reminder checker")

	ticker := time.NewTicker(10 * time.Second) // Check every 10 seconds
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				checkReminders(reminderService)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func checkReminders(reminderService *ReminderService) {
	reminders, err := reminderService.GetClosestReminders()
	if err != nil {
		return
	}

	for _, reminder := range reminders {
		// Send reminder

		// Print reminder to console
		fmt.Printf("Reminder: %v\n", reminder)
	}
}
