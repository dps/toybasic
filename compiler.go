package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var outfile *os.File
var writer *bufio.Writer

type Line struct {
	lineNum   int
	statement Node
}
type Node interface {
	Execute()
}

func ex(op Node, lineNum int) Line {
	fmt.Fprintf(writer, "goto label_%d;", lineNum)
	fmt.Fprintf(writer, "label_%d:", lineNum)
	fmt.Fprintln(writer, " // line ", lineNum)

	op.Execute()
	return Line{lineNum, op}
}

func opr(op int, nargs int, args ...interface{}) Node {
	if op == PRINT {
		return PrintOp{op, args[0].(Node)}
	} else if op == '+' || op == '-' || op == '*' || op == '/' {
		return InfixOp{op, args[0].(Node), args[1].(Node), string(op)}
	} else if op == '(' {
		return GroupingOp{op, args[0].(Node)}
	} else if op == 'l' {
		return ListOp{op, args[0].(Node), args[1].(Node)}
	} else if op == GOTO {
		return GotoOp{op, args[0].(Node)}
	} else if op == LET {
		return LetOp{op, args[0].(VarOp), args[1].(Node)}
	} else if op == IF {
		return ConditionalOp{op, args[0].(Node), args[2].(Node), args[1].(RelOp), args[3].(Node)}
	} else if op == END {
		return EndOp{op}
	}
	return Op{op, args[0].(string)}
}

type Op struct {
	opType   int
	operator string
}

func (op Op) Execute() {
	fmt.Fprint(writer, op.operator)
}

type LetOp struct {
	opType     int
	variable   VarOp
	expression Node
}

func (op LetOp) Execute() {
	regNum := (strings.ToUpper(op.variable.VariableName())[0] - 'A')
	fmt.Fprintf(writer, "registers[%d] = ", regNum)
	op.expression.Execute()
	fmt.Fprintln(writer)
}

type EndOp struct {
	opType int
}

func (op EndOp) Execute() {
	fmt.Fprintln(writer, "goto end")
}

type RelOp struct {
	opType int
}

func (op RelOp) Execute() {
	var relChar = map[int]string{
		GT: ">", LT: "<", LE: "<=", GE: ">=", EQ: "==", NE: "!=",
	}
	fmt.Fprintf(writer, relChar[op.opType])
}

type ConditionalOp struct {
	opType                int
	left                  Node
	right                 Node
	relop                 RelOp
	conditionalExpression Node
}

func (op ConditionalOp) Execute() {
	fmt.Fprintf(writer, `if (`)
	op.left.Execute()
	op.relop.Execute()
	op.right.Execute()
	fmt.Fprintln(writer, `) {`)
	op.conditionalExpression.Execute()
	fmt.Fprintln(writer, `}`)
}

type VarOp struct {
	opType   int
	variable string
}

func (op VarOp) VariableName() string {
	return op.variable
}
func (op VarOp) Execute() {
	regNum := (strings.ToUpper(op.variable)[0] - 'A')
	fmt.Fprintf(writer, "registers[%d]", regNum)
}

type GotoOp struct {
	opType     int
	expression Node
}

func (op GotoOp) Execute() {
	fmt.Fprintf(writer, "goto label_")
	op.expression.Execute()
}

type InfixOp struct {
	opType   int
	left     Node
	right    Node
	operator string
}

func (op InfixOp) Execute() {
	op.left.Execute()
	fmt.Fprint(writer, op.operator)
	op.right.Execute()
}

type GroupingOp struct {
	opType     int
	expression Node
}

func (op GroupingOp) Execute() {
	fmt.Fprint(writer, "(")
	op.expression.Execute()
	fmt.Fprint(writer, ")")
}

type ListOp struct {
	opType int
	head   Node
	tail   Node
}

func (op ListOp) Execute() {
	// "," is only used inside PRINT statements
	op.head.Execute()
	fmt.Fprint(writer, ", ")
	op.tail.Execute()
}

type StringOp struct {
	opType int
	val    string
}

func (op StringOp) Execute() {
	fmt.Fprint(writer, op.val)
}

type IntOp struct {
	opType int
	val    int
}

func (op IntOp) Execute() {
	fmt.Fprint(writer, op.val)
}

type PrintOp struct {
	opType     int
	expression Node
}

func (op PrintOp) Execute() {
	fmt.Fprint(writer, "fmt.Println(")
	op.expression.Execute()
	fmt.Fprintln(writer, ")")
}

func OpenOutput() {
	var err error
	outfile, err = os.Create("out.go")
	if err != nil {
		panic(err)
	}
	writer = bufio.NewWriter(outfile)
}
func WriteLeader() {
	fmt.Fprintln(writer, "package main")
	fmt.Fprint(writer, `
import "fmt"

func main() {
	var registers = make([]float64, 26)
	_ = registers

	goto start
start:
`)
}

func WriteTrailer() {
	fmt.Fprint(writer, `
	goto end
end:
}
	`)
	writer.Flush()
	outfile.Close()
}
