package mapper

import (
	"github.com/yoannduc/go-project-template/internal/domain"
	"github.com/yoannduc/go-project-template/internal/dtos"
)

type Mapper[D domain.Domain, T dtos.DTO[D, T]] interface {
	FromDomain(D) T
	ToDomain(T) D
	FromDomainList([]D) []T
	ToDomainList([]T) []D
}

type mapper[D domain.Domain, T dtos.DTO[D, T]] struct{}

func New[D domain.Domain, T dtos.DTO[D, T]]() Mapper[D, T] {
	return mapper[D, T]{}
}

func (m mapper[D, T]) FromDomain(dom D) T {
	var dto T
	return dto.FromDomain(dom)
}

func (m mapper[D, T]) ToDomain(dto T) D {
	return dto.ToDomain()
}

func (m mapper[D, T]) FromDomainList(doms []D) []T {
	dtos := make([]T, 0, len(doms))
	var dto T
	for _, dom := range doms {
		dtos = append(dtos, dto.FromDomain(dom))
	}

	return dtos
}

func (m mapper[D, T]) ToDomainList(dtos []T) []D {
	doms := make([]D, 0, len(dtos))
	for _, dto := range dtos {
		doms = append(doms, dto.ToDomain())
	}

	return doms
}
