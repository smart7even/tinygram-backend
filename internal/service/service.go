package service

type Services struct {
	Todo TodoService
}

type Deps struct {
	TodoRepo TodoRepo
}
