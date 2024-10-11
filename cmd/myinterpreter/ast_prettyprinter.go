package main

import (
	"fmt"
	"strings"
)

// AstPrettyPrinter struct implements the ExprVisitor interface.
type AstPrettyPrinter struct {
	// embed the stub to implement all functions
	StubExprVisitor
}

func (s *AstPrettyPrinter) VisitBinaryExpr(e *BinaryExpr) (result interface{}, err error) {
	return s.parenthesize(e.operator.Lexeme, e.left, e.right), nil
}

func (s *AstPrettyPrinter) VisitUnaryExpr(e *UnaryExpr) (result interface{}, err error) {
	return s.parenthesize(e.operator.Lexeme, e.right), nil
}

func (s *AstPrettyPrinter) VisitGroupingExpr(e *GroupingExpr) (result interface{}, err error) {
	return s.parenthesize("group", e.expr), nil
}

func (s *AstPrettyPrinter) VisitLiteralExpr(e *LiteralExpr) (result interface{}, err error) {
	return fmt.Sprintf("%v", e.value), nil
}

func (p *AstPrettyPrinter) parenthesize(name string, exprs ...Expr) string {
	var builder strings.Builder

	builder.WriteString("(")
	builder.WriteString(name)
	for _, expr := range exprs {
		builder.WriteString(" ")
		result, err := expr.Accept(p)
		if err != nil {
			result = "EROR"
		}
		builder.WriteString(fmt.Sprintf("%v", result))
	}
	builder.WriteString(")")

	return builder.String()
}
