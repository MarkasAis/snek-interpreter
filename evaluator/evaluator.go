package evaluator

import (
	"fmt"
	"snek/ast"
	"strings"
)

const indentString = "  "

func identLevel(indent int) string {
	return strings.Repeat(indentString, indent)
}

func indentPrint(fs string, indent int) {
	fmt.Printf("%s%s\n", identLevel(indent), fs)
}

func DebugPrint(node ast.Node, depth int) {
	switch n := node.(type) {
	case *ast.BlockNode:
		indentPrint("block", depth)
		DebugPrintAll(n.Statements, depth+1)
	case *ast.ExpressionsNode:
		indentPrint("expressions", depth)
		DebugPrintAll(n.Expressions, depth+1)
	case *ast.NumberNode:
		indentPrint("number", depth)
	case *ast.InfixNode:
		indentPrint("infix", depth)
		DebugPrint(n.Left, depth+1)
		DebugPrint(n.Right, depth+1)
	case *ast.PrefixNode:
		indentPrint("prefix", depth)
		DebugPrint(n.Right, depth+1)
	case *ast.AssignmentNode:
		indentPrint("assignment", depth)
		DebugPrint(n.Target, depth+1)
		DebugPrint(n.Value, depth+1)
	case *ast.IdentifierNode:
		indentPrint("identifier", depth)
	default:
		indentPrint("?", depth)
	}
}

func DebugPrintAll(nodes []ast.Node, depth int) {
	for _, n := range nodes {
		DebugPrint(n, depth)
	}
}
