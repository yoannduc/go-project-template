package dtos

import "github.com/yoannduc/go-project-template/internal/domain"

// Example is a data transfer object representation of an example.
type Example struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
}

// IsZero returns whether the Example is considered of zero value.
func (dto Example) IsZero() bool {
	if dto.ID <= 0 && dto.Label == "" {
		return true
	}

	return false
}

// FromDomain implements mapper.FromDomain.
func (dto Example) FromDomain(domain domain.Example) Example {
	dto.ID = domain.ID
	dto.Label = domain.Label

	return dto
}

// ToDomain implements mapper.ToDomain.
func (dto Example) ToDomain() domain.Example {
	return domain.Example{
		ID:     dto.ID,
		Label:  dto.Label,
		Hidden: "hidden",
	}
}
