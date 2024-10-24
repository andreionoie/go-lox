package main

import (
	"fmt"
	"os"
	"slices"
)

const tokenizeCommand = "tokenize"
const parseCommand = "parse"
const evaluateCommand = "evaluate"
const runCommand = "run"

var allowedCommands = []string{tokenizeCommand, parseCommand, evaluateCommand, runCommand}

var LoxHadError = false
var LoxHadRuntimeError = false

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

		prettyPrinter := &AstPrettyPrinter{}
		interpreter := &AstInterpreter{}

		switch command {
		case tokenizeCommand:
			for _, tok := range tokens {
				fmt.Println(tok)
			}
		case parseCommand:
			expr, err := parser.ParseExpr()
			if err != nil {
				fmt.Fprint(os.Stderr, err.Error())
				break
			}
			prettyPrinter.PrintExpr(expr)
		case evaluateCommand:
			expr, err := parser.ParseExpr()
			if err != nil {
				fmt.Fprint(os.Stderr, err.Error())
				break
			}
			interpreter.InterpretExpr(expr)
		case runCommand:
			stmts, err := parser.Parse()
			if err != nil {
				fmt.Fprint(os.Stderr, err.Error())
				break
			}
			interpreter.Interpret(stmts)
		}

		if LoxHadError {
			os.Exit(65)
		}
		if LoxHadRuntimeError {
			os.Exit(70)
		}
	}
}
