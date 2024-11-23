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

	VisitVariableExpr(v *VariableExpr) (result interface{}, err error)

	VisitAssignExpr(v *AssignExpr) (result interface{}, err error)

	VisitLogicalExpr(v *LogicalExpr) (result interface{}, err error)

	VisitCallExpr(v *CallExpr) (result interface{}, err error)
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

func (s StubExprVisitor) VisitVariableExpr(_ *VariableExpr) (result interface{}, err error) {
	return nil, errors.New("visit func for VariableExpr is not implemented")
}

func (s StubExprVisitor) VisitAssignExpr(_ *AssignExpr) (result interface{}, err error) {
	return nil, errors.New("visit func for AssignExpr is not implemented")
}

func (s StubExprVisitor) VisitLogicalExpr(_ *LogicalExpr) (result interface{}, err error) {
	return nil, errors.New("visit func for LogicalExpr is not implemented")
}

func (s StubExprVisitor) VisitCallExpr(_ *CallExpr) (result interface{}, err error) {
	return nil, errors.New("visit func for CallExpr is not implemented")
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

// define the subtype Variable (5.2.2 Metaprogramming the trees)
type VariableExpr struct {
	variableName Token
}

// each subtype implements the abstract accept() and calls the right visit method (5.3.3 Visitors for expressions)
func (b *VariableExpr) Accept(visitor ExprVisitor) (result interface{}, err error) {
	return visitor.VisitVariableExpr(b)
}

var _ Expr = (*VariableExpr)(nil)

// define the subtype Assign (5.2.2 Metaprogramming the trees)
type AssignExpr struct {
	variableName Token

	assignValue Expr
}

// each subtype implements the abstract accept() and calls the right visit method (5.3.3 Visitors for expressions)
func (b *AssignExpr) Accept(visitor ExprVisitor) (result interface{}, err error) {
	return visitor.VisitAssignExpr(b)
}

var _ Expr = (*AssignExpr)(nil)

// define the subtype Logical (5.2.2 Metaprogramming the trees)
type LogicalExpr struct {
	left Expr

	operator Token

	right Expr
}

// each subtype implements the abstract accept() and calls the right visit method (5.3.3 Visitors for expressions)
func (b *LogicalExpr) Accept(visitor ExprVisitor) (result interface{}, err error) {
	return visitor.VisitLogicalExpr(b)
}

var _ Expr = (*LogicalExpr)(nil)

// define the subtype Call (5.2.2 Metaprogramming the trees)
type CallExpr struct {
	callee Expr

	arguments []Expr

	closingParen Token
}

// each subtype implements the abstract accept() and calls the right visit method (5.3.3 Visitors for expressions)
func (b *CallExpr) Accept(visitor ExprVisitor) (result interface{}, err error) {
	return visitor.VisitCallExpr(b)
}

var _ Expr = (*CallExpr)(nil)

// define the base Stmt (5.2.2 Metaprogramming the trees)
type Stmt interface {
	// define the abstract accept() function (5.3.3 Visitors for expressions)
	Accept(visitor StmtVisitor) (result interface{}, err error)
}

// define the visitor interface (5.3.3 Visitors for expressions)
type StmtVisitor interface {
	VisitExpressionStmt(v *ExpressionStmt) (result interface{}, err error)

	VisitPrintStmt(v *PrintStmt) (result interface{}, err error)

	VisitVarStmt(v *VarStmt) (result interface{}, err error)

	VisitFunctionStmt(v *FunctionStmt) (result interface{}, err error)

	VisitReturnStmt(v *ReturnStmt) (result interface{}, err error)

	VisitBlockStmt(v *BlockStmt) (result interface{}, err error)

	VisitIfStmt(v *IfStmt) (result interface{}, err error)

	VisitWhileStmt(v *WhileStmt) (result interface{}, err error)

	VisitForStmt(v *ForStmt) (result interface{}, err error)
}

type StubStmtVisitor struct{}

// type assertion to ensure stub implements all
var _ StmtVisitor = StubStmtVisitor{}

func (s StubStmtVisitor) VisitExpressionStmt(_ *ExpressionStmt) (result interface{}, err error) {
	return nil, errors.New("visit func for ExpressionStmt is not implemented")
}

func (s StubStmtVisitor) VisitPrintStmt(_ *PrintStmt) (result interface{}, err error) {
	return nil, errors.New("visit func for PrintStmt is not implemented")
}

func (s StubStmtVisitor) VisitVarStmt(_ *VarStmt) (result interface{}, err error) {
	return nil, errors.New("visit func for VarStmt is not implemented")
}

func (s StubStmtVisitor) VisitFunctionStmt(_ *FunctionStmt) (result interface{}, err error) {
	return nil, errors.New("visit func for FunctionStmt is not implemented")
}

func (s StubStmtVisitor) VisitReturnStmt(_ *ReturnStmt) (result interface{}, err error) {
	return nil, errors.New("visit func for ReturnStmt is not implemented")
}

func (s StubStmtVisitor) VisitBlockStmt(_ *BlockStmt) (result interface{}, err error) {
	return nil, errors.New("visit func for BlockStmt is not implemented")
}

func (s StubStmtVisitor) VisitIfStmt(_ *IfStmt) (result interface{}, err error) {
	return nil, errors.New("visit func for IfStmt is not implemented")
}

func (s StubStmtVisitor) VisitWhileStmt(_ *WhileStmt) (result interface{}, err error) {
	return nil, errors.New("visit func for WhileStmt is not implemented")
}

func (s StubStmtVisitor) VisitForStmt(_ *ForStmt) (result interface{}, err error) {
	return nil, errors.New("visit func for ForStmt is not implemented")
}

// define the subtype Expression (5.2.2 Metaprogramming the trees)
type ExpressionStmt struct {
	expression Expr
}

// each subtype implements the abstract accept() and calls the right visit method (5.3.3 Visitors for expressions)
func (b *ExpressionStmt) Accept(visitor StmtVisitor) (result interface{}, err error) {
	return visitor.VisitExpressionStmt(b)
}

var _ Stmt = (*ExpressionStmt)(nil)

// define the subtype Print (5.2.2 Metaprogramming the trees)
type PrintStmt struct {
	expression Expr
}

// each subtype implements the abstract accept() and calls the right visit method (5.3.3 Visitors for expressions)
func (b *PrintStmt) Accept(visitor StmtVisitor) (result interface{}, err error) {
	return visitor.VisitPrintStmt(b)
}

var _ Stmt = (*PrintStmt)(nil)

// define the subtype Var (5.2.2 Metaprogramming the trees)
type VarStmt struct {
	varName Token

	initializerExpression Expr
}

// each subtype implements the abstract accept() and calls the right visit method (5.3.3 Visitors for expressions)
func (b *VarStmt) Accept(visitor StmtVisitor) (result interface{}, err error) {
	return visitor.VisitVarStmt(b)
}

var _ Stmt = (*VarStmt)(nil)

// define the subtype Function (5.2.2 Metaprogramming the trees)
type FunctionStmt struct {
	name Token

	parameters []Token

	body []Stmt
}

// each subtype implements the abstract accept() and calls the right visit method (5.3.3 Visitors for expressions)
func (b *FunctionStmt) Accept(visitor StmtVisitor) (result interface{}, err error) {
	return visitor.VisitFunctionStmt(b)
}

var _ Stmt = (*FunctionStmt)(nil)

// define the subtype Return (5.2.2 Metaprogramming the trees)
type ReturnStmt struct {
	expression Expr
}

// each subtype implements the abstract accept() and calls the right visit method (5.3.3 Visitors for expressions)
func (b *ReturnStmt) Accept(visitor StmtVisitor) (result interface{}, err error) {
	return visitor.VisitReturnStmt(b)
}

var _ Stmt = (*ReturnStmt)(nil)

// define the subtype Block (5.2.2 Metaprogramming the trees)
type BlockStmt struct {
	statements []Stmt
}

// each subtype implements the abstract accept() and calls the right visit method (5.3.3 Visitors for expressions)
func (b *BlockStmt) Accept(visitor StmtVisitor) (result interface{}, err error) {
	return visitor.VisitBlockStmt(b)
}

var _ Stmt = (*BlockStmt)(nil)

// define the subtype If (5.2.2 Metaprogramming the trees)
type IfStmt struct {
	condition Expr

	thenBranch Stmt

	elseBranch Stmt
}

// each subtype implements the abstract accept() and calls the right visit method (5.3.3 Visitors for expressions)
func (b *IfStmt) Accept(visitor StmtVisitor) (result interface{}, err error) {
	return visitor.VisitIfStmt(b)
}

var _ Stmt = (*IfStmt)(nil)

// define the subtype While (5.2.2 Metaprogramming the trees)
type WhileStmt struct {
	condition Expr

	loopBody Stmt
}

// each subtype implements the abstract accept() and calls the right visit method (5.3.3 Visitors for expressions)
func (b *WhileStmt) Accept(visitor StmtVisitor) (result interface{}, err error) {
	return visitor.VisitWhileStmt(b)
}

var _ Stmt = (*WhileStmt)(nil)

// define the subtype For (5.2.2 Metaprogramming the trees)
type ForStmt struct {
	init Stmt

	condition Expr

	iteration Expr

	loopBody Stmt
}

// each subtype implements the abstract accept() and calls the right visit method (5.3.3 Visitors for expressions)
func (b *ForStmt) Accept(visitor StmtVisitor) (result interface{}, err error) {
	return visitor.VisitForStmt(b)
}

var _ Stmt = (*ForStmt)(nil)
