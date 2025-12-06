package vm

import (
	"fmt"
	"math"

	"github.com/dennislee928/uad-lang/internal/ir"
)

// Value represents a runtime value in the VM
type Value struct {
	Kind  ir.ValueKind
	Data  interface{} // int64, float64, bool, string, *Function, []Value, map[string]Value
}

// NewIntValue creates an integer value
func NewIntValue(v int64) Value {
	return Value{Kind: ir.ValueInt, Data: v}
}

// NewFloatValue creates a float value
func NewFloatValue(v float64) Value {
	return Value{Kind: ir.ValueFloat, Data: v}
}

// NewBoolValue creates a boolean value
func NewBoolValue(v bool) Value {
	return Value{Kind: ir.ValueBool, Data: v}
}

// NewStringValue creates a string value
func NewStringValue(v string) Value {
	return Value{Kind: ir.ValueString, Data: v}
}

// NewNilValue creates a nil value
func NewNilValue() Value {
	return Value{Kind: ir.ValueNil, Data: nil}
}

// NewArrayValue creates an array value
func NewArrayValue(v []Value) Value {
	return Value{Kind: ir.ValueArray, Data: v}
}

// String returns a string representation of the value
func (v Value) String() string {
	switch v.Kind {
	case ir.ValueInt:
		return fmt.Sprintf("%d", v.Data.(int64))
	case ir.ValueFloat:
		return fmt.Sprintf("%g", v.Data.(float64))
	case ir.ValueBool:
		if v.Data.(bool) {
			return "true"
		}
		return "false"
	case ir.ValueString:
		return v.Data.(string)
	case ir.ValueNil:
		return "nil"
	case ir.ValueArray:
		return fmt.Sprintf("%v", v.Data)
	default:
		return fmt.Sprintf("<value %s>", v.Kind.String())
	}
}

// IsTruthy returns true if the value is truthy
func (v Value) IsTruthy() bool {
	switch v.Kind {
	case ir.ValueBool:
		return v.Data.(bool)
	case ir.ValueNil:
		return false
	case ir.ValueInt:
		return v.Data.(int64) != 0
	case ir.ValueFloat:
		return v.Data.(float64) != 0.0
	case ir.ValueString:
		return v.Data.(string) != ""
	default:
		return true
	}
}

// ==================== VM ====================

// VM represents the virtual machine
type VM struct {
	module      *ir.Module
	stack       []Value
	sp          int // Stack pointer
	globals     []Value
	callStack   []CallFrame
	
	// Built-in functions
	builtins map[string]func(*VM, []Value) (Value, error)
}

// CallFrame represents a call frame
type CallFrame struct {
	function *ir.Function
	ip       int32  // Instruction pointer
	bp       int    // Base pointer (for locals)
	locals   []Value
}

// New creates a new VM
func New(module *ir.Module) *VM {
	vm := &VM{
		module:    module,
		stack:     make([]Value, 0, 1024),
		sp:        0,
		globals:   make([]Value, len(module.Globals)),
		callStack: make([]CallFrame, 0, 256),
		builtins:  make(map[string]func(*VM, []Value) (Value, error)),
	}
	
	vm.initBuiltins()
	return vm
}

// Run executes the main function
func (vm *VM) Run() error {
	// Find main function
	mainFn, ok := vm.module.GetFunction("main")
	if !ok {
		return fmt.Errorf("main function not found")
	}
	
	// Call main
	_, err := vm.callFunction(mainFn, nil)
	return err
}

// ==================== Stack Operations ====================

func (vm *VM) push(v Value) {
	if vm.sp >= len(vm.stack) {
		vm.stack = append(vm.stack, v)
	} else {
		vm.stack[vm.sp] = v
	}
	vm.sp++
}

func (vm *VM) pop() Value {
	if vm.sp == 0 {
		panic("stack underflow")
	}
	vm.sp--
	return vm.stack[vm.sp]
}

func (vm *VM) peek(offset int) Value {
	return vm.stack[vm.sp-1-offset]
}

// ==================== Function Execution ====================

func (vm *VM) callFunction(fn *ir.Function, args []Value) (Value, error) {
	// Create call frame
	frame := CallFrame{
		function: fn,
		ip:       0,
		bp:       vm.sp - len(args),
		locals:   make([]Value, len(fn.Locals)+len(fn.Params)),
	}
	
	// Copy arguments to locals
	for i, arg := range args {
		frame.locals[fn.Params[i].Index] = arg
	}
	
	// Push frame
	vm.callStack = append(vm.callStack, frame)
	
	// Execute function
	result, err := vm.executeFrame()
	
	// Pop frame
	vm.callStack = vm.callStack[:len(vm.callStack)-1]
	
	return result, err
}

func (vm *VM) executeFrame() (Value, error) {
	frame := &vm.callStack[len(vm.callStack)-1]
	fn := frame.function
	
	for int(frame.ip) < len(fn.Instructions) {
		instr := fn.Instructions[frame.ip]
		frame.ip++
		
		if err := vm.executeInstruction(instr); err != nil {
			return NewNilValue(), err
		}
		
		// Check for return
		if instr.Op == ir.OpReturn {
			if vm.sp > 0 {
				return vm.pop(), nil
			}
			return NewNilValue(), nil
		}
	}
	
	// Implicit return nil
	return NewNilValue(), nil
}

func (vm *VM) executeInstruction(instr ir.Instruction) error {
	frame := &vm.callStack[len(vm.callStack)-1]
	
	switch instr.Op {
	case ir.OpNop:
		// Do nothing
		
	case ir.OpPop:
		vm.pop()
		
	case ir.OpDup:
		v := vm.peek(0)
		vm.push(v)
		
	case ir.OpSwap:
		a := vm.pop()
		b := vm.pop()
		vm.push(a)
		vm.push(b)
		
	// Constants
	case ir.OpConstInt:
		c, err := frame.function.GetConstant(instr.Operand)
		if err != nil {
			return err
		}
		vm.push(NewIntValue(c.Value.(int64)))
		
	case ir.OpConstFloat:
		c, err := frame.function.GetConstant(instr.Operand)
		if err != nil {
			return err
		}
		vm.push(NewFloatValue(c.Value.(float64)))
		
	case ir.OpConstBool:
		c, err := frame.function.GetConstant(instr.Operand)
		if err != nil {
			return err
		}
		vm.push(NewBoolValue(c.Value.(bool)))
		
	case ir.OpConstString:
		c, err := frame.function.GetConstant(instr.Operand)
		if err != nil {
			return err
		}
		vm.push(NewStringValue(c.Value.(string)))
		
	case ir.OpConstNil:
		vm.push(NewNilValue())
		
	// Arithmetic
	case ir.OpAdd:
		return vm.execAdd()
	case ir.OpSub:
		return vm.execSub()
	case ir.OpMul:
		return vm.execMul()
	case ir.OpDiv:
		return vm.execDiv()
	case ir.OpMod:
		return vm.execMod()
	case ir.OpNeg:
		return vm.execNeg()
		
	// Comparison
	case ir.OpEq:
		return vm.execEq()
	case ir.OpNeq:
		return vm.execNeq()
	case ir.OpLt:
		return vm.execLt()
	case ir.OpGt:
		return vm.execGt()
	case ir.OpLe:
		return vm.execLe()
	case ir.OpGe:
		return vm.execGe()
		
	// Logical
	case ir.OpAnd:
		return vm.execAnd()
	case ir.OpOr:
		return vm.execOr()
	case ir.OpNot:
		return vm.execNot()
		
	// Variables
	case ir.OpGetLocal:
		v := frame.locals[instr.Operand]
		vm.push(v)
		
	case ir.OpSetLocal:
		v := vm.pop()
		frame.locals[instr.Operand] = v
		
	case ir.OpGetGlobal:
		v := vm.globals[instr.Operand]
		vm.push(v)
		
	case ir.OpSetGlobal:
		v := vm.pop()
		vm.globals[instr.Operand] = v
		
	// Control flow
	case ir.OpJump:
		frame.ip = instr.Operand
		
	case ir.OpJumpIfFalse:
		cond := vm.pop()
		if !cond.IsTruthy() {
			frame.ip = instr.Operand
		}
		
	case ir.OpJumpIfTrue:
		cond := vm.pop()
		if cond.IsTruthy() {
			frame.ip = instr.Operand
		}
		
	case ir.OpCall:
		return vm.execCall(int(instr.Operand))
		
	case ir.OpReturn:
		// Handled in executeFrame
		
	// Memory operations
	case ir.OpNewArray:
		return vm.execNewArray(int(instr.Operand))
		
	case ir.OpGetIndex:
		return vm.execGetIndex()
		
	default:
		return fmt.Errorf("unimplemented opcode: %s", instr.Op.String())
	}
	
	return nil
}

// ==================== Instruction Implementations ====================

func (vm *VM) execAdd() error {
	b := vm.pop()
	a := vm.pop()
	
	if a.Kind == ir.ValueInt && b.Kind == ir.ValueInt {
		vm.push(NewIntValue(a.Data.(int64) + b.Data.(int64)))
	} else {
		// Promote to float
		af := vm.toFloat(a)
		bf := vm.toFloat(b)
		vm.push(NewFloatValue(af + bf))
	}
	
	return nil
}

func (vm *VM) execSub() error {
	b := vm.pop()
	a := vm.pop()
	
	if a.Kind == ir.ValueInt && b.Kind == ir.ValueInt {
		vm.push(NewIntValue(a.Data.(int64) - b.Data.(int64)))
	} else {
		af := vm.toFloat(a)
		bf := vm.toFloat(b)
		vm.push(NewFloatValue(af - bf))
	}
	
	return nil
}

func (vm *VM) execMul() error {
	b := vm.pop()
	a := vm.pop()
	
	if a.Kind == ir.ValueInt && b.Kind == ir.ValueInt {
		vm.push(NewIntValue(a.Data.(int64) * b.Data.(int64)))
	} else {
		af := vm.toFloat(a)
		bf := vm.toFloat(b)
		vm.push(NewFloatValue(af * bf))
	}
	
	return nil
}

func (vm *VM) execDiv() error {
	b := vm.pop()
	a := vm.pop()
	
	if a.Kind == ir.ValueInt && b.Kind == ir.ValueInt {
		bv := b.Data.(int64)
		if bv == 0 {
			return fmt.Errorf("division by zero")
		}
		vm.push(NewIntValue(a.Data.(int64) / bv))
	} else {
		bf := vm.toFloat(b)
		if bf == 0 {
			return fmt.Errorf("division by zero")
		}
		af := vm.toFloat(a)
		vm.push(NewFloatValue(af / bf))
	}
	
	return nil
}

func (vm *VM) execMod() error {
	b := vm.pop()
	a := vm.pop()
	
	if a.Kind != ir.ValueInt || b.Kind != ir.ValueInt {
		return fmt.Errorf("modulo requires integers")
	}
	
	bv := b.Data.(int64)
	if bv == 0 {
		return fmt.Errorf("modulo by zero")
	}
	
	vm.push(NewIntValue(a.Data.(int64) % bv))
	return nil
}

func (vm *VM) execNeg() error {
	v := vm.pop()
	
	if v.Kind == ir.ValueInt {
		vm.push(NewIntValue(-v.Data.(int64)))
	} else if v.Kind == ir.ValueFloat {
		vm.push(NewFloatValue(-v.Data.(float64)))
	} else {
		return fmt.Errorf("negation requires numeric type")
	}
	
	return nil
}

func (vm *VM) execEq() error {
	b := vm.pop()
	a := vm.pop()
	vm.push(NewBoolValue(vm.equals(a, b)))
	return nil
}

func (vm *VM) execNeq() error {
	b := vm.pop()
	a := vm.pop()
	vm.push(NewBoolValue(!vm.equals(a, b)))
	return nil
}

func (vm *VM) execLt() error {
	b := vm.pop()
	a := vm.pop()
	
	af := vm.toFloat(a)
	bf := vm.toFloat(b)
	vm.push(NewBoolValue(af < bf))
	return nil
}

func (vm *VM) execGt() error {
	b := vm.pop()
	a := vm.pop()
	
	af := vm.toFloat(a)
	bf := vm.toFloat(b)
	vm.push(NewBoolValue(af > bf))
	return nil
}

func (vm *VM) execLe() error {
	b := vm.pop()
	a := vm.pop()
	
	af := vm.toFloat(a)
	bf := vm.toFloat(b)
	vm.push(NewBoolValue(af <= bf))
	return nil
}

func (vm *VM) execGe() error {
	b := vm.pop()
	a := vm.pop()
	
	af := vm.toFloat(a)
	bf := vm.toFloat(b)
	vm.push(NewBoolValue(af >= bf))
	return nil
}

func (vm *VM) execAnd() error {
	b := vm.pop()
	a := vm.pop()
	vm.push(NewBoolValue(a.IsTruthy() && b.IsTruthy()))
	return nil
}

func (vm *VM) execOr() error {
	b := vm.pop()
	a := vm.pop()
	vm.push(NewBoolValue(a.IsTruthy() || b.IsTruthy()))
	return nil
}

func (vm *VM) execNot() error {
	v := vm.pop()
	vm.push(NewBoolValue(!v.IsTruthy()))
	return nil
}

func (vm *VM) execCall(argCount int) error {
	// Get function value
	fnVal := vm.pop()
	
	// Check if it's a builtin
	if fnVal.Kind == ir.ValueString {
		name := fnVal.Data.(string)
		if builtin, ok := vm.builtins[name]; ok {
			// Collect arguments
			args := make([]Value, argCount)
			for i := argCount - 1; i >= 0; i-- {
				args[i] = vm.pop()
			}
			
			// Call builtin
			result, err := builtin(vm, args)
			if err != nil {
				return err
			}
			
			vm.push(result)
			return nil
		}
	}
	
	// Try to find function in module
	if fnVal.Kind == ir.ValueString {
		name := fnVal.Data.(string)
		fn, ok := vm.module.GetFunction(name)
		if !ok {
			return fmt.Errorf("undefined function: %s", name)
		}
		
		// Collect arguments
		args := make([]Value, argCount)
		for i := argCount - 1; i >= 0; i-- {
			args[i] = vm.pop()
		}
		
		// Call function
		result, err := vm.callFunction(fn, args)
		if err != nil {
			return err
		}
		
		vm.push(result)
		return nil
	}
	
	return fmt.Errorf("cannot call non-function value")
}

func (vm *VM) execNewArray(size int) error {
	arr := make([]Value, size)
	for i := size - 1; i >= 0; i-- {
		arr[i] = vm.pop()
	}
	vm.push(NewArrayValue(arr))
	return nil
}

func (vm *VM) execGetIndex() error {
	index := vm.pop()
	target := vm.pop()
	
	if target.Kind == ir.ValueArray {
		arr := target.Data.([]Value)
		if index.Kind != ir.ValueInt {
			return fmt.Errorf("array index must be integer")
		}
		idx := index.Data.(int64)
		if idx < 0 || idx >= int64(len(arr)) {
			return fmt.Errorf("array index out of bounds")
		}
		vm.push(arr[idx])
		return nil
	}
	
	return fmt.Errorf("index operation requires array")
}

// ==================== Helper Functions ====================

func (vm *VM) toFloat(v Value) float64 {
	switch v.Kind {
	case ir.ValueInt:
		return float64(v.Data.(int64))
	case ir.ValueFloat:
		return v.Data.(float64)
	default:
		return 0
	}
}

func (vm *VM) equals(a, b Value) bool {
	if a.Kind != b.Kind {
		return false
	}
	
	switch a.Kind {
	case ir.ValueInt:
		return a.Data.(int64) == b.Data.(int64)
	case ir.ValueFloat:
		return a.Data.(float64) == b.Data.(float64)
	case ir.ValueBool:
		return a.Data.(bool) == b.Data.(bool)
	case ir.ValueString:
		return a.Data.(string) == b.Data.(string)
	case ir.ValueNil:
		return true
	default:
		return false
	}
}

// ==================== Built-in Functions ====================

func (vm *VM) initBuiltins() {
	vm.builtins["print"] = vm.builtinPrint
	vm.builtins["println"] = vm.builtinPrintln
	vm.builtins["sqrt"] = vm.builtinSqrt
	vm.builtins["abs"] = vm.builtinAbs
	vm.builtins["pow"] = vm.builtinPow
	vm.builtins["sin"] = vm.builtinSin
	vm.builtins["cos"] = vm.builtinCos
	vm.builtins["tan"] = vm.builtinTan
}

func (vm *VM) builtinPrint(self *VM, args []Value) (Value, error) {
	for i, arg := range args {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(arg.String())
	}
	return NewNilValue(), nil
}

func (vm *VM) builtinPrintln(self *VM, args []Value) (Value, error) {
	for i, arg := range args {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(arg.String())
	}
	fmt.Println()
	return NewNilValue(), nil
}

func (vm *VM) builtinSqrt(self *VM, args []Value) (Value, error) {
	if len(args) != 1 {
		return NewNilValue(), fmt.Errorf("sqrt expects 1 argument")
	}
	f := vm.toFloat(args[0])
	return NewFloatValue(math.Sqrt(f)), nil
}

func (vm *VM) builtinAbs(self *VM, args []Value) (Value, error) {
	if len(args) != 1 {
		return NewNilValue(), fmt.Errorf("abs expects 1 argument")
	}
	f := vm.toFloat(args[0])
	return NewFloatValue(math.Abs(f)), nil
}

func (vm *VM) builtinPow(self *VM, args []Value) (Value, error) {
	if len(args) != 2 {
		return NewNilValue(), fmt.Errorf("pow expects 2 arguments")
	}
	a := vm.toFloat(args[0])
	b := vm.toFloat(args[1])
	return NewFloatValue(math.Pow(a, b)), nil
}

func (vm *VM) builtinSin(self *VM, args []Value) (Value, error) {
	if len(args) != 1 {
		return NewNilValue(), fmt.Errorf("sin expects 1 argument")
	}
	f := vm.toFloat(args[0])
	return NewFloatValue(math.Sin(f)), nil
}

func (vm *VM) builtinCos(self *VM, args []Value) (Value, error) {
	if len(args) != 1 {
		return NewNilValue(), fmt.Errorf("cos expects 1 argument")
	}
	f := vm.toFloat(args[0])
	return NewFloatValue(math.Cos(f)), nil
}

func (vm *VM) builtinTan(self *VM, args []Value) (Value, error) {
	if len(args) != 1 {
		return NewNilValue(), fmt.Errorf("tan expects 1 argument")
	}
	f := vm.toFloat(args[0])
	return NewFloatValue(math.Tan(f)), nil
}

