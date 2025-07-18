package data

type Field struct {
	Name       string
	Type       string
	IsReadOnly bool
	IsHidden   bool
}

type Columns []Field
