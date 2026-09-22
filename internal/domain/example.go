package domain

type Example struct {
	ID     int
	Label  string
	Hidden string
}

func (Example) IsDomain() struct{} { return struct{}{} }
