package models

type Tree struct {
	ID     uint
	Levels []Level
}

type Level struct {
	Number uint
	Nodes  []Node
}

type Node struct {
	ID         uint
	Name       string
	AncestorID uint
}
