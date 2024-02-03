package domain

type Device struct {
	Id          string `json:"id"`
	UserId      string `json:"userId"`
	DeviceId    string `json:"deviceId"`
	DeviceOs    string `json:"deviceOs"`
	DeviceToken string `json:"deviceToken"`
}
