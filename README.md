# toybasic

A toy basic interpreter playground in golang. I'm making this just for fun.

Example program
```
10 PRINT "Hello, world."
20 LET x = (3 * 2) + 3
30 x = x + 1
40 IF x == 10 THEN PRINT "Ten!"
50 IF x > 20 THEN GOTO 70
60 GOTO 30
70 END
```

# The lexer.
I'm using `github.com/blynn/nex`

```
go run github.com/blynn/nex toybasic.nex
```

# The parser.
I intend to use `goyacc`

`go run golang.org/x/tools/cmd/goyacc -o toybasic.go -p ToyBasic toybasic.y && go build toybasic.go toybasic.nn.go`