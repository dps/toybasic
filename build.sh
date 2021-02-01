rm ./toybasic
go run github.com/blynn/nex lexer.nex
go run golang.org/x/tools/cmd/goyacc -o toybasic.go parser.y
go build -o toybasic toybasic.go lexer.nn.go compiler.go


echo "10 LET e = 3.14\n20 PRINT e" | ./toybasic
echo
echo ----- out.go -----
cat out.go
echo ------
echo ----- Executing compiled binary -----
go run out.go