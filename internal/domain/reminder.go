package domain

type Reminder struct {
	Id          int    `json:"id"`
	UserId      string `json:"userId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	RemindAt    string `json:"remindAt"`
	CreatedAt   string `json:"createdAt"`
}
