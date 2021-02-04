package service

// Service ...
type Service struct {
	Hello HelloService
}

// New ...
func New() Service {
	return Service{
		Hello: HelloService{},
	}
}
