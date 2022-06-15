package redis

type Mode int

const (
	Read = iota
	Write
	Misc
	Collaborate
)

func (m Mode) String() string {
	return [...]string{"read", "write", "misc", "collabrate"}[m]
}

type Database int

const (
	Project Database = iota
	Mapping
)
