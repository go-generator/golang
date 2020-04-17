package model

type FileInfo struct {
	Name       string
	StructName string
	Fields     []FieldInfo
	IDFields   []FieldInfo
}
