package sort

const (
	DirectionAscending = iota
	DirectionDescending
)

type Order struct {
	Field     string
	FieldType string
	Direction int
}
