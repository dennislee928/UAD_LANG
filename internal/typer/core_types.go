package typer

import (
	"fmt"
	"strings"
)

// Type represents a type in the .uad type system
type Type interface {
	String() string
	Equals(other Type) bool
	typeNode() // private marker method
}

// ==================== Primitive Types ====================

// PrimitiveType represents a primitive type
type PrimitiveType struct {
	Name string
}

func (p *PrimitiveType) String() string       { return p.Name }
func (p *PrimitiveType) Equals(other Type) bool {
	if o, ok := other.(*PrimitiveType); ok {
		return p.Name == o.Name
	}
	return false
}
func (p *PrimitiveType) typeNode() {}

// Predefined primitive types
var (
	IntType      = &PrimitiveType{Name: "Int"}
	FloatType    = &PrimitiveType{Name: "Float"}
	BoolType     = &PrimitiveType{Name: "Bool"}
	StringType   = &PrimitiveType{Name: "String"}
	DurationType = &PrimitiveType{Name: "Duration"}
	TimeType     = &PrimitiveType{Name: "Time"}
	NilType      = &PrimitiveType{Name: "Nil"}
	UnitType     = &PrimitiveType{Name: "Unit"} // for statements without value
	NeverType    = &PrimitiveType{Name: "Never"} // for expressions that never return
)

// ==================== Compound Types ====================

// ArrayType represents an array type [T]
type ArrayType struct {
	ElemType Type
}

func (a *ArrayType) String() string {
	return fmt.Sprintf("[%s]", a.ElemType.String())
}

func (a *ArrayType) Equals(other Type) bool {
	if o, ok := other.(*ArrayType); ok {
		return a.ElemType.Equals(o.ElemType)
	}
	return false
}

func (a *ArrayType) typeNode() {}

// MapType represents a map type Map[K, V]
type MapType struct {
	KeyType   Type
	ValueType Type
}

func (m *MapType) String() string {
	return fmt.Sprintf("Map[%s, %s]", m.KeyType.String(), m.ValueType.String())
}

func (m *MapType) Equals(other Type) bool {
	if o, ok := other.(*MapType); ok {
		return m.KeyType.Equals(o.KeyType) && m.ValueType.Equals(o.ValueType)
	}
	return false
}

func (m *MapType) typeNode() {}

// FunctionType represents a function type fn(T1, T2, ...) -> R
type FunctionType struct {
	ParamTypes []Type
	ReturnType Type
}

func (f *FunctionType) String() string {
	params := make([]string, len(f.ParamTypes))
	for i, p := range f.ParamTypes {
		params[i] = p.String()
	}
	return fmt.Sprintf("fn(%s) -> %s", strings.Join(params, ", "), f.ReturnType.String())
}

func (f *FunctionType) Equals(other Type) bool {
	if o, ok := other.(*FunctionType); ok {
		if len(f.ParamTypes) != len(o.ParamTypes) {
			return false
		}
		for i := range f.ParamTypes {
			if !f.ParamTypes[i].Equals(o.ParamTypes[i]) {
				return false
			}
		}
		return f.ReturnType.Equals(o.ReturnType)
	}
	return false
}

func (f *FunctionType) typeNode() {}

// ==================== User-Defined Types ====================

// StructType represents a struct type
type StructType struct {
	Name   string
	Fields map[string]Type // field name -> type
}

func (s *StructType) String() string {
	return s.Name
}

func (s *StructType) Equals(other Type) bool {
	if o, ok := other.(*StructType); ok {
		return s.Name == o.Name // Nominal typing for structs
	}
	return false
}

func (s *StructType) typeNode() {}

// GetField returns the type of a field, or nil if not found
func (s *StructType) GetField(name string) Type {
	return s.Fields[name]
}

// EnumType represents an enum type
type EnumType struct {
	Name     string
	Variants map[string][]Type // variant name -> associated types
}

func (e *EnumType) String() string {
	return e.Name
}

func (e *EnumType) Equals(other Type) bool {
	if o, ok := other.(*EnumType); ok {
		return e.Name == o.Name // Nominal typing for enums
	}
	return false
}

func (e *EnumType) typeNode() {}

// GetVariant returns the associated types of a variant, or nil if not found
func (e *EnumType) GetVariant(name string) []Type {
	return e.Variants[name]
}

// ==================== Domain Types (for ERH) ====================

// ActionType represents the Action type
type ActionType struct {
	ComplexityType Type
	TagsType       Type
}

func (a *ActionType) String() string {
	return "Action"
}

func (a *ActionType) Equals(other Type) bool {
	_, ok := other.(*ActionType)
	return ok
}

func (a *ActionType) typeNode() {}

// JudgeType represents the Judge type
type JudgeType struct {
	InputType  Type
	OutputType Type
}

func (j *JudgeType) String() string {
	return "Judge"
}

func (j *JudgeType) Equals(other Type) bool {
	_, ok := other.(*JudgeType)
	return ok
}

func (j *JudgeType) typeNode() {}

// AgentType represents the Agent type (for cognitive systems)
type AgentType struct {
	StateType Type
}

func (a *AgentType) String() string {
	return "Agent"
}

func (a *AgentType) Equals(other Type) bool {
	_, ok := other.(*AgentType)
	return ok
}

func (a *AgentType) typeNode() {}

// ==================== Type Variables (for generics & inference) ====================

// TypeVar represents a type variable (for type inference)
type TypeVar struct {
	Name string
	ID   int // Unique identifier
}

func (t *TypeVar) String() string {
	return fmt.Sprintf("'%s_%d", t.Name, t.ID)
}

func (t *TypeVar) Equals(other Type) bool {
	if o, ok := other.(*TypeVar); ok {
		return t.ID == o.ID
	}
	return false
}

func (t *TypeVar) typeNode() {}

// ==================== Type Aliases ====================

// AliasType represents a type alias
type AliasType struct {
	Name       string
	AliasedType Type
}

func (a *AliasType) String() string {
	return a.Name
}

func (a *AliasType) Equals(other Type) bool {
	// Aliases are transparent for type checking
	return a.AliasedType.Equals(other)
}

func (a *AliasType) typeNode() {}

// Resolve returns the underlying type
func (a *AliasType) Resolve() Type {
	if alias, ok := a.AliasedType.(*AliasType); ok {
		return alias.Resolve()
	}
	return a.AliasedType
}

// ==================== Type Utilities ====================

// IsNumeric checks if a type is numeric (Int or Float)
func IsNumeric(t Type) bool {
	return t.Equals(IntType) || t.Equals(FloatType)
}

// IsComparable checks if a type supports comparison operators
func IsComparable(t Type) bool {
	return IsNumeric(t) || t.Equals(StringType) || t.Equals(BoolType) || t.Equals(DurationType)
}

// IsEquatable checks if a type supports == and !=
func IsEquatable(t Type) bool {
	// Most types are equatable except functions
	_, isFunc := t.(*FunctionType)
	return !isFunc
}

// CanCoerce checks if we can implicitly coerce from -> to
func CanCoerce(from, to Type) bool {
	// Int can be coerced to Float
	if from.Equals(IntType) && to.Equals(FloatType) {
		return true
	}
	
	// Nil can be coerced to any optional type (future feature)
	if from.Equals(NilType) {
		// For now, allow nil to be assigned to maps and arrays
		_, isArray := to.(*ArrayType)
		_, isMap := to.(*MapType)
		return isArray || isMap
	}
	
	return false
}

// Unify attempts to unify two types, resolving type variables
func Unify(t1, t2 Type, subst map[int]Type) bool {
	// Resolve aliases
	if alias, ok := t1.(*AliasType); ok {
		t1 = alias.Resolve()
	}
	if alias, ok := t2.(*AliasType); ok {
		t2 = alias.Resolve()
	}
	
	// Check type variables
	if tv1, ok := t1.(*TypeVar); ok {
		if existing, found := subst[tv1.ID]; found {
			return Unify(existing, t2, subst)
		}
		subst[tv1.ID] = t2
		return true
	}
	
	if tv2, ok := t2.(*TypeVar); ok {
		if existing, found := subst[tv2.ID]; found {
			return Unify(t1, existing, subst)
		}
		subst[tv2.ID] = t1
		return true
	}
	
	// Structural unification
	if array1, ok := t1.(*ArrayType); ok {
		if array2, ok := t2.(*ArrayType); ok {
			return Unify(array1.ElemType, array2.ElemType, subst)
		}
		return false
	}
	
	if map1, ok := t1.(*MapType); ok {
		if map2, ok := t2.(*MapType); ok {
			return Unify(map1.KeyType, map2.KeyType, subst) &&
				Unify(map1.ValueType, map2.ValueType, subst)
		}
		return false
	}
	
	if fn1, ok := t1.(*FunctionType); ok {
		if fn2, ok := t2.(*FunctionType); ok {
			if len(fn1.ParamTypes) != len(fn2.ParamTypes) {
				return false
			}
			for i := range fn1.ParamTypes {
				if !Unify(fn1.ParamTypes[i], fn2.ParamTypes[i], subst) {
					return false
				}
			}
			return Unify(fn1.ReturnType, fn2.ReturnType, subst)
		}
		return false
	}
	
	// Direct equality check
	return t1.Equals(t2)
}

// ApplySubstitution applies a type variable substitution to a type
func ApplySubstitution(t Type, subst map[int]Type) Type {
	if tv, ok := t.(*TypeVar); ok {
		if replacement, found := subst[tv.ID]; found {
			// Recursively apply substitution
			return ApplySubstitution(replacement, subst)
		}
		return t
	}
	
	if array, ok := t.(*ArrayType); ok {
		return &ArrayType{ElemType: ApplySubstitution(array.ElemType, subst)}
	}
	
	if mapType, ok := t.(*MapType); ok {
		return &MapType{
			KeyType:   ApplySubstitution(mapType.KeyType, subst),
			ValueType: ApplySubstitution(mapType.ValueType, subst),
		}
	}
	
	if fn, ok := t.(*FunctionType); ok {
		paramTypes := make([]Type, len(fn.ParamTypes))
		for i, p := range fn.ParamTypes {
			paramTypes[i] = ApplySubstitution(p, subst)
		}
		return &FunctionType{
			ParamTypes: paramTypes,
			ReturnType: ApplySubstitution(fn.ReturnType, subst),
		}
	}
	
	return t
}

// NewTypeVar creates a new unique type variable
var typeVarCounter = 0

func NewTypeVar(name string) *TypeVar {
	typeVarCounter++
	return &TypeVar{Name: name, ID: typeVarCounter}
}

// ==================== Type Constructors ====================

// NewArrayType creates a new array type
func NewArrayType(elemType Type) *ArrayType {
	return &ArrayType{ElemType: elemType}
}

// NewMapType creates a new map type
func NewMapType(keyType, valueType Type) *MapType {
	return &MapType{KeyType: keyType, ValueType: valueType}
}

// NewFunctionType creates a new function type
func NewFunctionType(paramTypes []Type, returnType Type) *FunctionType {
	return &FunctionType{ParamTypes: paramTypes, ReturnType: returnType}
}

// NewStructType creates a new struct type
func NewStructType(name string, fields map[string]Type) *StructType {
	return &StructType{Name: name, Fields: fields}
}

// NewEnumType creates a new enum type
func NewEnumType(name string, variants map[string][]Type) *EnumType {
	return &EnumType{Name: name, Variants: variants}
}

