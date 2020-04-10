package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/lexer"
	"monkey/parser"
)

// PROMPT the text that shows on the console
const PROMPT = ">> "

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
func Start(in io.Reader, out io.Writer) {
	// create the scanner
	scanner := bufio.NewScanner(in)
	for {
		// print prompt: >>
		fmt.Printf(PROMPT)
		// advance to the next token
		scanned := scanner.Scan()
		if !scanned {
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

		// print our program string
		io.WriteString(out, program.String())
		io.WriteString(out, "\n")
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
