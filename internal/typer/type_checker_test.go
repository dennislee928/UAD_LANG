package typer

import (
	"testing"

	"github.com/dennislee928/uad-lang/internal/lexer"
	"github.com/dennislee928/uad-lang/internal/parser"
)

func parseAndCheck(t *testing.T, source string) (*TypeChecker, error) {
	l := lexer.New(source, "test.uad")
	tokens := l.AllTokens()
	
	p := parser.New(tokens, "test.uad")
	module, err := p.ParseModule()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	
	tc := NewTypeChecker()
	err = tc.Check(module)
	
	return tc, err
}

func TestTypeChecker_SimpleLet(t *testing.T) {
	source := `
		fn main() {
			let x = 42;
			let y = 3.14;
			let s = "hello";
			let b = true;
		}
	`
	
	_, err := parseAndCheck(t, source)
	if err != nil {
		t.Fatalf("Type check failed: %v", err)
	}
}

func TestTypeChecker_LetWithTypeAnnotation(t *testing.T) {
	source := `
		fn main() {
			let x: Int = 42;
			let y: Float = 3.14;
		}
	`
	
	_, err := parseAndCheck(t, source)
	if err != nil {
		t.Fatalf("Type check failed: %v", err)
	}
}

func TestTypeChecker_TypeMismatch(t *testing.T) {
	source := `
		fn main() {
			let x: Int = "hello";
		}
	`
	
	_, err := parseAndCheck(t, source)
	if err == nil {
		t.Fatal("Expected type error, got nil")
	}
}

func TestTypeChecker_BinaryArithmetic(t *testing.T) {
	source := `
		fn main() {
			let x = 1 + 2;
			let y = 3.0 * 4.0;
			let z = 5 - 6;
		}
	`
	
	_, err := parseAndCheck(t, source)
	if err != nil {
		t.Fatalf("Type check failed: %v", err)
	}
}

func TestTypeChecker_BinaryComparison(t *testing.T) {
	source := `
		fn main() {
			let a = 1 < 2;
			let b = 3.0 >= 4.0;
			let c = 5 == 6;
			let d = true != false;
		}
	`
	
	_, err := parseAndCheck(t, source)
	if err != nil {
		t.Fatalf("Type check failed: %v", err)
	}
}

func TestTypeChecker_BinaryLogical(t *testing.T) {
	source := `
		fn main() {
			let a = true && false;
			let b = true || false;
		}
	`
	
	_, err := parseAndCheck(t, source)
	if err != nil {
		t.Fatalf("Type check failed: %v", err)
	}
}

func TestTypeChecker_InvalidBinaryOp(t *testing.T) {
	source := `
		fn main() {
			let x = "hello" + 42;
		}
	`
	
	_, err := parseAndCheck(t, source)
	if err == nil {
		t.Fatal("Expected type error, got nil")
	}
}

func TestTypeChecker_FunctionCall(t *testing.T) {
	source := `
		fn add(x: Int, y: Int) -> Int {
			return x + y;
		}
		
		fn main() {
			let result = add(1, 2);
		}
	`
	
	_, err := parseAndCheck(t, source)
	if err != nil {
		t.Fatalf("Type check failed: %v", err)
	}
}

func TestTypeChecker_FunctionCallWrongArgType(t *testing.T) {
	source := `
		fn add(x: Int, y: Int) -> Int {
			return x + y;
		}
		
		fn main() {
			let result = add("hello", 2);
		}
	`
	
	_, err := parseAndCheck(t, source)
	if err == nil {
		t.Fatal("Expected type error, got nil")
	}
}

func TestTypeChecker_FunctionCallWrongArgCount(t *testing.T) {
	source := `
		fn add(x: Int, y: Int) -> Int {
			return x + y;
		}
		
		fn main() {
			let result = add(1);
		}
	`
	
	_, err := parseAndCheck(t, source)
	if err == nil {
		t.Fatal("Expected type error, got nil")
	}
}

func TestTypeChecker_IfExpression(t *testing.T) {
	source := `
		fn main() {
			let x = if true { 10 } else { 20 };
		}
	`
	
	_, err := parseAndCheck(t, source)
	if err != nil {
		t.Fatalf("Type check failed: %v", err)
	}
}

func TestTypeChecker_IfWithDifferentTypes(t *testing.T) {
	source := `
		fn main() {
			let x = if true { 10 } else { "hello" };
		}
	`
	
	_, err := parseAndCheck(t, source)
	if err == nil {
		t.Fatal("Expected type error, got nil")
	}
}

func TestTypeChecker_WhileLoop(t *testing.T) {
	source := `
		fn main() {
			let x = 0;
			while x < 10 {
				x = x + 1;
			}
		}
	`
	
	_, err := parseAndCheck(t, source)
	if err != nil {
		t.Fatalf("Type check failed: %v", err)
	}
}

func TestTypeChecker_ForLoop(t *testing.T) {
	source := `
		fn main() {
			let arr = [1, 2, 3];
			for x in arr {
				let y = x + 1;
			}
		}
	`
	
	_, err := parseAndCheck(t, source)
	if err != nil {
		t.Fatalf("Type check failed: %v", err)
	}
}

func TestTypeChecker_StructDeclaration(t *testing.T) {
	source := `
		struct Point {
			x: Float,
			y: Float,
		}
		
		fn main() {
			let p = Point { x: 1.0, y: 2.0 };
		}
	`
	
	_, err := parseAndCheck(t, source)
	if err != nil {
		t.Fatalf("Type check failed: %v", err)
	}
}

func TestTypeChecker_StructFieldAccess(t *testing.T) {
	source := `
		struct Point {
			x: Float,
			y: Float,
		}
		
		fn main() {
			let p = Point { x: 1.0, y: 2.0 };
			let x_val = p.x;
		}
	`
	
	_, err := parseAndCheck(t, source)
	if err != nil {
		t.Fatalf("Type check failed: %v", err)
	}
}

func TestTypeChecker_StructMissingField(t *testing.T) {
	source := `
		struct Point {
			x: Float,
			y: Float,
		}
		
		fn main() {
			let p = Point { x: 1.0 };
		}
	`
	
	_, err := parseAndCheck(t, source)
	if err == nil {
		t.Fatal("Expected type error for missing field, got nil")
	}
}

func TestTypeChecker_StructInvalidField(t *testing.T) {
	source := `
		struct Point {
			x: Float,
			y: Float,
		}
		
		fn main() {
			let p = Point { x: 1.0, y: 2.0, z: 3.0 };
		}
	`
	
	_, err := parseAndCheck(t, source)
	if err == nil {
		t.Fatal("Expected type error for invalid field, got nil")
	}
}

func TestTypeChecker_ArrayLiteral(t *testing.T) {
	source := `
		fn main() {
			let arr = [1, 2, 3];
		}
	`
	
	_, err := parseAndCheck(t, source)
	if err != nil {
		t.Fatalf("Type check failed: %v", err)
	}
}

func TestTypeChecker_ArrayIndex(t *testing.T) {
	source := `
		fn main() {
			let arr = [1, 2, 3];
			let first = arr[0];
		}
	`
	
	_, err := parseAndCheck(t, source)
	if err != nil {
		t.Fatalf("Type check failed: %v", err)
	}
}

func TestTypeChecker_ArrayMixedTypes(t *testing.T) {
	source := `
		fn main() {
			let arr = [1, "hello", 3];
		}
	`
	
	_, err := parseAndCheck(t, source)
	if err == nil {
		t.Fatal("Expected type error for mixed array types, got nil")
	}
}

func TestTypeChecker_ReturnType(t *testing.T) {
	source := `
		fn get_number() -> Int {
			return 42;
		}
		
		fn main() {
			let x = get_number();
		}
	`
	
	_, err := parseAndCheck(t, source)
	if err != nil {
		t.Fatalf("Type check failed: %v", err)
	}
}

func TestTypeChecker_ReturnTypeMismatch(t *testing.T) {
	source := `
		fn get_number() -> Int {
			return "hello";
		}
	`
	
	_, err := parseAndCheck(t, source)
	if err == nil {
		t.Fatal("Expected type error for return type mismatch, got nil")
	}
}

func TestTypeChecker_BuiltinFunctions(t *testing.T) {
	source := `
		fn main() {
			print("Hello");
			let x = sqrt(9.0);
			let y = abs(-5.0);
		}
	`
	
	_, err := parseAndCheck(t, source)
	if err != nil {
		t.Fatalf("Type check failed: %v", err)
	}
}

func TestTypeChecker_UnaryOperators(t *testing.T) {
	source := `
		fn main() {
			let x = -42;
			let y = !true;
		}
	`
	
	_, err := parseAndCheck(t, source)
	if err != nil {
		t.Fatalf("Type check failed: %v", err)
	}
}

func TestTypeChecker_ComplexExpression(t *testing.T) {
	source := `
		fn main() {
			let x = (1 + 2) * 3 - 4 / 2;
			let y = x > 5 && x < 10;
		}
	`
	
	_, err := parseAndCheck(t, source)
	if err != nil {
		t.Fatalf("Type check failed: %v", err)
	}
}

func TestTypeChecker_NestedScopes(t *testing.T) {
	source := `
		fn main() {
			let x = 10;
			{
				let y = 20;
				let z = x + y;
			}
		}
	`
	
	_, err := parseAndCheck(t, source)
	if err != nil {
		t.Fatalf("Type check failed: %v", err)
	}
}

func TestTypeChecker_RecursiveFunction(t *testing.T) {
	source := `
		fn factorial(n: Int) -> Int {
			if n <= 1 {
				return 1;
			} else {
				return n * factorial(n - 1);
			}
		}
		
		fn main() {
			let result = factorial(5);
		}
	`
	
	_, err := parseAndCheck(t, source)
	if err != nil {
		t.Fatalf("Type check failed: %v", err)
	}
}

func TestTypeChecker_EnumDeclaration(t *testing.T) {
	source := `
		enum Result {
			Ok(Int),
			Err(String),
		}
		
		fn main() {
			let x = 42;
		}
	`
	
	_, err := parseAndCheck(t, source)
	if err != nil {
		t.Fatalf("Type check failed: %v", err)
	}
}

func TestTypeChecker_TypeAlias(t *testing.T) {
	source := `
		type Number = Int;
		
		fn main() {
			let x: Number = 42;
		}
	`
	
	_, err := parseAndCheck(t, source)
	if err != nil {
		t.Fatalf("Type check failed: %v", err)
	}
}

func BenchmarkTypeChecker_SimpleProgram(b *testing.B) {
	source := `
		fn add(x: Int, y: Int) -> Int {
			return x + y;
		}
		
		fn main() {
			let result = add(1, 2);
		}
	`
	
	l := lexer.New(source, "bench.uad")
	tokens := l.AllTokens()
	
	p := parser.New(tokens, "bench.uad")
	module, _ := p.ParseModule()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tc := NewTypeChecker()
		_ = tc.Check(module)
	}
}

func BenchmarkTypeChecker_ComplexProgram(b *testing.B) {
	source := `
		struct Point {
			x: Float,
			y: Float,
		}
		
		fn distance(p1: Point, p2: Point) -> Float {
			let dx = p2.x - p1.x;
			let dy = p2.y - p1.y;
			return sqrt(dx * dx + dy * dy);
		}
		
		fn main() {
			let points = [
				Point { x: 0.0, y: 0.0 },
				Point { x: 3.0, y: 4.0 },
			];
			
			for p in points {
				let d = distance(points[0], p);
			}
		}
	`
	
	l := lexer.New(source, "bench.uad")
	tokens := l.AllTokens()
	
	p := parser.New(tokens, "bench.uad")
	module, _ := p.ParseModule()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tc := NewTypeChecker()
		_ = tc.Check(module)
	}
}

