package ports

import (
	"context"

	"github.com/yoannduc/go-project-template/internal/dtos"
)

type ExampleService interface {
	GetAll(context.Context) ([]dtos.Example, error)
	GetByLabelContaining(context.Context, string) ([]dtos.Example, error)
	GetByID(context.Context, int) (dtos.Example, error)
	Post(context.Context, dtos.Example) (dtos.Example, error)
	Patch(context.Context, int, dtos.Example) (dtos.Example, error)
	Delete(context.Context, int) (dtos.Example, error)
}
