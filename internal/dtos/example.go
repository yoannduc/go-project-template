package dtos

import "github.com/yoannduc/go-project-template/internal/domain"

type Example struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
}

func (dto Example) IsZero() bool {
	if dto.ID <= 0 && dto.Label == "" {
		return true
	}

	return false
}

func (dto Example) FromDomain(domain domain.Example) Example {
	dto.ID = domain.ID
	dto.Label = domain.Label

	return dto
}

func (dto Example) ToDomain() domain.Example {
	return domain.Example{
		ID:     dto.ID,
		Label:  dto.Label,
		Hidden: "hidden",
	}
}
