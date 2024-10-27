package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"io"
	"os"
	"text/template"
)

type Grammar struct {
	AST_DEFINITIONS []AstDefinition `json:"astDefinitions"`
}

type AstDefinition struct {
	BASE_NAME   string       `json:"baseName"`
	PRODUCTIONS []Production `json:"productions"`
}

type Production struct {
	HEAD string     `json:"head"`
	BODY []BodyItem `json:"body"`
}

type BodyItem struct {
	TYPE string `json:"type"`
	NAME string `json:"name"`
}

const codeTemplate = `
// Code generated by ast_codegen.go; DO NOT EDIT.

package main
import "errors"

{{ range .AST_DEFINITIONS }}
{{ $L_AST_DEFINITION := . }}
{{ $L_VISITOR_INTERFACE := print $L_AST_DEFINITION.BASE_NAME "Visitor" }}
	// define the base {{ $L_AST_DEFINITION.BASE_NAME }} (5.2.2 Metaprogramming the trees)
	type {{ $L_AST_DEFINITION.BASE_NAME }} interface {
		// define the abstract accept() function (5.3.3 Visitors for expressions)
		Accept(visitor {{ $L_VISITOR_INTERFACE }}) (result interface{}, err error)
	}

	// define the visitor interface (5.3.3 Visitors for expressions)
	type {{ $L_VISITOR_INTERFACE }} interface {
		{{ range $L_AST_DEFINITION.PRODUCTIONS }}
		{{ $L_PRODUCTION_NAME := print .HEAD $L_AST_DEFINITION.BASE_NAME }}
			Visit{{ $L_PRODUCTION_NAME }}(v *{{ $L_PRODUCTION_NAME }}) (result interface{}, err error)
		{{ end }}
	}

	type Stub{{ $L_VISITOR_INTERFACE }} struct{}
	// type assertion to ensure stub implements all
	var _ {{ $L_VISITOR_INTERFACE }} = Stub{{ $L_VISITOR_INTERFACE }} {}

	{{ range $L_AST_DEFINITION.PRODUCTIONS }}
	{{ $L_PRODUCTION_NAME := print .HEAD $L_AST_DEFINITION.BASE_NAME }}
		func (s Stub{{ $L_VISITOR_INTERFACE }}) Visit{{ $L_PRODUCTION_NAME }}(_ *{{ $L_PRODUCTION_NAME }}) (result interface{}, err error) {
			return nil, errors.New("visit func for {{ $L_PRODUCTION_NAME }} is not implemented")
		}
	{{ end }}

	{{ range $L_AST_DEFINITION.PRODUCTIONS }}
	{{ $L_PRODUCTION_NAME := print .HEAD $L_AST_DEFINITION.BASE_NAME }}
		// define the subtype {{ .HEAD }} (5.2.2 Metaprogramming the trees)
		type {{ $L_PRODUCTION_NAME }} struct {
			{{ range .BODY }}
				{{ .NAME }} {{ .TYPE }}
			{{ end }}
		}

		// each subtype implements the abstract accept() and calls the right visit method (5.3.3 Visitors for expressions)
		func (b *{{ $L_PRODUCTION_NAME }}) Accept(visitor {{ $L_VISITOR_INTERFACE }}) (result interface{}, err error) {
			return visitor.Visit{{ $L_PRODUCTION_NAME }}(b)
		}

		var _ {{ $L_AST_DEFINITION.BASE_NAME }} = (* {{ $L_PRODUCTION_NAME }} ) (nil)
	{{ end }}

{{ end }}
`

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a grammar file path as an argument")
		return
	}

	grammarFile, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer grammarFile.Close()

	grammarFileBytes, err := io.ReadAll(grammarFile)
	if err != nil {
		panic(err)
	}

	var grammar Grammar
	err = json.Unmarshal(grammarFileBytes, &grammar)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Grammar: %+v\n", grammar)

	tmpl, err := template.New("myTempl").Parse(codeTemplate)
	if err != nil {
		panic(err)
	}

	var codegenBuf bytes.Buffer
	err = tmpl.Execute(&codegenBuf, grammar)
	if err != nil {
		panic(err)
	}

	codegenFormatted, err := format.Source(codegenBuf.Bytes())
	if err != nil {
		panic(err)
	}

	codegenOut, err := os.Create("ast.generated.go")
	if err != nil {
		panic(err)
	}
	defer codegenOut.Close()

	_, err = codegenOut.Write(codegenFormatted)
	if err != nil {
		panic(err)
	}
}
