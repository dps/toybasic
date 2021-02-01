rm ./toybasic
go run github.com/blynn/nex lexer.nex
go run golang.org/x/tools/cmd/goyacc -o toybasic.go parser.y
go build -o toybasic toybasic.go lexer.nn.go compiler.go


echo "10 PRINT 2 * 2 + 3" | ./toybasic
echo
echo ----- out.go -----
cat out.go
echo ------
echo ----- Executing compiled binary -----
go run out.go