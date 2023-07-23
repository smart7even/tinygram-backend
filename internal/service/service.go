package service

type Services struct {
	Todo    TodoService
	User    UserService
	Chat    ChatService
	Message MessageService
	Auth    AuthService
	Event   EventService
}

type Deps struct {
	TodoRepo    TodoRepo
	UserRepo    UserRepo
	ChatRepo    ChatRepo
	MessageRepo MessageRepo
	EventRepo   EventRepo
}
