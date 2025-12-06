package ast

import (
	"github.com/dennislee928/uad-lang/internal/common"
)

// Node is the base interface for all AST nodes
type Node interface {
	Span() common.Span
	node() // Private method to prevent external implementation
}

// ==================== Base Types ====================

// baseNode provides common functionality for all nodes
type baseNode struct {
	span common.Span
}

func (n *baseNode) Span() common.Span { return n.span }
func (n *baseNode) node()             {}

// ==================== Expressions ====================

// Expr represents an expression node
type Expr interface {
	Node
	exprNode()
}

// Ident represents an identifier
type Ident struct {
	baseNode
	Name string
}

func (i *Ident) exprNode() {}

// NewIdent creates a new identifier
func NewIdent(name string, span common.Span) *Ident {
	return &Ident{
		baseNode: baseNode{span},
		Name:     name,
	}
}

// LiteralKind represents the kind of a literal
type LiteralKind int

const (
	LitInt LiteralKind = iota
	LitFloat
	LitString
	LitBool
	LitDuration
	LitNil
)

// Literal represents a literal value
type Literal struct {
	baseNode
	Kind  LiteralKind
	Value string // String representation of the value
}

func (l *Literal) exprNode() {}

// NewLiteral creates a new literal
func NewLiteral(kind LiteralKind, value string, span common.Span) *Literal {
	return &Literal{
		baseNode: baseNode{span},
		Kind:     kind,
		Value:    value,
	}
}

// BinaryOp represents a binary operator
type BinaryOp int

const (
	// Arithmetic
	OpAdd BinaryOp = iota
	OpSub
	OpMul
	OpDiv
	OpMod
	
	// Comparison
	OpEq
	OpNeq
	OpLt
	OpGt
	OpLe
	OpGe
	
	// Logical
	OpAnd
	OpOr
)

// BinaryExpr represents a binary expression
type BinaryExpr struct {
	baseNode
	Op    BinaryOp
	Left  Expr
	Right Expr
}

func (b *BinaryExpr) exprNode() {}

// NewBinaryExpr creates a new binary expression
func NewBinaryExpr(op BinaryOp, left, right Expr, span common.Span) *BinaryExpr {
	return &BinaryExpr{
		baseNode: baseNode{span},
		Op:       op,
		Left:     left,
		Right:    right,
	}
}

// UnaryOp represents a unary operator
type UnaryOp int

const (
	OpNeg UnaryOp = iota // -
	OpNot                // !
)

// UnaryExpr represents a unary expression
type UnaryExpr struct {
	baseNode
	Op   UnaryOp
	Expr Expr
}

func (u *UnaryExpr) exprNode() {}

// NewUnaryExpr creates a new unary expression
func NewUnaryExpr(op UnaryOp, expr Expr, span common.Span) *UnaryExpr {
	return &UnaryExpr{
		baseNode: baseNode{span},
		Op:       op,
		Expr:     expr,
	}
}

// CallExpr represents a function call
type CallExpr struct {
	baseNode
	Func Expr
	Args []Expr
}

func (c *CallExpr) exprNode() {}

// NewCallExpr creates a new call expression
func NewCallExpr(fn Expr, args []Expr, span common.Span) *CallExpr {
	return &CallExpr{
		baseNode: baseNode{span},
		Func:     fn,
		Args:     args,
	}
}

// IfExpr represents an if expression
type IfExpr struct {
	baseNode
	Cond Expr
	Then *BlockExpr
	Else Expr // Can be nil, IfExpr, or BlockExpr
}

func (i *IfExpr) exprNode() {}

// NewIfExpr creates a new if expression
func NewIfExpr(cond Expr, then *BlockExpr, els Expr, span common.Span) *IfExpr {
	return &IfExpr{
		baseNode: baseNode{span},
		Cond:     cond,
		Then:     then,
		Else:     els,
	}
}

// MatchExpr represents a match expression
type MatchExpr struct {
	baseNode
	Expr Expr
	Arms []*MatchArm
}

func (m *MatchExpr) exprNode() {}

// MatchArm represents a match arm
type MatchArm struct {
	baseNode
	Pattern Pattern
	Expr    Expr
}

// NewMatchExpr creates a new match expression
func NewMatchExpr(expr Expr, arms []*MatchArm, span common.Span) *MatchExpr {
	return &MatchExpr{
		baseNode: baseNode{span},
		Expr:     expr,
		Arms:     arms,
	}
}

// NewMatchArm creates a new match arm
func NewMatchArm(pattern Pattern, expr Expr, span common.Span) *MatchArm {
	return &MatchArm{
		baseNode: baseNode{span},
		Pattern:  pattern,
		Expr:     expr,
	}
}

// BlockExpr represents a block expression
type BlockExpr struct {
	baseNode
	Stmts []Stmt
	Expr  Expr // Optional trailing expression (value of block)
}

func (b *BlockExpr) exprNode() {}

// NewBlockExpr creates a new block expression
func NewBlockExpr(stmts []Stmt, expr Expr, span common.Span) *BlockExpr {
	return &BlockExpr{
		baseNode: baseNode{span},
		Stmts:    stmts,
		Expr:     expr,
	}
}

// StructLiteral represents a struct literal
type StructLiteral struct {
	baseNode
	Name   *Ident
	Fields []*FieldInit
}

func (s *StructLiteral) exprNode() {}

// FieldInit represents a field initialization in a struct literal
type FieldInit struct {
	baseNode
	Name  *Ident
	Value Expr
}

// NewStructLiteral creates a new struct literal
func NewStructLiteral(name *Ident, fields []*FieldInit, span common.Span) *StructLiteral {
	return &StructLiteral{
		baseNode: baseNode{span},
		Name:     name,
		Fields:   fields,
	}
}

// NewFieldInit creates a new field initialization
func NewFieldInit(name *Ident, value Expr, span common.Span) *FieldInit {
	return &FieldInit{
		baseNode: baseNode{span},
		Name:     name,
		Value:    value,
	}
}

// ArrayLiteral represents an array literal
type ArrayLiteral struct {
	baseNode
	Elements []Expr
}

func (a *ArrayLiteral) exprNode() {}

// NewArrayLiteral creates a new array literal
func NewArrayLiteral(elements []Expr, span common.Span) *ArrayLiteral {
	return &ArrayLiteral{
		baseNode: baseNode{span},
		Elements: elements,
	}
}

// MapLiteral represents a map literal
type MapLiteral struct {
	baseNode
	Entries []*MapEntry
}

func (m *MapLiteral) exprNode() {}

// MapEntry represents a key-value pair in a map literal
type MapEntry struct {
	baseNode
	Key   Expr
	Value Expr
}

// NewMapLiteral creates a new map literal
func NewMapLiteral(entries []*MapEntry, span common.Span) *MapLiteral {
	return &MapLiteral{
		baseNode: baseNode{span},
		Entries:  entries,
	}
}

// NewMapEntry creates a new map entry
func NewMapEntry(key, value Expr, span common.Span) *MapEntry {
	return &MapEntry{
		baseNode: baseNode{span},
		Key:      key,
		Value:    value,
	}
}

// FieldAccess represents field access (obj.field)
type FieldAccess struct {
	baseNode
	Expr  Expr
	Field *Ident
}

func (f *FieldAccess) exprNode() {}

// NewFieldAccess creates a new field access
func NewFieldAccess(expr Expr, field *Ident, span common.Span) *FieldAccess {
	return &FieldAccess{
		baseNode: baseNode{span},
		Expr:     expr,
		Field:    field,
	}
}

// IndexExpr represents array/map indexing (arr[index])
type IndexExpr struct {
	baseNode
	Expr  Expr
	Index Expr
}

func (i *IndexExpr) exprNode() {}

// NewIndexExpr creates a new index expression
func NewIndexExpr(expr, index Expr, span common.Span) *IndexExpr {
	return &IndexExpr{
		baseNode: baseNode{span},
		Expr:     expr,
		Index:    index,
	}
}

// ParenExpr represents a parenthesized expression
type ParenExpr struct {
	baseNode
	Expr Expr
}

func (p *ParenExpr) exprNode() {}

// NewParenExpr creates a new parenthesized expression
func NewParenExpr(expr Expr, span common.Span) *ParenExpr {
	return &ParenExpr{
		baseNode: baseNode{span},
		Expr:     expr,
	}
}

// ==================== Patterns ====================

// Pattern represents a pattern for matching
type Pattern interface {
	Node
	patternNode()
}

// LiteralPattern represents a literal pattern
type LiteralPattern struct {
	baseNode
	Literal *Literal
}

func (l *LiteralPattern) patternNode() {}

// NewLiteralPattern creates a new literal pattern
func NewLiteralPattern(lit *Literal, span common.Span) *LiteralPattern {
	return &LiteralPattern{
		baseNode: baseNode{span},
		Literal:  lit,
	}
}

// IdentPattern represents an identifier pattern (binds to variable)
type IdentPattern struct {
	baseNode
	Name *Ident
}

func (i *IdentPattern) patternNode() {}

// NewIdentPattern creates a new identifier pattern
func NewIdentPattern(name *Ident, span common.Span) *IdentPattern {
	return &IdentPattern{
		baseNode: baseNode{span},
		Name:     name,
	}
}

// WildcardPattern represents a wildcard pattern (_)
type WildcardPattern struct {
	baseNode
}

func (w *WildcardPattern) patternNode() {}

// NewWildcardPattern creates a new wildcard pattern
func NewWildcardPattern(span common.Span) *WildcardPattern {
	return &WildcardPattern{
		baseNode: baseNode{span},
	}
}

// StructPattern represents a struct pattern
type StructPattern struct {
	baseNode
	Name   *Ident
	Fields []*FieldPattern
}

func (s *StructPattern) patternNode() {}

// FieldPattern represents a field pattern in a struct pattern
type FieldPattern struct {
	baseNode
	Name    *Ident
	Pattern Pattern // Can be nil (shorthand)
}

// NewStructPattern creates a new struct pattern
func NewStructPattern(name *Ident, fields []*FieldPattern, span common.Span) *StructPattern {
	return &StructPattern{
		baseNode: baseNode{span},
		Name:     name,
		Fields:   fields,
	}
}

// NewFieldPattern creates a new field pattern
func NewFieldPattern(name *Ident, pattern Pattern, span common.Span) *FieldPattern {
	return &FieldPattern{
		baseNode: baseNode{span},
		Name:     name,
		Pattern:  pattern,
	}
}

// EnumPattern represents an enum variant pattern
type EnumPattern struct {
	baseNode
	Name     *Ident
	Patterns []Pattern // Patterns for variant data
}

func (e *EnumPattern) patternNode() {}

// NewEnumPattern creates a new enum pattern
func NewEnumPattern(name *Ident, patterns []Pattern, span common.Span) *EnumPattern {
	return &EnumPattern{
		baseNode: baseNode{span},
		Name:     name,
		Patterns: patterns,
	}
}

// ==================== Statements ====================

// Stmt represents a statement node
type Stmt interface {
	Node
	stmtNode()
}

// LetStmt represents a variable declaration
type LetStmt struct {
	baseNode
	Name     *Ident
	TypeExpr TypeExpr // Can be nil (type inference)
	Value    Expr
}

func (l *LetStmt) stmtNode() {}

// NewLetStmt creates a new let statement
func NewLetStmt(name *Ident, typeExpr TypeExpr, value Expr, span common.Span) *LetStmt {
	return &LetStmt{
		baseNode: baseNode{span},
		Name:     name,
		TypeExpr: typeExpr,
		Value:    value,
	}
}

// ExprStmt represents an expression statement
type ExprStmt struct {
	baseNode
	Expr Expr
}

func (e *ExprStmt) stmtNode() {}

// NewExprStmt creates a new expression statement
func NewExprStmt(expr Expr, span common.Span) *ExprStmt {
	return &ExprStmt{
		baseNode: baseNode{span},
		Expr:     expr,
	}
}

// ReturnStmt represents a return statement
type ReturnStmt struct {
	baseNode
	Value Expr // Can be nil
}

func (r *ReturnStmt) stmtNode() {}

// NewReturnStmt creates a new return statement
func NewReturnStmt(value Expr, span common.Span) *ReturnStmt {
	return &ReturnStmt{
		baseNode: baseNode{span},
		Value:    value,
	}
}

// AssignStmt represents an assignment statement
type AssignStmt struct {
	baseNode
	Target Expr // Identifier or field access
	Value  Expr
}

func (a *AssignStmt) stmtNode() {}

// NewAssignStmt creates a new assignment statement
func NewAssignStmt(target, value Expr, span common.Span) *AssignStmt {
	return &AssignStmt{
		baseNode: baseNode{span},
		Target:   target,
		Value:    value,
	}
}

// WhileStmt represents a while loop
type WhileStmt struct {
	baseNode
	Cond Expr
	Body *BlockExpr
}

func (w *WhileStmt) stmtNode() {}

// NewWhileStmt creates a new while statement
func NewWhileStmt(cond Expr, body *BlockExpr, span common.Span) *WhileStmt {
	return &WhileStmt{
		baseNode: baseNode{span},
		Cond:     cond,
		Body:     body,
	}
}

// ForStmt represents a for loop
type ForStmt struct {
	baseNode
	Var  *Ident
	Iter Expr
	Body *BlockExpr
}

func (f *ForStmt) stmtNode() {}

// NewForStmt creates a new for statement
func NewForStmt(varName *Ident, iter Expr, body *BlockExpr, span common.Span) *ForStmt {
	return &ForStmt{
		baseNode: baseNode{span},
		Var:      varName,
		Iter:     iter,
		Body:     body,
	}
}

// BreakStmt represents a break statement
type BreakStmt struct {
	baseNode
}

func (b *BreakStmt) stmtNode() {}

// NewBreakStmt creates a new break statement
func NewBreakStmt(span common.Span) *BreakStmt {
	return &BreakStmt{
		baseNode: baseNode{span},
	}
}

// ContinueStmt represents a continue statement
type ContinueStmt struct {
	baseNode
}

func (c *ContinueStmt) stmtNode() {}

// NewContinueStmt creates a new continue statement
func NewContinueStmt(span common.Span) *ContinueStmt {
	return &ContinueStmt{
		baseNode: baseNode{span},
	}
}

// ==================== Declarations ====================

// Decl represents a declaration node
type Decl interface {
	Node
	declNode()
}

// FnDecl represents a function declaration
type FnDecl struct {
	baseNode
	Name       *Ident
	Params     []*Param
	ReturnType TypeExpr // Can be nil
	Body       *BlockExpr
}

func (f *FnDecl) declNode() {}

// Param represents a function parameter
type Param struct {
	baseNode
	Name     *Ident
	TypeExpr TypeExpr
}

// NewFnDecl creates a new function declaration
func NewFnDecl(name *Ident, params []*Param, returnType TypeExpr, body *BlockExpr, span common.Span) *FnDecl {
	return &FnDecl{
		baseNode:   baseNode{span},
		Name:       name,
		Params:     params,
		ReturnType: returnType,
		Body:       body,
	}
}

// NewParam creates a new parameter
func NewParam(name *Ident, typeExpr TypeExpr, span common.Span) *Param {
	return &Param{
		baseNode: baseNode{span},
		Name:     name,
		TypeExpr: typeExpr,
	}
}

// StructDecl represents a struct declaration
type StructDecl struct {
	baseNode
	Name   *Ident
	Fields []*Field
}

func (s *StructDecl) declNode() {}

// Field represents a struct field
type Field struct {
	baseNode
	Name     *Ident
	TypeExpr TypeExpr
}

// NewStructDecl creates a new struct declaration
func NewStructDecl(name *Ident, fields []*Field, span common.Span) *StructDecl {
	return &StructDecl{
		baseNode: baseNode{span},
		Name:     name,
		Fields:   fields,
	}
}

// NewField creates a new field
func NewField(name *Ident, typeExpr TypeExpr, span common.Span) *Field {
	return &Field{
		baseNode: baseNode{span},
		Name:     name,
		TypeExpr: typeExpr,
	}
}

// EnumDecl represents an enum declaration
type EnumDecl struct {
	baseNode
	Name     *Ident
	Variants []*Variant
}

func (e *EnumDecl) declNode() {}

// Variant represents an enum variant
type Variant struct {
	baseNode
	Name  *Ident
	Types []TypeExpr // Associated data types
}

// NewEnumDecl creates a new enum declaration
func NewEnumDecl(name *Ident, variants []*Variant, span common.Span) *EnumDecl {
	return &EnumDecl{
		baseNode: baseNode{span},
		Name:     name,
		Variants: variants,
	}
}

// NewVariant creates a new variant
func NewVariant(name *Ident, types []TypeExpr, span common.Span) *Variant {
	return &Variant{
		baseNode: baseNode{span},
		Name:     name,
		Types:    types,
	}
}

// TypeAlias represents a type alias declaration
type TypeAlias struct {
	baseNode
	Name     *Ident
	TypeExpr TypeExpr
}

func (t *TypeAlias) declNode() {}

// NewTypeAlias creates a new type alias
func NewTypeAlias(name *Ident, typeExpr TypeExpr, span common.Span) *TypeAlias {
	return &TypeAlias{
		baseNode: baseNode{span},
		Name:     name,
		TypeExpr: typeExpr,
	}
}

// ImportDecl represents an import declaration
type ImportDecl struct {
	baseNode
	Path string
}

func (i *ImportDecl) declNode() {}

// NewImportDecl creates a new import declaration
func NewImportDecl(path string, span common.Span) *ImportDecl {
	return &ImportDecl{
		baseNode: baseNode{span},
		Path:     path,
	}
}

// ==================== Type Expressions ====================

// TypeExpr represents a type expression
type TypeExpr interface {
	Node
	typeExprNode()
}

// NamedType represents a named type (identifier)
type NamedType struct {
	baseNode
	Name *Ident
}

func (n *NamedType) typeExprNode() {}

// NewNamedType creates a new named type
func NewNamedType(name *Ident, span common.Span) *NamedType {
	return &NamedType{
		baseNode: baseNode{span},
		Name:     name,
	}
}

// ArrayType represents an array type ([T])
type ArrayType struct {
	baseNode
	ElemType TypeExpr
}

func (a *ArrayType) typeExprNode() {}

// NewArrayType creates a new array type
func NewArrayType(elemType TypeExpr, span common.Span) *ArrayType {
	return &ArrayType{
		baseNode: baseNode{span},
		ElemType: elemType,
	}
}

// MapType represents a map type (Map[K, V])
type MapType struct {
	baseNode
	KeyType   TypeExpr
	ValueType TypeExpr
}

func (m *MapType) typeExprNode() {}

// NewMapType creates a new map type
func NewMapType(keyType, valueType TypeExpr, span common.Span) *MapType {
	return &MapType{
		baseNode:  baseNode{span},
		KeyType:   keyType,
		ValueType: valueType,
	}
}

// FunctionType represents a function type (fn(T1, T2) -> T3)
type FunctionType struct {
	baseNode
	ParamTypes []TypeExpr
	ReturnType TypeExpr
}

func (f *FunctionType) typeExprNode() {}

// NewFunctionType creates a new function type
func NewFunctionType(paramTypes []TypeExpr, returnType TypeExpr, span common.Span) *FunctionType {
	return &FunctionType{
		baseNode:   baseNode{span},
		ParamTypes: paramTypes,
		ReturnType: returnType,
	}
}

// ==================== Module ====================

// Module represents a complete source file
type Module struct {
	baseNode
	Decls []Decl
}

// NewModule creates a new module
func NewModule(decls []Decl, span common.Span) *Module {
	return &Module{
		baseNode: baseNode{span},
		Decls:    decls,
	}
}

