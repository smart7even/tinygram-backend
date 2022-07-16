package service

type Services struct {
	Todo TodoService
	User UserService
}

type Deps struct {
	TodoRepo TodoRepo
}
