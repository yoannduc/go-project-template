package example

import (
	"context"
	"fmt"
	"strings"

	"github.com/yoannduc/go-project-template/internal/domain"
	"github.com/yoannduc/go-project-template/internal/ports"
	"github.com/yoannduc/go-project-template/pkg/memorydb"
)

func cmp(a, b domain.Example) bool {
	if a.ID == b.ID {
		return false
	}

	return strings.EqualFold(a.Label, b.Label)
}

type memoryRepository struct {
	db memorydb.MemoryDB[domain.Example]
}

func NewMemoryRepository(db memorydb.MemoryDB[domain.Example]) ports.ExampleRepository {
	return memoryRepository{
		db: db,
	}
}

func (repo memoryRepository) FindAll(ctx context.Context) ([]domain.Example, error) {
	return repo.db.FindAll(ctx)
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
	v, err := repo.db.FindByID(ctx, uint64(id))
	if err != nil {
		return domain.Example{}, fmt.Errorf("memory repository: %w", err)
	}

	return v, nil
}

func (repo memoryRepository) Create(ctx context.Context, dom domain.Example) (domain.Example, error) {
	v, err := repo.db.Create(ctx, dom, cmp)
	if err != nil {
		return domain.Example{}, fmt.Errorf("memory repository: %w", err)
	}

	return v, nil
}

func (repo memoryRepository) Update(ctx context.Context, id int, dom domain.Example) (domain.Example, error) {
	v, err := repo.db.Update(ctx, uint64(id), dom, cmp)
	if err != nil {
		return domain.Example{}, fmt.Errorf("memory repository: %w", err)
	}

	return v, nil
}

func (repo memoryRepository) Delete(ctx context.Context, id int) (domain.Example, error) {
	v, err := repo.db.Delete(ctx, uint64(id))
	if err != nil {
		return domain.Example{}, fmt.Errorf("memory repository: %w", err)
	}

	return v, nil
}
