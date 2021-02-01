package main

import (
	"bufio"
	"fmt"
	"os"
)

var outfile *os.File
var writer *bufio.Writer

type Line struct {
	lineNum   int
	statement Node
}
type Node interface {
	Execute()
	Type() int
}

// Generic Op is used for pass through operators (e.g. math) that
// work the same way in BASIC as they do in Go.
type Op struct {
	opType   int
	operator string
}

func (op Op) Type() int {
	return op.opType
}
func (op Op) Execute() {
	fmt.Fprint(writer, op.operator)
}

type InfixOp struct {
	opType   int
	left     Node
	right    Node
	operator string
}

func (op InfixOp) Type() int {
	return op.opType
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

func (op GroupingOp) Type() int {
	return op.opType
}
func (op GroupingOp) Execute() {
	fmt.Fprint(writer, "(")
	op.expression.Execute()
	fmt.Fprint(writer, ")")
}

type StringOp struct {
	opType int
	val    string
}

func (op StringOp) Type() int {
	return op.opType
}
func (op StringOp) Execute() {
	fmt.Fprint(writer, op.val)
}

type IntOp struct {
	opType int
	val    int
}

func (op IntOp) Type() int {
	return op.opType
}
func (op IntOp) Execute() {
	fmt.Fprint(writer, op.val)
}

type PrintOp struct {
	opType     int
	expression Node
}

func (op PrintOp) Type() int {
	return op.opType
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
	var registers = make([]int, 26)
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
