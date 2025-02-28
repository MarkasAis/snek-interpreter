package ast

import "bytes"

const NIL_STRING = "<nil>"

func safeString(n Node) string {
	if n == nil {
		return NIL_STRING
	}
	return n.String()
}

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
	out.WriteString(safeString(n.Target))
	out.WriteString(" = ")
	out.WriteString(safeString(n.Value))
	return out.String()
}

type InfixNode struct {
	Left     Node
	Operator string
	Right    Node
}

func (n *InfixNode) String() string {
	var out bytes.Buffer
	out.WriteString(safeString(n.Left))
	out.WriteString(" ")
	out.WriteString(n.Operator)
	out.WriteString(" ")
	out.WriteString(safeString(n.Right))
	return out.String()
}
