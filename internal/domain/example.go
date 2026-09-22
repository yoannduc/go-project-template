package domain

// Example is the domain representation of an example.
type Example struct {
	ID     int
	Label  string
	Hidden string
}

// IsDomain implements mapper.Domain.
func (Example) IsDomain() struct{} { return struct{}{} }
