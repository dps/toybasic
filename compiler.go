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
	Generate()
}

func ex(op Node, lineNum int) Line {
	fmt.Fprintf(writer, "goto label_%d;", lineNum)
	fmt.Fprintf(writer, "label_%d:", lineNum)
	fmt.Fprintln(writer, " // line ", lineNum)

	op.Generate()
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

func (op Op) Generate() {
	fmt.Fprint(writer, op.operator)
}

type LetOp struct {
	opType     int
	variable   VarOp
	expression Node
}

func (op LetOp) Generate() {
	regNum := (strings.ToUpper(op.variable.VariableName())[0] - 'A')
	fmt.Fprintf(writer, "registers[%d] = ", regNum)
	op.expression.Generate()
	fmt.Fprintln(writer)
}

type EndOp struct {
	opType int
}

func (op EndOp) Generate() {
	fmt.Fprintln(writer, "goto end")
}

type RelOp struct {
	opType int
}

func (op RelOp) Generate() {
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

func (op ConditionalOp) Generate() {
	fmt.Fprintf(writer, `if (`)
	op.left.Generate()
	op.relop.Generate()
	op.right.Generate()
	fmt.Fprintln(writer, `) {`)
	op.conditionalExpression.Generate()
	fmt.Fprintln(writer, `}`)
}

type VarOp struct {
	opType   int
	variable string
}

func (op VarOp) VariableName() string {
	return op.variable
}
func (op VarOp) Generate() {
	regNum := (strings.ToUpper(op.variable)[0] - 'A')
	fmt.Fprintf(writer, "registers[%d]", regNum)
}

type GotoOp struct {
	opType     int
	expression Node
}

func (op GotoOp) Generate() {
	fmt.Fprintf(writer, "goto label_")
	op.expression.Generate()
	fmt.Fprintf(writer, ";")
}

type InfixOp struct {
	opType   int
	left     Node
	right    Node
	operator string
}

func (op InfixOp) Generate() {
	op.left.Generate()
	fmt.Fprint(writer, op.operator)
	op.right.Generate()
}

type GroupingOp struct {
	opType     int
	expression Node
}

func (op GroupingOp) Generate() {
	fmt.Fprint(writer, "(")
	op.expression.Generate()
	fmt.Fprint(writer, ")")
}

type ListOp struct {
	opType int
	head   Node
	tail   Node
}

func (op ListOp) Generate() {
	// "," is only used inside PRINT statements
	op.head.Generate()
	fmt.Fprint(writer, ", ")
	op.tail.Generate()
}

type StringOp struct {
	opType int
	val    string
}

func (op StringOp) Generate() {
	fmt.Fprint(writer, op.val)
}

type IntOp struct {
	opType int
	val    int
}

func (op IntOp) Generate() {
	fmt.Fprint(writer, op.val)
}

type PrintOp struct {
	opType     int
	expression Node
}

func (op PrintOp) Generate() {
	fmt.Fprint(writer, "fmt.Println(")
	op.expression.Generate()
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
	OpenOutput()
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
