package ast

import (
	"bytes"
	"strings"
)

const (
	NIL_STRING = "<nil>"
	INDENT     = "    "
)

func safeString(n Node) string {
	if n == nil {
		return NIL_STRING
	}
	return n.String()
}

type Node interface {
	String() string
	Write(w *ASTWriter)
}

type ASTWriter struct {
	out    bytes.Buffer
	indent int
}

func NewASTWriter() *ASTWriter {
	return &ASTWriter{}
}

func (w *ASTWriter) writeIndent() {
	w.out.WriteString(strings.Repeat(INDENT, w.indent))
}

func (w *ASTWriter) WriteString(s string) {
	w.out.WriteString(s)
}

func (w *ASTWriter) WriteLine(s string) {
	w.writeIndent()
	w.out.WriteString(s)
	w.out.WriteString("\n")
}

func (w *ASTWriter) Indent() {
	w.indent++
}

func (w *ASTWriter) Dedent() {
	if w.indent > 0 {
		w.indent--
	}
}

func (w *ASTWriter) String() string {
	return w.out.String()
}

type BlockNode struct {
	Statements []Node
}

func (n *BlockNode) String() string {
	w := NewASTWriter()
	n.Write(w)
	return w.String()
}

func (n *BlockNode) Write(w *ASTWriter) {
	for _, stmt := range n.Statements {
		stmt.Write(w)
	}
}

type IdentifierNode struct {
	Name string
}

func (n *IdentifierNode) String() string { return n.Name }

func (n *IdentifierNode) Write(w *ASTWriter) {
	w.WriteString(n.Name)
}

type NumberNode struct {
	Value string
}

func (n *NumberNode) String() string { return n.Value }

func (n *NumberNode) Write(w *ASTWriter) {
	w.WriteString(n.Value)
}

type AssignmentNode struct {
	Target Node
	Value  Node
}

func (n *AssignmentNode) String() string {
	w := NewASTWriter()
	n.Write(w)
	return w.String()
}

func (n *AssignmentNode) Write(w *ASTWriter) {
	w.writeIndent()
	n.Target.Write(w)
	w.WriteString(" = ")
	n.Value.Write(w)
	w.WriteString("\n")
}

type InfixNode struct {
	Left     Node
	Operator string
	Right    Node
}

func (n *InfixNode) String() string {
	w := NewASTWriter()
	n.Write(w)
	return w.String()
}

func (n *InfixNode) Write(w *ASTWriter) {
	n.Left.Write(w)
	w.WriteString(" " + n.Operator + " ")
	n.Right.Write(w)
}

type IfNode struct {
	Condition Node
	Body      Node
	Else      Node
}

func (n *IfNode) String() string {
	w := NewASTWriter()
	n.Write(w)
	return w.String()
}

func (n *IfNode) Write(w *ASTWriter) {
	w.WriteLine("if " + safeString(n.Condition) + ":")
	w.Indent()
	n.Body.Write(w)
	w.Dedent()
	if n.Else != nil {
		w.WriteLine("else:")
		w.Indent()
		n.Else.Write(w)
		w.Dedent()
	}
}

type WhileNode struct {
	Condition Node
	Body      Node
	Else      Node
}

func (n *WhileNode) String() string {
	w := NewASTWriter()
	n.Write(w)
	return w.String()
}

func (n *WhileNode) Write(w *ASTWriter) {
	w.WriteLine("while " + safeString(n.Condition) + ":")
	w.Indent()
	n.Body.Write(w)
	w.Dedent()
	if n.Else != nil {
		w.WriteLine("else:")
		w.Indent()
		n.Else.Write(w)
		w.Dedent()
	}
}

type ControlNode struct {
	Type string
}

func (n *ControlNode) String() string {
	w := NewASTWriter()
	n.Write(w)
	return w.String()
}

func (n *ControlNode) Write(w *ASTWriter) {
	w.WriteLine(n.Type)
}
