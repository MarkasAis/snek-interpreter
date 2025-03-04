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
	Target   Node
	Operator string
	Value    Node
}

func (n *AssignmentNode) String() string {
	w := NewASTWriter()
	n.Write(w)
	return w.String()
}

func (n *AssignmentNode) Write(w *ASTWriter) {
	w.writeIndent()
	n.Target.Write(w)
	w.WriteString(" " + n.Operator + " ")
	n.Value.Write(w)
	w.WriteString("\n")
}

type PrefixNode struct {
	Operator string
	Right    Node
}

func (n *PrefixNode) String() string {
	w := NewASTWriter()
	n.Write(w)
	return w.String()
}

func (n *PrefixNode) Write(w *ASTWriter) {
	w.WriteString("(" + n.Operator + " ")
	n.Right.Write(w)
	w.WriteString(")")
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
	w.WriteString("(")
	n.Left.Write(w)
	w.WriteString(" " + n.Operator + " ")
	n.Right.Write(w)
	w.WriteString(")")
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

type ReturnNode struct {
	Value Node
}

func (n *ReturnNode) String() string {
	w := NewASTWriter()
	n.Write(w)
	return w.String()
}

func (n *ReturnNode) Write(w *ASTWriter) {
	w.WriteLine("return " + safeString(n.Value))
}

type ForNode struct {
	Targets Node
	Values  Node
	Body    Node
	Else    Node
}

func (n *ForNode) String() string {
	w := NewASTWriter()
	n.Write(w)
	return w.String()
}

func (n *ForNode) Write(w *ASTWriter) {
	w.WriteLine("for " + safeString(n.Targets) + " in " + safeString(n.Values) + ":")
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

type FunctionDefNode struct {
	Name   Node
	Params []Node
	Body   Node
}

func (n *FunctionDefNode) String() string {
	w := NewASTWriter()
	n.Write(w)
	return w.String()
}

func (n *FunctionDefNode) Write(w *ASTWriter) {
	w.WriteString("def " + safeString(n.Name) + "(")

	for i, param := range n.Params {
		param.Write(w)
		if i < len(n.Params)-1 {
			w.WriteString(", ")
		}
	}

	w.WriteLine("):")
	w.Indent()
	n.Body.Write(w)
	w.Dedent()
}

type ParamNode struct {
	Name         Node
	DefaultValue Node
}

func (n *ParamNode) String() string {
	w := NewASTWriter()
	n.Write(w)
	return w.String()
}

func (n *ParamNode) Write(w *ASTWriter) {
	w.WriteString(safeString(n.Name))
	if n.DefaultValue != nil {
		w.WriteString("=" + safeString(n.DefaultValue))
	}
}

type CallNode struct {
	Function Node
	Args     []Node
}

func (n *CallNode) String() string {
	w := NewASTWriter()
	n.Write(w)
	return w.String()
}

func (n *CallNode) Write(w *ASTWriter) {
	w.WriteString(safeString(n.Function) + "(")

	for i, arg := range n.Args {
		arg.Write(w)
		if i < len(n.Args)-1 {
			w.WriteString(", ")
		}
	}

	w.WriteString(")")
}

type SliceNode struct {
	Left  Node
	Index Node
}

func (n *SliceNode) String() string {
	w := NewASTWriter()
	n.Write(w)
	return w.String()
}

func (n *SliceNode) Write(w *ASTWriter) {
	w.WriteString(safeString(n.Left) + "[")
	n.Index.Write(w)
	w.WriteString("]")
}

type ExpressionsNode struct {
	Expressions []Node
}

func (n *ExpressionsNode) String() string {
	w := NewASTWriter()
	n.Write(w)
	return w.String()
}

func (n *ExpressionsNode) Write(w *ASTWriter) {
	for i, exp := range n.Expressions {
		exp.Write(w)
		if i < len(n.Expressions)-1 {
			w.WriteString(", ")
		}
	}
}
