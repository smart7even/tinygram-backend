package service

import "github.com/smart7even/golang-do/internal/domain"

type DeviceRepo interface {
	Create(device *domain.Device) error
	ReadAll(userId string) ([]domain.Device, error)
	Read(id int, deviceId string) (domain.Device, error)
	ReadByDeviceId(deviceId string) (domain.Device, error)
	Update(device *domain.Device) error
	Delete(id int, userId string) error
}

func NewDeviceService(deviceRepo DeviceRepo) *DeviceService {
	return &DeviceService{
		deviceRepo: deviceRepo,
	}
}

type DeviceService struct {
	deviceRepo DeviceRepo
}

func (s *DeviceService) Create(device *domain.Device) error {
	// First check if device with device_id already exists
	savedDevice, err := s.deviceRepo.ReadByDeviceId(device.DeviceId)

	// If device with device_id already exists, update it
	if err == nil {
		savedDevice.DeviceToken = device.DeviceToken
		savedDevice.DeviceOs = device.DeviceOs
		savedDevice.UserId = device.UserId
		return s.deviceRepo.Update(&savedDevice)
	}

	return s.deviceRepo.Create(device)
}

func (s *DeviceService) ReadAll(userId string) ([]domain.Device, error) {
	return s.deviceRepo.ReadAll(userId)
}

func (s *DeviceService) Read(id int, deviceId string) (domain.Device, error) {
	return s.deviceRepo.Read(id, deviceId)
}

func (s *DeviceService) Update(device *domain.Device) error {
	return s.deviceRepo.Update(device)
}

func (s *DeviceService) Delete(id int, userId string) error {
	return s.deviceRepo.Delete(id, userId)
}
