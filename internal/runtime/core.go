package runtime

import "fmt"

// core.go provides core runtime concepts and execution logic.
// This module defines the fundamental building blocks for UAD program execution,
// including value representation, environment management, and execution context.

// Value represents a runtime value in UAD.
// All values during execution are represented by this type.
type Value interface {
	Type() ValueType
	String() string
}

// ValueType represents the type of a runtime value.
type ValueType int

const (
	ValueTypeNil ValueType = iota
	ValueTypeInt
	ValueTypeFloat
	ValueTypeBool
	ValueTypeString
	ValueTypeArray
	ValueTypeMap
	ValueTypeStruct
	ValueTypeFunction
	ValueTypeDuration
	ValueTypeTime
)

// ==================== Concrete Value Types ====================

// NilValue represents the nil value.
type NilValue struct{}

func (v *NilValue) Type() ValueType { return ValueTypeNil }
func (v *NilValue) String() string  { return "nil" }

// IntValue represents an integer value.
type IntValue struct {
	Value int64
}

func (v *IntValue) Type() ValueType { return ValueTypeInt }
func (v *IntValue) String() string  { return fmt.Sprintf("%d", v.Value) }

// FloatValue represents a floating-point value.
type FloatValue struct {
	Value float64
}

func (v *FloatValue) Type() ValueType { return ValueTypeFloat }
func (v *FloatValue) String() string  { return fmt.Sprintf("%f", v.Value) }

// BoolValue represents a boolean value.
type BoolValue struct {
	Value bool
}

func (v *BoolValue) Type() ValueType { return ValueTypeBool }
func (v *BoolValue) String() string {
	if v.Value {
		return "true"
	}
	return "false"
}

// StringValue represents a string value.
type StringValue struct {
	Value string
}

func (v *StringValue) Type() ValueType { return ValueTypeString }
func (v *StringValue) String() string  { return v.Value }

// ArrayValue represents an array value.
type ArrayValue struct {
	Elements []Value
}

func (v *ArrayValue) Type() ValueType { return ValueTypeArray }
func (v *ArrayValue) String() string {
	if len(v.Elements) == 0 {
		return "[]"
	}
	result := "["
	for i, elem := range v.Elements {
		if i > 0 {
			result += ", "
		}
		result += elem.String()
	}
	result += "]"
	return result
}

// MapValue represents a map value.
type MapValue struct {
	Entries map[string]Value
}

func (v *MapValue) Type() ValueType { return ValueTypeMap }
func (v *MapValue) String() string {
	if len(v.Entries) == 0 {
		return "{}"
	}
	result := "{"
	first := true
	for key, val := range v.Entries {
		if !first {
			result += ", "
		}
		result += fmt.Sprintf("%s: %s", key, val.String())
		first = false
	}
	result += "}"
	return result
}

// StructValue represents a struct instance.
type StructValue struct {
	TypeName string
	Fields   map[string]Value
}

func (v *StructValue) Type() ValueType { return ValueTypeStruct }
func (v *StructValue) String() string {
	return fmt.Sprintf("%s {...}", v.TypeName)
}

// FunctionValue represents a function value.
type FunctionValue struct {
	Name       string
	Parameters []string
	Body       interface{} // AST node or native function
	Closure    *Environment
}

func (v *FunctionValue) Type() ValueType { return ValueTypeFunction }
func (v *FunctionValue) String() string {
	return fmt.Sprintf("fn %s(...)", v.Name)
}

// DurationValue represents a time duration.
type DurationValue struct {
	Seconds int64
}

func (v *DurationValue) Type() ValueType { return ValueTypeDuration }
func (v *DurationValue) String() string {
	// Convert seconds to human-readable format
	s := v.Seconds
	if s == 0 {
		return "0s"
	}
	
	days := s / 86400
	s %= 86400
	hours := s / 3600
	s %= 3600
	minutes := s / 60
	seconds := s % 60
	
	result := ""
	if days > 0 {
		result += fmt.Sprintf("%dd", days)
	}
	if hours > 0 {
		result += fmt.Sprintf("%dh", hours)
	}
	if minutes > 0 {
		result += fmt.Sprintf("%dm", minutes)
	}
	if seconds > 0 {
		result += fmt.Sprintf("%ds", seconds)
	}
	
	return result
}

// TimeValue represents an absolute time point.
type TimeValue struct {
	UnixTimestamp int64
}

func (v *TimeValue) Type() ValueType { return ValueTypeTime }
func (v *TimeValue) String() string {
	// Return ISO 8601 format (basic implementation)
	// For full implementation, would use time.Unix(v.UnixTimestamp, 0).Format(time.RFC3339)
	// But avoiding time package dependency in runtime/core for now
	return fmt.Sprintf("@%d", v.UnixTimestamp)
}

// ==================== Environment (Variable Bindings) ====================

// Environment represents a lexical environment for variable bindings.
// Environments form a chain (linked list) to support nested scopes.
type Environment struct {
	bindings map[string]Value
	parent   *Environment
}

// NewEnvironment creates a new environment with an optional parent.
func NewEnvironment(parent *Environment) *Environment {
	return &Environment{
		bindings: make(map[string]Value),
		parent:   parent,
	}
}

// Define adds a new binding to the current environment.
// Returns an error if the variable is already defined in this scope.
func (e *Environment) Define(name string, value Value) error {
	if _, exists := e.bindings[name]; exists {
		return fmt.Errorf("variable '%s' already defined in this scope", name)
	}
	e.bindings[name] = value
	return nil
}

// Set updates an existing binding in the environment chain.
// Returns an error if the variable is not found.
func (e *Environment) Set(name string, value Value) error {
	if _, exists := e.bindings[name]; exists {
		e.bindings[name] = value
		return nil
	}
	if e.parent != nil {
		return e.parent.Set(name, value)
	}
	return fmt.Errorf("undefined variable '%s'", name)
}

// Get retrieves a value from the environment chain.
// Returns nil if the variable is not found.
func (e *Environment) Get(name string) (Value, bool) {
	if value, exists := e.bindings[name]; exists {
		return value, true
	}
	if e.parent != nil {
		return e.parent.Get(name)
	}
	return nil, false
}

// ==================== Execution Context ====================

// ExecutionContext holds the state during program execution.
// This includes the current environment, call stack, and output buffer.
type ExecutionContext struct {
	GlobalEnv *Environment
	CurrentEnv *Environment
	CallStack []string // Function call stack for debugging
	OutputBuffer string // Accumulated output from print statements
}

// NewExecutionContext creates a new execution context.
func NewExecutionContext() *ExecutionContext {
	globalEnv := NewEnvironment(nil)
	return &ExecutionContext{
		GlobalEnv:    globalEnv,
		CurrentEnv:   globalEnv,
		CallStack:    []string{},
		OutputBuffer: "",
	}
}

// PushScope creates a new nested scope.
func (ctx *ExecutionContext) PushScope() {
	ctx.CurrentEnv = NewEnvironment(ctx.CurrentEnv)
}

// PopScope returns to the parent scope.
func (ctx *ExecutionContext) PopScope() {
	if ctx.CurrentEnv.parent != nil {
		ctx.CurrentEnv = ctx.CurrentEnv.parent
	}
}

// PushCall adds a function call to the call stack.
func (ctx *ExecutionContext) PushCall(funcName string) {
	ctx.CallStack = append(ctx.CallStack, funcName)
}

// PopCall removes the top function call from the call stack.
func (ctx *ExecutionContext) PopCall() {
	if len(ctx.CallStack) > 0 {
		ctx.CallStack = ctx.CallStack[:len(ctx.CallStack)-1]
	}
}

// AppendOutput adds text to the output buffer.
func (ctx *ExecutionContext) AppendOutput(text string) {
	ctx.OutputBuffer += text
}

// GetOutput returns the accumulated output.
func (ctx *ExecutionContext) GetOutput() string {
	return ctx.OutputBuffer
}

// ClearOutput clears the output buffer.
func (ctx *ExecutionContext) ClearOutput() {
	ctx.OutputBuffer = ""
}

// ==================== Helper Functions ====================

// IsTruthy determines if a value is considered "true" in a boolean context.
// In UAD: false and nil are falsy, everything else is truthy.
func IsTruthy(v Value) bool {
	if v == nil {
		return false
	}
	switch val := v.(type) {
	case *NilValue:
		return false
	case *BoolValue:
		return val.Value
	default:
		return true
	}
}

// ValuesEqual checks if two values are equal.
func ValuesEqual(a, b Value) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if a.Type() != b.Type() {
		return false
	}

	switch av := a.(type) {
	case *IntValue:
		return av.Value == b.(*IntValue).Value
	case *FloatValue:
		return av.Value == b.(*FloatValue).Value
	case *BoolValue:
		return av.Value == b.(*BoolValue).Value
	case *StringValue:
		return av.Value == b.(*StringValue).Value
	case *NilValue:
		return true
	default:
		// For complex types, use pointer equality for now
		return a == b
	}
}

