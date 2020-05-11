## Writing an interpreter in Go
golang project that implements an interpreter for the Monkey programming language using book written by Thorsten Ball

Clone project
```bash
git clone https://github.com/tewshi/interpreter-with-go.git
```

then
```bash
go run main.go
```

this will spawn the REPL console for testing


to build the app binaary
```bash
go build -o bin/app
```

this will build the app for the platform.


to run the built binary
```bash
./bin/app
```

this will spawn the REPL console for coding


to run tests
```bash
go test ./evaluator ./lexer ./object ./ast ./parser
```

this will run all the tests