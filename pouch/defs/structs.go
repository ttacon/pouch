package defs

type StructInfo struct {
	Name              string
	Table             string
	Fields            []FieldInfo
	IDField           string
	HasAutoGenIDField bool
}

type FieldInfo struct {
	Name         string
	Column       string
	IsPrimaryKey bool
	IsPointer    bool
	Type         string
}
