package mapper

// A Domain is an interface to represent a Domain with arbitrary
// method to dinstinguate from any in type signatures.
type Domain interface {
	IsDomain() struct{}
}

// A DTO is an interface to represent a data transfer object.
// It has methods to transform a Domain object to its DTO
// representation and from DTO to its Domain counterpart.
type DTO[D Domain, T DTO[D, T]] interface {
	FromDomain(D) T
	ToDomain() D
}

// A Mapper is an object that provides methods to easily map
// between a Domain object, its DTO equivalent and vice versa.
// It also provides methods to map between slices of both types.
type Mapper[D Domain, T DTO[D, T]] interface {
	FromDomain(D) T
	ToDomain(T) D
	FromDomainList([]D) []T
	ToDomainList([]T) []D
}

// A mapper is the concrete type that implements Mapper.
type mapper[D Domain, T DTO[D, T]] struct{}

// New returns a Mapper for types Domain and DTO.
func New[D Domain, T DTO[D, T]]() Mapper[D, T] {
	return mapper[D, T]{}
}

// FromDomain transforms a Domain object to a DTO object.
// It uses DTO's FromDomain method to do the transformation.
func (m mapper[D, T]) FromDomain(dom D) T {
	var dto T
	return dto.FromDomain(dom)
}

// ToDomain transforms a DTO object to a Domain object.
// It uses DTO's ToDomain method to do the transformation.
func (m mapper[D, T]) ToDomain(dto T) D {
	return dto.ToDomain()
}

// FromDomainList transforms a slice of Domain object to a slice
// of DTO object. It uses DTO's FromDomain method to transform
// individual objects.
func (m mapper[D, T]) FromDomainList(doms []D) []T {
	dtos := make([]T, 0, len(doms))
	var dto T
	for _, dom := range doms {
		dtos = append(dtos, dto.FromDomain(dom))
	}

	return dtos
}

// ToDomainList transforms a slice of DTO object to a slice
// of Domain object. It uses DTO's ToDomain method to transform
// individual objects.
func (m mapper[D, T]) ToDomainList(dtos []T) []D {
	doms := make([]D, 0, len(dtos))
	for _, dto := range dtos {
		doms = append(doms, dto.ToDomain())
	}

	return doms
}
