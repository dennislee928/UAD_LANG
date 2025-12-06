package interpreter

import (
	"fmt"

	"github.com/dennislee928/uad-lang/internal/ast"
	"github.com/dennislee928/uad-lang/internal/typer"
)

// ValueKind represents the kind of a runtime value
type ValueKind int

const (
	ValueInt ValueKind = iota
	ValueFloat
	ValueBool
	ValueString
	ValueNil
	ValueArray
	ValueMap
	ValueStruct
	ValueFunction
	ValueBuiltin
)

// Value represents a runtime value
type Value interface {
	Kind() ValueKind
	String() string
	Type() typer.Type
}

// ==================== Primitive Values ====================

// IntValue represents an integer value
type IntValue struct {
	Value int64
}

func (v *IntValue) Kind() ValueKind     { return ValueInt }
func (v *IntValue) String() string      { return fmt.Sprintf("%d", v.Value) }
func (v *IntValue) Type() typer.Type    { return typer.IntType }

// FloatValue represents a float value
type FloatValue struct {
	Value float64
}

func (v *FloatValue) Kind() ValueKind   { return ValueFloat }
func (v *FloatValue) String() string    { return fmt.Sprintf("%g", v.Value) }
func (v *FloatValue) Type() typer.Type  { return typer.FloatType }

// BoolValue represents a boolean value
type BoolValue struct {
	Value bool
}

func (v *BoolValue) Kind() ValueKind    { return ValueBool }
func (v *BoolValue) String() string {
	if v.Value {
		return "true"
	}
	return "false"
}
func (v *BoolValue) Type() typer.Type   { return typer.BoolType }

// StringValue represents a string value
type StringValue struct {
	Value string
}

func (v *StringValue) Kind() ValueKind  { return ValueString }
func (v *StringValue) String() string   { return v.Value } // Return raw value without quotes
func (v *StringValue) Type() typer.Type { return typer.StringType }

// NilValue represents a nil value
type NilValue struct{}

func (v *NilValue) Kind() ValueKind     { return ValueNil }
func (v *NilValue) String() string      { return "nil" }
func (v *NilValue) Type() typer.Type    { return typer.NilType }

// ==================== Compound Values ====================

// ArrayValue represents an array value
type ArrayValue struct {
	Elements []Value
	ElemType typer.Type
}

func (v *ArrayValue) Kind() ValueKind { return ValueArray }
func (v *ArrayValue) String() string {
	s := "["
	for i, elem := range v.Elements {
		if i > 0 {
			s += ", "
		}
		s += elem.String()
	}
	s += "]"
	return s
}
func (v *ArrayValue) Type() typer.Type {
	return typer.NewArrayType(v.ElemType)
}

// MapValue represents a map value
type MapValue struct {
	Entries  map[string]Value // Using string as key for simplicity
	KeyType   typer.Type
	ValueType typer.Type
}

func (v *MapValue) Kind() ValueKind { return ValueMap }
func (v *MapValue) String() string {
	s := "{"
	first := true
	for k, val := range v.Entries {
		if !first {
			s += ", "
		}
		s += fmt.Sprintf("%s: %s", k, val.String())
		first = false
	}
	s += "}"
	return s
}
func (v *MapValue) Type() typer.Type {
	return typer.NewMapType(v.KeyType, v.ValueType)
}

// StructValue represents a struct value
type StructValue struct {
	TypeName string
	Fields   map[string]Value
	StructType *typer.StructType
}

func (v *StructValue) Kind() ValueKind { return ValueStruct }
func (v *StructValue) String() string {
	s := v.TypeName + " {"
	first := true
	for name, val := range v.Fields {
		if !first {
			s += ", "
		}
		s += fmt.Sprintf("%s: %s", name, val.String())
		first = false
	}
	s += "}"
	return s
}
func (v *StructValue) Type() typer.Type {
	return v.StructType
}

// ==================== Function Values ====================

// FunctionValue represents a user-defined function
type FunctionValue struct {
	Name   string
	Params []*ast.Param
	Body   *ast.BlockExpr
	Env    *Environment // Closure environment
	FnType *typer.FunctionType
}

func (v *FunctionValue) Kind() ValueKind { return ValueFunction }
func (v *FunctionValue) String() string {
	return fmt.Sprintf("<function %s>", v.Name)
}
func (v *FunctionValue) Type() typer.Type {
	return v.FnType
}

// BuiltinFunction represents a built-in function
type BuiltinFunction struct {
	Name string
	Fn   func(*Interpreter, []Value) (Value, error)
	FnType *typer.FunctionType
}

func (v *BuiltinFunction) Kind() ValueKind { return ValueBuiltin }
func (v *BuiltinFunction) String() string {
	return fmt.Sprintf("<builtin %s>", v.Name)
}
func (v *BuiltinFunction) Type() typer.Type {
	return v.FnType
}

// ==================== Value Constructors ====================

// NewIntValue creates a new integer value
func NewIntValue(v int64) *IntValue {
	return &IntValue{Value: v}
}

// NewFloatValue creates a new float value
func NewFloatValue(v float64) *FloatValue {
	return &FloatValue{Value: v}
}

// NewBoolValue creates a new boolean value
func NewBoolValue(v bool) *BoolValue {
	return &BoolValue{Value: v}
}

// NewStringValue creates a new string value
func NewStringValue(v string) *StringValue {
	return &StringValue{Value: v}
}

// NewNilValue creates a new nil value
func NewNilValue() *NilValue {
	return &NilValue{}
}

// NewArrayValue creates a new array value
func NewArrayValue(elements []Value, elemType typer.Type) *ArrayValue {
	return &ArrayValue{Elements: elements, ElemType: elemType}
}

// NewMapValue creates a new map value
func NewMapValue(keyType, valueType typer.Type) *MapValue {
	return &MapValue{
		Entries:   make(map[string]Value),
		KeyType:   keyType,
		ValueType: valueType,
	}
}

// NewStructValue creates a new struct value
func NewStructValue(typeName string, fields map[string]Value, structType *typer.StructType) *StructValue {
	return &StructValue{
		TypeName:   typeName,
		Fields:     fields,
		StructType: structType,
	}
}

// NewFunctionValue creates a new function value
func NewFunctionValue(name string, params []*ast.Param, body *ast.BlockExpr, env *Environment, fnType *typer.FunctionType) *FunctionValue {
	return &FunctionValue{
		Name:   name,
		Params: params,
		Body:   body,
		Env:    env,
		FnType: fnType,
	}
}

// NewBuiltinFunction creates a new builtin function
func NewBuiltinFunction(name string, fn func(*Interpreter, []Value) (Value, error), fnType *typer.FunctionType) *BuiltinFunction {
	return &BuiltinFunction{
		Name:   name,
		Fn:     fn,
		FnType: fnType,
	}
}

// ==================== Value Utilities ====================

// IsTruthy checks if a value is truthy
func IsTruthy(v Value) bool {
	switch val := v.(type) {
	case *BoolValue:
		return val.Value
	case *NilValue:
		return false
	case *IntValue:
		return val.Value != 0
	case *FloatValue:
		return val.Value != 0.0
	case *StringValue:
		return val.Value != ""
	default:
		return true
	}
}

// IsEqual checks if two values are equal
func IsEqual(a, b Value) bool {
	if a.Kind() != b.Kind() {
		return false
	}
	
	switch v1 := a.(type) {
	case *IntValue:
		v2 := b.(*IntValue)
		return v1.Value == v2.Value
	case *FloatValue:
		v2 := b.(*FloatValue)
		return v1.Value == v2.Value
	case *BoolValue:
		v2 := b.(*BoolValue)
		return v1.Value == v2.Value
	case *StringValue:
		v2 := b.(*StringValue)
		return v1.Value == v2.Value
	case *NilValue:
		return true
	default:
		return false // Complex types need more sophisticated comparison
	}
}

// ToInt converts a value to an integer
func ToInt(v Value) (int64, error) {
	switch val := v.(type) {
	case *IntValue:
		return val.Value, nil
	case *FloatValue:
		return int64(val.Value), nil
	default:
		return 0, fmt.Errorf("cannot convert %s to Int", v.Type().String())
	}
}

// ToFloat converts a value to a float
func ToFloat(v Value) (float64, error) {
	switch val := v.(type) {
	case *IntValue:
		return float64(val.Value), nil
	case *FloatValue:
		return val.Value, nil
	default:
		return 0, fmt.Errorf("cannot convert %s to Float", v.Type().String())
	}
}

// ToString converts a value to a string
func ToString(v Value) string {
	return v.String()
}

