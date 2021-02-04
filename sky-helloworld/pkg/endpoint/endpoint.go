package endpoint

import "github.com/WiFeng/go-sky/sky-helloworld/pkg/service"

// Endpoints ...
type Endpoints struct {
	Hello HelloEndpoints
}

// New ...
func New(s service.Service) Endpoints {

	return Endpoints{
		Hello: NewHelloEndpoints(s),
	}
}
