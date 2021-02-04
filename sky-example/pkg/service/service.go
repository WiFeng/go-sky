package service

// Service ...
type Service struct {
	Example ExampleService
}

// New ...
func New() Service {
	return Service{
		Example: ExampleService{},
	}
}
