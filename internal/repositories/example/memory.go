package example

import (
	"context"
	"fmt"
	"strings"

	"github.com/yoannduc/go-project-template/internal/domain"
	"github.com/yoannduc/go-project-template/internal/ports"
	"github.com/yoannduc/go-project-template/pkg/memorydb"
)

type memoryRepository struct {
	db memorydb.MemoryDB[domain.Example]
}

func NewMemoryRepository(db memorydb.MemoryDB[domain.Example]) ports.ExampleRepository {
	return memoryRepository{
		db: db,
	}
}

func (repo memoryRepository) FindAll(ctx context.Context) ([]domain.Example, error) {
	v, err := repo.db.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("memory repository: %w", err)
	}

	return v, nil
}

func (repo memoryRepository) FindByLabelContaining(ctx context.Context, search string) ([]domain.Example, error) {
	v, err := repo.db.FindFilterFunc(ctx, func(dom domain.Example) bool {
		return strings.Contains(dom.Label, search)
	})
	if err != nil {
		return nil, fmt.Errorf("memory repository: %w", err)
	}

	return v, nil
}

func (repo memoryRepository) FindByID(ctx context.Context, id int) (domain.Example, error) {
	v, err := repo.db.FindByID(ctx, id)
	if err != nil {
		return domain.Example{}, fmt.Errorf("memory repository: %w", err)
	}

	return v, nil
}

func (repo memoryRepository) Create(ctx context.Context, dom domain.Example) (domain.Example, error) {
	v, err := repo.db.Create(ctx, dom, func(a, b domain.Example) int {
		return strings.Compare(a.Label, b.Label)
	})
	if err != nil {
		return domain.Example{}, fmt.Errorf("memory repository: %w", err)
	}

	return v, nil
}

func (repo memoryRepository) Update(ctx context.Context, id int, dom domain.Example) (domain.Example, error) {
	v, err := repo.db.Update(ctx, id, dom, func(a, b domain.Example) int {
		return strings.Compare(a.Label, b.Label)
	})
	if err != nil {
		return domain.Example{}, fmt.Errorf("memory repository: %w", err)
	}

	return v, nil
}

func (repo memoryRepository) Delete(ctx context.Context, id int) (domain.Example, error) {
	v, err := repo.db.Delete(ctx, id)
	if err != nil {
		return domain.Example{}, fmt.Errorf("memory repository: %w", err)
	}

	return v, nil
}
