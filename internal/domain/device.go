package domain

type Device struct {
	Id          int    `json:"id"`
	UserId      string `json:"userId"`
	DeviceId    string `json:"deviceId"`
	DeviceOs    string `json:"deviceOs"`
	DeviceToken string `json:"deviceToken"`
}
