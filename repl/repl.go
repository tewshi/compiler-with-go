package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"strings"
)

// PROMPT the text that shows on the console
const PROMPT = "mk>> "

// MONKEYSMILE monkey face meme
const MONKEYSMILE = `         __,__
   .--. .-" "-. .--.
/ .. \/ .-. .-. \/ .. \
| | | '| / Y \ |' | | |
|| \ \ \ 0 | 0 / / / ||
 \ '- ,\.-"""-./, -' /
  ''-' /_ ^ ^ _\ '-''
      | \._ _./ |
      \ \ '~' / /
      '._'-=-'_.'
        '-----'
`

// MONKEYFROWN monkey face meme
const MONKEYFROWN = `         __,__
   .--. .-" "-. .--.
/ .. \/ .-. .-. \/ .. \
| | | '| / Y \ |' | | |
|| \ \ \ * | * / / / ||
 \ '- ,\.-"""-./, -' /
  ''-' /_ ^ ^ _\ '-''
      \ \.~~~./ /
      | \.~~~./ |
      '._'-=-'_.'
        '-----'
`

// Start is the main function that starts this repl
func Start(in io.Reader, out io.Writer, username string) {
	// create the scanner
	scanner := bufio.NewScanner(in)

	// create the environment for storage
	env := object.NewEnvironment()

L:
	for {
		// print prompt: >>
		fmt.Printf(PROMPT)
		// advance to the next token
		scanned := scanner.Scan()
		if !scanned {
			fmt.Printf("\nGoodbye %s\n", username)
			return
		}

		// get the typed content from the scanner
		line := scanner.Text()

		// create our lexer with the input
		l := lexer.NewLexer(line)
		// create parser from lexer
		p := parser.NewParser(l)
		// parse the program fed to the parser
		program := p.ParseProgram()

		// print errors if any
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		// print our evaluated program
		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			line = strings.TrimSpace(line)
			if (line == "quit()" || line == "exit()" || line == "quit" || line == "exit") && evaluated.Type() == object.ERROROBJ {
				error := evaluated.(*object.Error)
				switch error.Message {
				case "identifier not found: quit", "identifier not found: exit":
					if line == "quit()" || line == "exit()" {
						fmt.Printf("Goodbye %s\n", username)
						break L
					}
					fmt.Printf("use %s() or Ctrl-D (i.e. EOF) to exit\n", line)
					continue L
				}
			}

			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, MONKEYFROWN)
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
