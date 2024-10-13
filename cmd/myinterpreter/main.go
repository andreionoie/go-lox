package main

import (
	"fmt"
	"os"
	"slices"
)

const tokenizeCommand = "tokenize"
const parseCommand = "parse"
const evaluateCommand = "evaluate"

var allowedCommands = []string{tokenizeCommand, parseCommand, evaluateCommand}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: ./your_program.sh %s <filename>\n", allowedCommands)
		os.Exit(1)
	}

	command := os.Args[1]

	if !slices.Contains(allowedCommands, command) {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	if len(fileContents) >= 0 {
		scanner := Scanner{Source: []rune(string(fileContents))}
		tokens := scanner.ScanTokens()
		parser := Parser{Tokens: tokens}

		switch command {
		case tokenizeCommand:
			for _, tok := range tokens {
				fmt.Println(tok)
			}
		case parseCommand:
			expr, err := parser.Parse()
			if err != nil {
				fmt.Fprint(os.Stderr, err.Error())
				break
			}
			printer := &AstPrettyPrinter{}

			result, _ := expr.Accept(printer)
			fmt.Println(result)
		case evaluateCommand:
			expr, err := parser.Parse()
			if err != nil {
				fmt.Fprint(os.Stderr, err.Error())
				break
			}
			printer := &AstInterpreter{}

			result, err := expr.Accept(printer)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(70)
			}
			if result == nil {
				fmt.Println("nil")
			} else {
				fmt.Println(result)
			}
		}

		if scanner.HadErrors || parser.HadErrors {
			os.Exit(65)
		}
	}
}
