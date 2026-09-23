package example

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/yoannduc/go-project-template/internal/domain"
	"github.com/yoannduc/go-project-template/internal/ports"
	"github.com/yoannduc/go-project-template/pkg/memorydb"
)

var (
	errIDLowerZero = errors.New("ID cannot be lower than zero")
)

// cmp is the function applied to create and update to ensure
// uniqueness. It will be applied to every element on the list
// until found or end of list. a is current list item and b is
// the fixed item it will be compared to. It returns true if Label
// matches but ID does not (uniqueness on Label that will not bloc
// update on same item).
func cmp(a, b domain.Example) bool {
	if a.ID == b.ID {
		return false
	}

	return strings.EqualFold(a.Label, b.Label)
}

// Concrete type to implement ports.ExampleRepository.
type memoryRepository struct {
	db memorydb.MemoryDB[domain.Example]
}

// NewMemoryRepository returns a ports.ExampleRepository fetching
// data from a memorydb.MemoryDB instance.
func NewMemoryRepository(db memorydb.MemoryDB[domain.Example]) ports.ExampleRepository {
	return memoryRepository{
		db: db,
	}
}

// FindAll implements ports.ExampleRepository FindAll.
// It returns all instances from db.
// It uses memorydb.MemoryDB FindAll.
func (repo memoryRepository) FindAll(ctx context.Context) ([]domain.Example, error) {
	return repo.db.FindAll(ctx)
}

// FindByLabelContaining implements ports.ExampleRepository
// FindByLabelContaining. It returns every Example whose Label's
// value contains search. It uses memorydb.MemoryDB FindFilterFunc
// applying closure filter function.
func (repo memoryRepository) FindByLabelContaining(ctx context.Context, search string) ([]domain.Example, error) {
	v, err := repo.db.FindFilterFunc(ctx, func(dom domain.Example) bool {
		return strings.Contains(dom.Label, search)
	})
	if err != nil {
		return nil, fmt.Errorf("memory repository: %w", err)
	}

	return v, nil
}

// FindByID implements ports.ExampleRepository FindByID.
// It uses memorydb.MemoryDB FindByID.
func (repo memoryRepository) FindByID(ctx context.Context, id int) (domain.Example, error) {
	if id < 0 {
		return domain.Example{}, fmt.Errorf("memory repository: %w", errIDLowerZero)
	}
	v, err := repo.db.FindByID(ctx, uint64(id))
	if err != nil {
		return domain.Example{}, fmt.Errorf("memory repository: %w", err)
	}

	return v, nil
}

// Create implements ports.ExampleRepository Create.
// It returns created Example.
// It uses memorydb.MemoryDB Create applying cmp for uniqueness check.
func (repo memoryRepository) Create(ctx context.Context, dom domain.Example) (domain.Example, error) {
	v, err := repo.db.Create(ctx, dom, cmp)
	if err != nil {
		return domain.Example{}, fmt.Errorf("memory repository: %w", err)
	}

	return v, nil
}

// Update implements ports.ExampleRepository Update.
// It returns updated Example.
// It uses memorydb.MemoryDB Update applying cmp for uniqueness check.
func (repo memoryRepository) Update(ctx context.Context, id int, dom domain.Example) (domain.Example, error) {
	if id < 0 {
		return domain.Example{}, fmt.Errorf("memory repository: %w", errIDLowerZero)
	}
	v, err := repo.db.Update(ctx, uint64(id), dom, cmp)
	if err != nil {
		return domain.Example{}, fmt.Errorf("memory repository: %w", err)
	}

	return v, nil
}

// Delete implements ports.ExampleRepository Delete.
// It returns deleted Example.
// It uses memorydb.MemoryDB Delete.
func (repo memoryRepository) Delete(ctx context.Context, id int) (domain.Example, error) {
	if id < 0 {
		return domain.Example{}, fmt.Errorf("memory repository: %w", errIDLowerZero)
	}
	v, err := repo.db.Delete(ctx, uint64(id))
	if err != nil {
		return domain.Example{}, fmt.Errorf("memory repository: %w", err)
	}

	return v, nil
}
