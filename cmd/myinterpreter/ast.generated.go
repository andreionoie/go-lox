// Code generated by ast_codegen.go; DO NOT EDIT.

package main

import "errors"

// define the base Expr (5.2.2 Metaprogramming the trees)
type Expr interface {
	// define the abstract accept() function (5.3.3 Visitors for expressions)
	Accept(visitor ExprVisitor) (result interface{}, err error)
}

// define the visitor interface (5.3.3 Visitors for expressions)
type ExprVisitor interface {
	VisitBinaryExpr(v *BinaryExpr) (result interface{}, err error)

	VisitUnaryExpr(v *UnaryExpr) (result interface{}, err error)

	VisitGroupingExpr(v *GroupingExpr) (result interface{}, err error)

	VisitLiteralExpr(v *LiteralExpr) (result interface{}, err error)
}

type StubExprVisitor struct{}

// type assertion to ensure stub implements all
var _ ExprVisitor = StubExprVisitor{}

func (s StubExprVisitor) VisitBinaryExpr(_ *BinaryExpr) (result interface{}, err error) {
	return nil, errors.New("visit func for BinaryExpr is not implemented")
}

func (s StubExprVisitor) VisitUnaryExpr(_ *UnaryExpr) (result interface{}, err error) {
	return nil, errors.New("visit func for UnaryExpr is not implemented")
}

func (s StubExprVisitor) VisitGroupingExpr(_ *GroupingExpr) (result interface{}, err error) {
	return nil, errors.New("visit func for GroupingExpr is not implemented")
}

func (s StubExprVisitor) VisitLiteralExpr(_ *LiteralExpr) (result interface{}, err error) {
	return nil, errors.New("visit func for LiteralExpr is not implemented")
}

// define the subtype Binary (5.2.2 Metaprogramming the trees)
type BinaryExpr struct {
	left Expr

	operator Token

	right Expr
}

// each subtype implements the abstract accept() and calls the right visit method (5.3.3 Visitors for expressions)
func (b *BinaryExpr) Accept(visitor ExprVisitor) (result interface{}, err error) {
	return visitor.VisitBinaryExpr(b)
}

var _ Expr = (*BinaryExpr)(nil)

// define the subtype Unary (5.2.2 Metaprogramming the trees)
type UnaryExpr struct {
	operator Token

	right Expr
}

// each subtype implements the abstract accept() and calls the right visit method (5.3.3 Visitors for expressions)
func (b *UnaryExpr) Accept(visitor ExprVisitor) (result interface{}, err error) {
	return visitor.VisitUnaryExpr(b)
}

var _ Expr = (*UnaryExpr)(nil)

// define the subtype Grouping (5.2.2 Metaprogramming the trees)
type GroupingExpr struct {
	expr Expr
}

// each subtype implements the abstract accept() and calls the right visit method (5.3.3 Visitors for expressions)
func (b *GroupingExpr) Accept(visitor ExprVisitor) (result interface{}, err error) {
	return visitor.VisitGroupingExpr(b)
}

var _ Expr = (*GroupingExpr)(nil)

// define the subtype Literal (5.2.2 Metaprogramming the trees)
type LiteralExpr struct {
	value interface{}
}

// each subtype implements the abstract accept() and calls the right visit method (5.3.3 Visitors for expressions)
func (b *LiteralExpr) Accept(visitor ExprVisitor) (result interface{}, err error) {
	return visitor.VisitLiteralExpr(b)
}

var _ Expr = (*LiteralExpr)(nil)
