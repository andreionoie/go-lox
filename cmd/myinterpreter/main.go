package main

import (
	"fmt"
	"os"
	"slices"
)

const tokenizeCommand = "tokenize"
const parseCommand = "parse"

var allowedCommands = []string{tokenizeCommand, parseCommand}

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
		scanner := Scanner{
			Source: []rune(string(fileContents)),
		}

		tokens := scanner.ScanTokens()

		switch command {
		case tokenizeCommand:
			for _, tok := range tokens {
				fmt.Println(tok)
			}
		case parseCommand:
			for _, tok := range tokens {
				if tok.Literal != nil {
					fmt.Print(tok.GetLiteralAsString())
				} else {
					fmt.Print(tok.Lexeme)
				}
			}
			fmt.Println()

			// Construct the expression: (-123) * (45.67)
			expression := &BinaryExpr{
				left: &UnaryExpr{
					operator: Token{
						Type:    Minus,
						Lexeme:  "-",
						Literal: nil,
						Line:    1,
					},
					right: &LiteralExpr{
						value: 123,
					},
				},
				operator: Token{
					Type:    Star,
					Lexeme:  "*",
					Literal: nil,
					Line:    1,
				},
				right: &GroupingExpr{
					expr: &LiteralExpr{
						value: 45.67,
					},
				},
			}

			printer := &AstPrettyPrinter{}
			result, _ := expression.Accept(printer)
			fmt.Println(result)
		}

		if scanner.HadErrors {
			os.Exit(65)
		}
	}
}
