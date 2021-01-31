package main

import (
	"bufio"
	"fmt"
	"os"
)

var outfile *os.File
var writer *bufio.Writer

type Node interface {
	Execute()
	Type() int
}

type Op struct {
	opType   int
	operands []Op
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
	fmt.Fprintln(writer, "package toybasic")
	fmt.Fprint(writer, `
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
