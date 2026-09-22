package dtos

import "github.com/yoannduc/go-project-template/internal/domain"

type DTO[dom domain.Domain, dto DTO[dom, dto]] interface {
	FromDomain(dom) dto
	ToDomain() dom
}
