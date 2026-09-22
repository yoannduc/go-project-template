package example

import (
	"context"
	"fmt"

	"github.com/yoannduc/go-project-template/internal/domain"
	"github.com/yoannduc/go-project-template/internal/dtos"
	"github.com/yoannduc/go-project-template/internal/ports"
	"github.com/yoannduc/go-project-template/pkg/mapper"
)

type exampleService struct {
	repo  ports.ExampleRepository
	mappr mapper.Mapper[domain.Example, dtos.Example]
}

func New(repo ports.ExampleRepository, mapper mapper.Mapper[domain.Example, dtos.Example]) ports.ExampleService {
	return exampleService{
		repo:  repo,
		mappr: mapper,
	}
}

func (srv exampleService) GetAll(ctx context.Context) ([]dtos.Example, error) {
	v, err := srv.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("memory service: %w", err)
	}

	return srv.mappr.FromDomainList(v), nil
}

func (srv exampleService) GetByLabelContaining(ctx context.Context, search string) ([]dtos.Example, error) {
	v, err := srv.repo.FindByLabelContaining(ctx, search)
	if err != nil {
		return nil, fmt.Errorf("memory service: %w", err)
	}

	return srv.mappr.FromDomainList(v), nil
}

func (srv exampleService) GetByID(ctx context.Context, id int) (dtos.Example, error) {
	v, err := srv.repo.FindByID(ctx, id)
	if err != nil {
		return dtos.Example{}, fmt.Errorf("memory service: %w", err)
	}

	return srv.mappr.FromDomain(v), nil
}

func (srv exampleService) Post(ctx context.Context, dto dtos.Example) (dtos.Example, error) {
	v, err := srv.repo.Create(ctx, srv.mappr.ToDomain(dto))
	if err != nil {
		return dtos.Example{}, fmt.Errorf("memory service: %w", err)
	}

	return srv.mappr.FromDomain(v), nil
}

func (srv exampleService) Patch(ctx context.Context, id int, dto dtos.Example) (dtos.Example, error) {
	v, err := srv.repo.Update(ctx, id, srv.mappr.ToDomain(dto))
	if err != nil {
		return dtos.Example{}, fmt.Errorf("memory service: %w", err)
	}

	return srv.mappr.FromDomain(v), nil
}

func (srv exampleService) Delete(ctx context.Context, id int) (dtos.Example, error) {
	v, err := srv.repo.Delete(ctx, id)
	if err != nil {
		return dtos.Example{}, fmt.Errorf("memory service: %w", err)
	}

	return srv.mappr.FromDomain(v), nil
}
