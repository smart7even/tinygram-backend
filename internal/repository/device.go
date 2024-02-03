package repository

import (
	"database/sql"

	"github.com/smart7even/golang-do/internal/domain"
)

type PGDeviceRepo struct {
	db *sql.DB
}

func NewPGDeviceRepo(db *sql.DB) *PGDeviceRepo {
	return &PGDeviceRepo{
		db: db,
	}
}

func (r *PGDeviceRepo) Create(device *domain.Device) error {
	_, err := r.db.Exec("INSERT INTO devices (device_id, device_token, device_os, user_id) VALUES ($1, $2, $3, $4)", device.DeviceId, device.DeviceToken, device.DeviceOs, device.UserId)
	return err
}

func (r *PGDeviceRepo) ReadAll(userId string) ([]domain.Device, error) {
	rows, err := r.db.Query("SELECT id, device_id, device_token, device_os, user_id FROM devices WHERE user_id = $1", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	devices := []domain.Device{}
	for rows.Next() {
		var device domain.Device
		if err := rows.Scan(&device.Id, &device.DeviceId, &device.DeviceToken, &device.DeviceOs, &device.UserId); err != nil {
			return nil, err
		}
		devices = append(devices, device)
	}

	return devices, nil
}

func (r *PGDeviceRepo) Read(id int, deviceId string) (domain.Device, error) {
	var device domain.Device
	err := r.db.QueryRow("SELECT id, device_id, device_token, device_os, user_id FROM devices WHERE id = $1 AND device_id = $2", id, deviceId).Scan(&device.Id, &device.DeviceId, &device.DeviceToken, &device.DeviceOs, &device.UserId)
	return device, err
}

func (r *PGDeviceRepo) Update(device *domain.Device) error {
	// Update user_id, device_token, device_os by device_id
	_, err := r.db.Exec("UPDATE devices SET user_id = $1, device_token = $2, device_os = $3 WHERE device_id = $4", device.UserId, device.DeviceToken, device.DeviceOs, device.DeviceId)
	return err
}

func (r *PGDeviceRepo) Delete(id int, userId string) error {
	_, err := r.db.Exec("DELETE FROM devices WHERE id = $1 AND user_id = $2", id, userId)
	return err
}
