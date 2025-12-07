package ir

import (
	"fmt"
	"strconv"

	"github.com/dennislee928/uad-lang/internal/ast"
	"github.com/dennislee928/uad-lang/internal/typer"
)

// Builder builds IR from AST
type Builder struct {
	module      *Module
	currentFunc *Function
	typeChecker *typer.TypeChecker
	
	// Variable tracking
	locals      map[string]int32  // name -> local index
	localCount  int32
	
	// Loop tracking (for break/continue)
	loopStack   []*loopContext
}

type loopContext struct {
	breakLabel    *Label
	continueLabel *Label
}

// NewBuilder creates a new IR builder
func NewBuilder(typeChecker *typer.TypeChecker) *Builder {
	return &Builder{
		module:      NewModule(),
		typeChecker: typeChecker,
		locals:      make(map[string]int32),
	}
}

// Build builds IR from an AST module
func (b *Builder) Build(astModule *ast.Module) (*Module, error) {
	// First pass: collect all function declarations
	for _, decl := range astModule.Decls {
		if fnDecl, ok := decl.(*ast.FnDecl); ok {
			fn := NewFunction(fnDecl.Name.Name, fnDecl.Span())
			b.module.AddFunction(fn)
		}
	}
	
	// Second pass: build function bodies
	for _, decl := range astModule.Decls {
		if err := b.buildDecl(decl); err != nil {
			return nil, err
		}
	}
	
	return b.module, nil
}

// ==================== Declarations ====================

func (b *Builder) buildDecl(decl ast.Decl) error {
	switch d := decl.(type) {
	case *ast.FnDecl:
		return b.buildFnDecl(d)
	case *ast.StructDecl, *ast.EnumDecl, *ast.TypeAlias:
		// Type declarations don't generate IR
		return nil
	default:
		return fmt.Errorf("unknown declaration type: %T", d)
	}
}

func (b *Builder) buildFnDecl(decl *ast.FnDecl) error {
	fn, ok := b.module.GetFunction(decl.Name.Name)
	if !ok {
		return fmt.Errorf("function '%s' not found", decl.Name.Name)
	}
	
	// Set current function
	b.currentFunc = fn
	b.locals = make(map[string]int32)
	b.localCount = 0
	
	// Get function type from type checker
	fnSym, ok := b.typeChecker.Env().Lookup(decl.Name.Name)
	if !ok {
		return fmt.Errorf("function '%s' not in type environment", decl.Name.Name)
	}
	
	fnType, ok := fnSym.Type.(*typer.FunctionType)
	if !ok {
		return fmt.Errorf("'%s' is not a function", decl.Name.Name)
	}
	
	fn.ReturnType = fnType.ReturnType
	
	// Add parameters
	for i, param := range decl.Params {
		paramType := fnType.ParamTypes[i]
		local := NewLocal(param.Name.Name, paramType, b.localCount)
		b.localCount++
		fn.AddParam(local)
		b.locals[param.Name.Name] = local.Index
	}
	
	// Build function body
	if err := b.buildBlockExpr(decl.Body); err != nil {
		return err
	}
	
	// Add implicit return if needed
	if len(fn.Instructions) == 0 || fn.Instructions[len(fn.Instructions)-1].Op != OpReturn {
		fn.AddInstruction(NewInstruction(OpConstNil, 0, decl.Span()))
		fn.AddInstruction(NewInstruction(OpReturn, 0, decl.Span()))
	}
	
	return nil
}

// ==================== Statements ====================

func (b *Builder) buildStmt(stmt ast.Stmt) error {
	switch s := stmt.(type) {
	case *ast.LetStmt:
		return b.buildLetStmt(s)
	case *ast.ExprStmt:
		if err := b.buildExpr(s.Expr); err != nil {
			return err
		}
		// Pop the expression result since it's not used
		b.currentFunc.AddInstruction(NewInstruction(OpPop, 0, s.Span()))
		return nil
	case *ast.ReturnStmt:
		return b.buildReturnStmt(s)
	case *ast.AssignStmt:
		return b.buildAssignStmt(s)
	case *ast.WhileStmt:
		return b.buildWhileStmt(s)
	case *ast.ForStmt:
		return b.buildForStmt(s)
	case *ast.BreakStmt:
		return b.buildBreakStmt(s)
	case *ast.ContinueStmt:
		return b.buildContinueStmt(s)
	default:
		return fmt.Errorf("unknown statement type: %T", s)
	}
}

func (b *Builder) buildLetStmt(stmt *ast.LetStmt) error {
	// Build value expression
	if err := b.buildExpr(stmt.Value); err != nil {
		return err
	}
	
	// Get type from type checker
	exprType := b.typeChecker.GetExprType(stmt.Value)
	if exprType == nil {
		exprType = typer.UnitType
	}
	
	// Allocate local variable
	local := NewLocal(stmt.Name.Name, exprType, b.localCount)
	b.localCount++
	b.currentFunc.AddLocal(local)
	b.locals[stmt.Name.Name] = local.Index
	
	// Store value in local
	b.currentFunc.AddInstruction(NewInstruction(OpSetLocal, local.Index, stmt.Span()))
	
	return nil
}

func (b *Builder) buildReturnStmt(stmt *ast.ReturnStmt) error {
	if stmt.Value == nil {
		// Return nil
		b.currentFunc.AddInstruction(NewInstruction(OpConstNil, 0, stmt.Span()))
	} else {
		if err := b.buildExpr(stmt.Value); err != nil {
			return err
		}
	}
	
	b.currentFunc.AddInstruction(NewInstruction(OpReturn, 0, stmt.Span()))
	return nil
}

func (b *Builder) buildAssignStmt(stmt *ast.AssignStmt) error {
	// Build value expression
	if err := b.buildExpr(stmt.Value); err != nil {
		return err
	}
	
	// Handle different assignment targets
	switch target := stmt.Target.(type) {
	case *ast.Ident:
		// Simple variable assignment
		idx, ok := b.locals[target.Name]
		if !ok {
			return fmt.Errorf("undefined variable '%s'", target.Name)
		}
		b.currentFunc.AddInstruction(NewInstruction(OpSetLocal, idx, stmt.Span()))
		return nil
		
	case *ast.FieldAccess:
		// Build target expression
		if err := b.buildExpr(target.Expr); err != nil {
			return err
		}
		// Build value (already on stack)
		// Field name as constant
		constIdx := b.currentFunc.AddConstant(NewConstant(ValueString, target.Field.Name))
		b.currentFunc.AddInstruction(NewInstruction(OpSetField, constIdx, stmt.Span()))
		return nil
		
	case *ast.IndexExpr:
		// Build target expression
		if err := b.buildExpr(target.Expr); err != nil {
			return err
		}
		// Build index expression
		if err := b.buildExpr(target.Index); err != nil {
			return err
		}
		// Value is already on stack from earlier
		b.currentFunc.AddInstruction(NewInstruction(OpSetIndex, 0, stmt.Span()))
		return nil
		
	default:
		return fmt.Errorf("invalid assignment target: %T", target)
	}
}

func (b *Builder) buildWhileStmt(stmt *ast.WhileStmt) error {
	// Labels
	loopStart := NewLabel("while_start")
	loopEnd := NewLabel("while_end")
	
	// Push loop context
	b.loopStack = append(b.loopStack, &loopContext{
		breakLabel:    loopEnd,
		continueLabel: loopStart,
	})
	defer func() { b.loopStack = b.loopStack[:len(b.loopStack)-1] }()
	
	// Loop start
	startAddr := int32(len(b.currentFunc.Instructions))
	loopStart.Resolve(startAddr)
	
	// Build condition
	if err := b.buildExpr(stmt.Cond); err != nil {
		return err
	}
	
	// Jump to end if false
	jumpInstr := NewInstruction(OpJumpIfFalse, 0, stmt.Cond.Span())
	jumpIdx := b.currentFunc.AddInstruction(jumpInstr)
	
	// Build body
	if err := b.buildBlockExpr(stmt.Body); err != nil {
		return err
	}
	
	// Jump back to start
	b.currentFunc.AddInstruction(NewInstruction(OpJump, startAddr, stmt.Span()))
	
	// Resolve end label
	endAddr := int32(len(b.currentFunc.Instructions))
	loopEnd.Resolve(endAddr)
	
	// Patch jump instruction
	b.currentFunc.Instructions[jumpIdx].Operand = endAddr
	
	return nil
}

func (b *Builder) buildForStmt(stmt *ast.ForStmt) error {
	// Build iterator expression
	if err := b.buildExpr(stmt.Iter); err != nil {
		return err
	}
	
	// TODO: Implement for loop properly
	// For now, just pop the iterator value
	b.currentFunc.AddInstruction(NewInstruction(OpPop, 0, stmt.Span()))
	
	return fmt.Errorf("for loops not yet fully implemented in IR")
}

func (b *Builder) buildBreakStmt(stmt *ast.BreakStmt) error {
	if len(b.loopStack) == 0 {
		return fmt.Errorf("break outside of loop")
	}
	
	loop := b.loopStack[len(b.loopStack)-1]
	if loop.breakLabel.IsResolved() {
		b.currentFunc.AddInstruction(NewInstruction(OpJump, loop.breakLabel.Address, stmt.Span()))
	} else {
		// Will be patched later
		b.currentFunc.AddInstruction(NewInstruction(OpJump, -1, stmt.Span()))
	}
	
	return nil
}

func (b *Builder) buildContinueStmt(stmt *ast.ContinueStmt) error {
	if len(b.loopStack) == 0 {
		return fmt.Errorf("continue outside of loop")
	}
	
	loop := b.loopStack[len(b.loopStack)-1]
	if loop.continueLabel.IsResolved() {
		b.currentFunc.AddInstruction(NewInstruction(OpJump, loop.continueLabel.Address, stmt.Span()))
	} else {
		// Will be patched later
		b.currentFunc.AddInstruction(NewInstruction(OpJump, -1, stmt.Span()))
	}
	
	return nil
}

// ==================== Expressions ====================

func (b *Builder) buildExpr(expr ast.Expr) error {
	switch e := expr.(type) {
	case *ast.Ident:
		return b.buildIdent(e)
	case *ast.Literal:
		return b.buildLiteral(e)
	case *ast.BinaryExpr:
		return b.buildBinaryExpr(e)
	case *ast.UnaryExpr:
		return b.buildUnaryExpr(e)
	case *ast.CallExpr:
		return b.buildCallExpr(e)
	case *ast.IfExpr:
		return b.buildIfExpr(e)
	case *ast.BlockExpr:
		return b.buildBlockExpr(e)
	case *ast.StructLiteral:
		return b.buildStructLiteral(e)
	case *ast.ArrayLiteral:
		return b.buildArrayLiteral(e)
	case *ast.FieldAccess:
		return b.buildFieldAccess(e)
	case *ast.IndexExpr:
		return b.buildIndexExpr(e)
	case *ast.ParenExpr:
		return b.buildExpr(e.Expr)
	default:
		return fmt.Errorf("unknown expression type: %T", e)
	}
}

func (b *Builder) buildIdent(expr *ast.Ident) error {
	idx, ok := b.locals[expr.Name]
	if !ok {
		// Try global or builtin
		return fmt.Errorf("undefined variable '%s'", expr.Name)
	}
	
	b.currentFunc.AddInstruction(NewInstruction(OpGetLocal, idx, expr.Span()))
	return nil
}

func (b *Builder) buildLiteral(expr *ast.Literal) error {
	switch expr.Kind {
	case ast.LitInt:
		val, err := strconv.ParseInt(expr.Value, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid integer literal: %s", expr.Value)
		}
		constIdx := b.currentFunc.AddConstant(NewConstant(ValueInt, val))
		b.currentFunc.AddInstruction(NewInstruction(OpConstInt, constIdx, expr.Span()))
		
	case ast.LitFloat:
		val, err := strconv.ParseFloat(expr.Value, 64)
		if err != nil {
			return fmt.Errorf("invalid float literal: %s", expr.Value)
		}
		constIdx := b.currentFunc.AddConstant(NewConstant(ValueFloat, val))
		b.currentFunc.AddInstruction(NewInstruction(OpConstFloat, constIdx, expr.Span()))
		
	case ast.LitString:
		// Unquote the string
		val, err := strconv.Unquote(expr.Value)
		if err != nil {
			val = expr.Value
		}
		constIdx := b.currentFunc.AddConstant(NewConstant(ValueString, val))
		b.currentFunc.AddInstruction(NewInstruction(OpConstString, constIdx, expr.Span()))
		
	case ast.LitBool:
		val := expr.Value == "true"
		constIdx := b.currentFunc.AddConstant(NewConstant(ValueBool, val))
		b.currentFunc.AddInstruction(NewInstruction(OpConstBool, constIdx, expr.Span()))
		
	case ast.LitNil:
		b.currentFunc.AddInstruction(NewInstruction(OpConstNil, 0, expr.Span()))
		
	default:
		return fmt.Errorf("unknown literal kind: %v", expr.Kind)
	}
	
	return nil
}

func (b *Builder) buildBinaryExpr(expr *ast.BinaryExpr) error {
	// Build left operand
	if err := b.buildExpr(expr.Left); err != nil {
		return err
	}
	
	// Build right operand
	if err := b.buildExpr(expr.Right); err != nil {
		return err
	}
	
	// Emit operation
	var op OpCode
	switch expr.Op {
	case ast.OpAdd:
		op = OpAdd
	case ast.OpSub:
		op = OpSub
	case ast.OpMul:
		op = OpMul
	case ast.OpDiv:
		op = OpDiv
	case ast.OpMod:
		op = OpMod
	case ast.OpEq:
		op = OpEq
	case ast.OpNeq:
		op = OpNeq
	case ast.OpLt:
		op = OpLt
	case ast.OpGt:
		op = OpGt
	case ast.OpLe:
		op = OpLe
	case ast.OpGe:
		op = OpGe
	case ast.OpAnd:
		op = OpAnd
	case ast.OpOr:
		op = OpOr
	default:
		return fmt.Errorf("unknown binary operator: %v", expr.Op)
	}
	
	b.currentFunc.AddInstruction(NewInstruction(op, 0, expr.Span()))
	return nil
}

func (b *Builder) buildUnaryExpr(expr *ast.UnaryExpr) error {
	// Build operand
	if err := b.buildExpr(expr.Expr); err != nil {
		return err
	}
	
	// Emit operation
	var op OpCode
	switch expr.Op {
	case ast.OpNeg:
		op = OpNeg
	case ast.OpNot:
		op = OpNot
	default:
		return fmt.Errorf("unknown unary operator: %v", expr.Op)
	}
	
	b.currentFunc.AddInstruction(NewInstruction(op, 0, expr.Span()))
	return nil
}

func (b *Builder) buildCallExpr(expr *ast.CallExpr) error {
	// Build arguments (pushed in order)
	for _, arg := range expr.Args {
		if err := b.buildExpr(arg); err != nil {
			return err
		}
	}
	
	// Build function expression
	if err := b.buildExpr(expr.Func); err != nil {
		return err
	}
	
	// Emit call instruction
	b.currentFunc.AddInstruction(NewInstruction(OpCall, int32(len(expr.Args)), expr.Span()))
	return nil
}

func (b *Builder) buildIfExpr(expr *ast.IfExpr) error {
	// Build condition
	if err := b.buildExpr(expr.Cond); err != nil {
		return err
	}
	
	// Jump to else if false
	jumpToElse := NewInstruction(OpJumpIfFalse, 0, expr.Cond.Span())
	jumpToElseIdx := b.currentFunc.AddInstruction(jumpToElse)
	
	// Build then branch
	if err := b.buildBlockExpr(expr.Then); err != nil {
		return err
	}
	
	// Jump to end
	jumpToEnd := NewInstruction(OpJump, 0, expr.Span())
	jumpToEndIdx := b.currentFunc.AddInstruction(jumpToEnd)
	
	// Patch jump to else
	elseAddr := int32(len(b.currentFunc.Instructions))
	b.currentFunc.Instructions[jumpToElseIdx].Operand = elseAddr
	
	// Build else branch
	if expr.Else != nil {
		if elseIf, ok := expr.Else.(*ast.IfExpr); ok {
			if err := b.buildIfExpr(elseIf); err != nil {
				return err
			}
		} else if elseBlock, ok := expr.Else.(*ast.BlockExpr); ok {
			if err := b.buildBlockExpr(elseBlock); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("invalid else clause: %T", expr.Else)
		}
	} else {
		// No else clause, push nil
		b.currentFunc.AddInstruction(NewInstruction(OpConstNil, 0, expr.Span()))
	}
	
	// Patch jump to end
	endAddr := int32(len(b.currentFunc.Instructions))
	b.currentFunc.Instructions[jumpToEndIdx].Operand = endAddr
	
	return nil
}

func (b *Builder) buildBlockExpr(expr *ast.BlockExpr) error {
	// Build statements
	for _, stmt := range expr.Stmts {
		if err := b.buildStmt(stmt); err != nil {
			return err
		}
	}
	
	// Build final expression
	if expr.Expr != nil {
		if err := b.buildExpr(expr.Expr); err != nil {
			return err
		}
	} else {
		// Block without expression returns nil
		b.currentFunc.AddInstruction(NewInstruction(OpConstNil, 0, expr.Span()))
	}
	
	return nil
}

func (b *Builder) buildStructLiteral(expr *ast.StructLiteral) error {
	// TODO: Implement struct literals
	return fmt.Errorf("struct literals not yet implemented in IR")
}

func (b *Builder) buildArrayLiteral(expr *ast.ArrayLiteral) error {
	// Build all elements
	for _, elem := range expr.Elements {
		if err := b.buildExpr(elem); err != nil {
			return err
		}
	}
	
	// Create array
	b.currentFunc.AddInstruction(NewInstruction(OpNewArray, int32(len(expr.Elements)), expr.Span()))
	return nil
}

func (b *Builder) buildFieldAccess(expr *ast.FieldAccess) error {
	// Build target expression
	if err := b.buildExpr(expr.Expr); err != nil {
		return err
	}
	
	// Get field index (stored as constant)
	constIdx := b.currentFunc.AddConstant(NewConstant(ValueString, expr.Field.Name))
	b.currentFunc.AddInstruction(NewInstruction(OpGetField, constIdx, expr.Span()))
	
	return nil
}

func (b *Builder) buildIndexExpr(expr *ast.IndexExpr) error {
	// Build target expression
	if err := b.buildExpr(expr.Expr); err != nil {
		return err
	}
	
	// Build index expression
	if err := b.buildExpr(expr.Index); err != nil {
		return err
	}
	
	// Get element
	b.currentFunc.AddInstruction(NewInstruction(OpGetIndex, 0, expr.Span()))
	return nil
}

