%{
package main
import "fmt"

var registers = make([]int, 26)
var line = 0

type Op struct {
    opType int;
    operands []interface{};
}

type Line struct {
    lineNum int;
    statement Op;
}


func ex(op Op, lineNum int) Line {
    fmt.Printf("Line", lineNum);
    return Line{lineNum, op};
}

func opr(op int, nargs int, args ...interface{}) Op {
    fmt.Printf("Op", op, nargs, args);
    return Op{op, nil}
}

func variable(name rune) Op {
    fmt.Printf("Variable", name);
    return Op{VARIABLE, nil}
}

func integer(val int) Op {
    fmt.Printf("Integer", val)
    return Op{INTEGER, nil}
}

func decimal(val float64) Op {
    fmt.Printf("Float", val)
    return Op{DECIMAL, nil}
}

func basString(val string) Op {
    fmt.Printf("String", val)
    return Op{STRING, nil}
}

%}

%union {
    v rune   /* Variable */
    s string /* String */
    num int  /* Integer constant. */
    dec float64  /* Decimal constant. */
    node Op /* Node object. */
};

%token <num> INTEGER
%token <s> STRING
%token <v> VARIABLE
%token <dec> DECIMAL

%token PRINT IF GOTO LET END THEN CR

%left LT LE GT GE EQ NE 
%left '+' '-'
%left '*' '/'

%type <node> line statement expression term factor number v s
%type <node> relop expr_list

%%

program:
    block                               {}
    ;

block:
    block line                          {}
    | line                              {}
    ;

line:
    INTEGER statement CR                { ex($2,$1); line = $1;}
    ;

statement:
    PRINT expr_list                     { $$ = opr(PRINT, 1, $2); }
    | IF expression relop expression THEN statement { $$ = opr(IF, 4, $2, $3, $4, $6); }
    | GOTO expression                   { $$ = opr(GOTO, 1, $2); }
    | LET v '=' expression            { $$ = opr(LET, 2, $2, $4); }
    | END                               { $$ = opr(END, 0);  }
    ;

expr_list:
    expr_list ','  expression           { $$ = opr('l', 2, $1, $3); }
    | expression                        { $$ = $1; }
    ;

expression:
    expression '+' term                 { $$ = opr('+', 2, $1, $3); }
    | expression '-' term               { $$ = opr('-', 2, $1, $3); }
    | term                              { $$ = $1; }
    | s                                 { $$ = $1; }
    ;

term:
    term '*' factor                     { $$ = opr('*', 2, $1, $3); }
    | term '/' factor                   { $$ = opr('/', 2, $1, $3); }
    | factor                            { $$ = $1; }
    ;

factor:
    v                                 { $$ = $1; }
    | number                            { $$ = $1; }
    | '(' expression ')'                { $$ = opr('(', 1, $2); }
    ;

number:
    INTEGER                             { $$ = integer($1); }
    | DECIMAL                           { $$ = decimal($1); }
    ;

v:
    VARIABLE                            { $$ = variable($1); }
    ;

s:
    STRING                              { $$ = basString($1);}
    ;

relop:
    LT                                  { $$ = opr(LT, 0); }
    | LE                                { $$ = opr(LE, 0); }
    | GT                                { $$ = opr(GT, 0); }
    | GE                                { $$ = opr(GE, 0); }
    | EQ                                { $$ = opr(EQ, 0); }
    | NE                                { $$ = opr(NE, 0); }
    ;

%%

