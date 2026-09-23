package example

import (
	"context"
	"fmt"

	"github.com/yoannduc/go-project-template/internal/domain"
	"github.com/yoannduc/go-project-template/internal/dtos"
	"github.com/yoannduc/go-project-template/internal/ports"
	"github.com/yoannduc/go-project-template/pkg/mapper"
)

// exampleService is the concrete type to implement ports.ExampleService.
type exampleService struct {
	repo  ports.ExampleRepository
	mappr mapper.Mapper[domain.Example, dtos.Example]
}

// New returns a ports.ExampleService using inputed
// ports.ExampleRepository and mapper.Mapper.
func New(repo ports.ExampleRepository, mapper mapper.Mapper[domain.Example, dtos.Example]) ports.ExampleService {
	return exampleService{
		repo:  repo,
		mappr: mapper,
	}
}

// GetAll implements ports.ExampleService GetAll.
// It returns all Example fetched from repo, mapped
// to their dto representation. It uses ports.ExampleRepository
// FindAll and mapper.Mapper FromDomainList.
func (srv exampleService) GetAll(ctx context.Context) ([]dtos.Example, error) {
	v, err := srv.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("memory service: %w", err)
	}

	return srv.mappr.FromDomainList(v), nil
}

// GetByLabelContaining implements ports.ExampleService
// GetByLabelContaining. It returns every Example that
// matches search param, mapped to their dto representation.
// It uses ports.ExampleRepository FindByLabelContaining and
// mapper.Mapper FromDomainList.
func (srv exampleService) GetByLabelContaining(ctx context.Context, search string) ([]dtos.Example, error) {
	v, err := srv.repo.FindByLabelContaining(ctx, search)
	if err != nil {
		return nil, fmt.Errorf("memory service: %w", err)
	}

	return srv.mappr.FromDomainList(v), nil
}

// GetByID implements ports.ExampleService GetByID.
// It returns found Example if any mapped, to its dto representation.
// It uses ports.ExampleRepository FindByID and mapper.Mapper FromDomain.
func (srv exampleService) GetByID(ctx context.Context, id int) (dtos.Example, error) {
	v, err := srv.repo.FindByID(ctx, id)
	if err != nil {
		return dtos.Example{}, fmt.Errorf("memory service: %w", err)
	}

	return srv.mappr.FromDomain(v), nil
}

// Create implements ports.ExampleService Create.
// It returns created Example if any, mapped to its dto representation.
// It uses ports.ExampleRepository Create and mapper.Mapper FromDomain.
func (srv exampleService) Create(ctx context.Context, dto dtos.Example) (dtos.Example, error) {
	v, err := srv.repo.Create(ctx, srv.mappr.ToDomain(dto))
	if err != nil {
		return dtos.Example{}, fmt.Errorf("memory service: %w", err)
	}

	return srv.mappr.FromDomain(v), nil
}

// Update implements ports.ExampleService Update.
// It returns updated Example if any, mapped to its dto representation.
// It uses ports.ExampleRepository Update and mapper.Mapper FromDomain.
func (srv exampleService) Update(ctx context.Context, id int, dto dtos.Example) (dtos.Example, error) {
	v, err := srv.repo.Update(ctx, id, srv.mappr.ToDomain(dto))
	if err != nil {
		return dtos.Example{}, fmt.Errorf("memory service: %w", err)
	}

	return srv.mappr.FromDomain(v), nil
}

// Delete implements ports.ExampleService Delete.
// It returns deleted Example if any, mapped to its dto representation.
// It uses ports.ExampleRepository Delete and mapper.Mapper FromDomain.
func (srv exampleService) Delete(ctx context.Context, id int) (dtos.Example, error) {
	v, err := srv.repo.Delete(ctx, id)
	if err != nil {
		return dtos.Example{}, fmt.Errorf("memory service: %w", err)
	}

	return srv.mappr.FromDomain(v), nil
}
