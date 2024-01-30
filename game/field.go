package game

type Field struct {
	// Type of this field
	Type int
	// Tile on this field, nil if empty
	Tile *Tile
}

func NewField(fieldType int) *Field {
	field := &Field{
		Type: fieldType,
	}
	return field
}
