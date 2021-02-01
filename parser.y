%{
package main
import "fmt"
%}

%union {
    v string    /* Variable */
    s string    /* String */
    num int     /* Integer constant */
    dec float64 /* Decimal constant */
    node Node   /* Node in the AST */
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
    INTEGER statement CR                { ex($2,$1); }
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
    INTEGER                             { $$ = IntOp{INTEGER, $1}; }
    | DECIMAL                           { $$ = Op{DECIMAL, fmt.Sprintf("%f", $1)}; }
    ;

v:
    VARIABLE                            { $$ = VarOp{VARIABLE, $1}; }
    ;

s:
    STRING                              { $$ = StringOp{STRING, $1};}
    ;

relop:
    LT                                  { $$ = RelOp{LT}; }
    | LE                                { $$ = RelOp{LE}; }
    | GT                                { $$ = RelOp{GT}; }
    | GE                                { $$ = RelOp{GE}; }
    | EQ                                { $$ = RelOp{EQ}; }
    | NE                                { $$ = RelOp{NE}; }
    ;

%%