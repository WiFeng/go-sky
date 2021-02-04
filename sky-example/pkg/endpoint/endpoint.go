package endpoint

import "github.com/WiFeng/go-sky/sky-example/pkg/service"

// Endpoints ...
type Endpoints struct {
	Example ExampleEndpoints
}

// New ...
func New(s service.Service) Endpoints {

	return Endpoints{
		Example: NewExampleEndpoints(s),
	}
}
