package ast

import "bytes"

type Node interface {
	// TokenLiteral() string
	String() string
}

type IdentifierNode struct {
	Name string
}

func (n *IdentifierNode) String() string { return n.Name }

type NumberNode struct {
	Value string
}

func (n *NumberNode) String() string { return n.Value }

type AssignmentNode struct {
	Target Node
	Value  Node
}

func (n *AssignmentNode) String() string {
	var out bytes.Buffer
	out.WriteString(n.Target.String())
	out.WriteString(" = ")
	out.WriteString(n.Value.String())
	return out.String()
}
