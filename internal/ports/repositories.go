package ports

import (
	"context"

	"github.com/yoannduc/go-project-template/internal/domain"
)

// ExampleRepository is the api contract for the repository that
// handles Examples.
type ExampleRepository interface {
	FindAll(context.Context) ([]domain.Example, error)
	FindByLabelContaining(context.Context, string) ([]domain.Example, error)
	FindByID(context.Context, int) (domain.Example, error)
	Create(context.Context, domain.Example) (domain.Example, error)
	Update(context.Context, int, domain.Example) (domain.Example, error)
	Delete(context.Context, int) (domain.Example, error)
}
