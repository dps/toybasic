%{
package main
import "fmt"

var line = 0


func integer(val int) Node {
    fmt.Printf("Integer", val)
    return IntOp{INTEGER, val}
}

func decimal(val float64) Op {
    fmt.Printf("Float", val)
    return Op{DECIMAL, fmt.Sprintf("%f", val)}
}

func basString(val string) Node {
    fmt.Printf("String", val)
    return StringOp{STRING, val}
}

%}

%union {
    v string   /* Variable */
    s string /* String */
    num int  /* Integer constant. */
    dec float64  /* Decimal constant. */
    node Node /* Node object. */
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
    VARIABLE                            { $$ = VarOp{VARIABLE, $1}; }
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

