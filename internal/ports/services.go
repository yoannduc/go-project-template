package ports

import (
	"context"

	"github.com/yoannduc/go-project-template/internal/dtos"
)

// ExampleService is the api contract for the service that handles Examples.
type ExampleService interface {
	GetAll(context.Context) ([]dtos.Example, error)
	GetByLabelContaining(context.Context, string) ([]dtos.Example, error)
	GetByID(context.Context, int) (dtos.Example, error)
	Create(context.Context, dtos.Example) (dtos.Example, error)
	Update(context.Context, int, dtos.Example) (dtos.Example, error)
	Delete(context.Context, int) (dtos.Example, error)
}
