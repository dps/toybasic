# toybasic

This is a toy compiler for a simple dialect of BASIC in golang. I made this just for fun.

Example program
```VB
10 PRINT "Hello, world."
20 LET x = (3 * 2) + 3
30 LET x = x + 1
40 IF x == 10 THEN PRINT "Ten!"
45 PRINT x
50 IF x >= 15 THEN GOTO 70
60 GOTO 30
70 END
```
Example output
```
$ ./toybasic <hello.bas
$ go run out.go
Hello, world.
Ten!
10
11
12
13
14
15
```
# The lexer
The lexer uses `github.com/blynn/nex`. I really like `nex`'s awk-like syntax.
```
$ go run github.com/blynn/nex lexer.nex
```

# The parser
The parser uses `goyacc`
```
$ go run golang.org/x/tools/cmd/goyacc -o toybasic.go parser.y
```

# Build the compiler
```
$ go build -o toybasic toybasic.go lexer.nn.go compiler.go
```