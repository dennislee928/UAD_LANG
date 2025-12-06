package lexer

import (
	"testing"
)

func TestLexer_Keywords(t *testing.T) {
	source := `fn let if else return while for break continue struct enum type
	import module pub true false nil match case when in action judge agent event emit
	action_class erh_profile scenario cognitive_siem from dataset prime_threshold fit_alpha`
	
	expected := []TokenType{
		TokenFn, TokenLet, TokenIf, TokenElse, TokenReturn, TokenWhile, TokenFor,
		TokenBreak, TokenContinue, TokenStruct, TokenEnum, TokenTypeDecl, TokenImport,
		TokenModule, TokenPub, TokenTrue, TokenFalse, TokenNil, TokenMatch,
		TokenCase, TokenWhen, TokenIn, TokenAction, TokenJudge, TokenAgent,
		TokenEvent, TokenEmit, TokenActionClass, TokenErhProfile, TokenScenario,
		TokenCognitiveSiem, TokenFrom, TokenDataset, TokenPrimeThreshold, TokenFitAlpha,
		TokenEOF,
	}
	
	l := New(source, "test.uad")
	for i, exp := range expected {
		tok := l.NextToken()
		if tok.Type != exp {
			t.Errorf("Token %d: expected %s, got %s", i, exp, tok.Type)
		}
	}
}

func TestLexer_Identifiers(t *testing.T) {
	source := `myVar _internal MyStruct calculate_alpha x1 y2`
	
	expected := []struct {
		typ    TokenType
		lexeme string
	}{
		{TokenIdent, "myVar"},
		{TokenIdent, "_internal"},
		{TokenIdent, "MyStruct"},
		{TokenIdent, "calculate_alpha"},
		{TokenIdent, "x1"},
		{TokenIdent, "y2"},
		{TokenEOF, ""},
	}
	
	l := New(source, "test.uad")
	for i, exp := range expected {
		tok := l.NextToken()
		if tok.Type != exp.typ {
			t.Errorf("Token %d: expected type %s, got %s", i, exp.typ, tok.Type)
		}
		if tok.Lexeme != exp.lexeme {
			t.Errorf("Token %d: expected lexeme '%s', got '%s'", i, exp.lexeme, tok.Lexeme)
		}
	}
}

func TestLexer_Numbers(t *testing.T) {
	tests := []struct {
		source string
		typ    TokenType
		lexeme string
	}{
		{"42", TokenInt, "42"},
		{"0", TokenInt, "0"},
		{"123456", TokenInt, "123456"},
		{"3.14", TokenFloat, "3.14"},
		{"0.5", TokenFloat, "0.5"},
		{"1.5e10", TokenFloat, "1.5e10"},
		{"2.0E+5", TokenFloat, "2.0E+5"},
		{"1.5e-10", TokenFloat, "1.5e-10"},
		{"0xFF", TokenInt, "0xFF"},
		{"0xABCD", TokenInt, "0xABCD"},
		{"0b1010", TokenInt, "0b1010"},
		{"0b11111111", TokenInt, "0b11111111"},
		{"10s", TokenDuration, "10s"},
		{"5m", TokenDuration, "5m"},
		{"2h", TokenDuration, "2h"},
		{"90d", TokenDuration, "90d"},
	}
	
	for _, tt := range tests {
		l := New(tt.source, "test.uad")
		tok := l.NextToken()
		if tok.Type != tt.typ {
			t.Errorf("Source '%s': expected type %s, got %s", tt.source, tt.typ, tok.Type)
		}
		if tok.Lexeme != tt.lexeme {
			t.Errorf("Source '%s': expected lexeme '%s', got '%s'", tt.source, tt.lexeme, tok.Lexeme)
		}
	}
}

func TestLexer_Strings(t *testing.T) {
	tests := []struct {
		source string
		lexeme string
	}{
		{`"hello"`, `"hello"`},
		{`"Hello, world!"`, `"Hello, world!"`},
		{`"Line 1\nLine 2"`, `"Line 1\nLine 2"`},
		{`"Tab\there"`, `"Tab\there"`},
		{`"Quote: \""`, `"Quote: \""`},
		{`"Backslash: \\"`, `"Backslash: \\"`},
		{`""`, `""`},
	}
	
	for _, tt := range tests {
		l := New(tt.source, "test.uad")
		tok := l.NextToken()
		if tok.Type != TokenString {
			t.Errorf("Source '%s': expected TokenString, got %s", tt.source, tok.Type)
		}
		if tok.Lexeme != tt.lexeme {
			t.Errorf("Source '%s': expected lexeme '%s', got '%s'", tt.source, tt.lexeme, tok.Lexeme)
		}
	}
}

func TestLexer_Operators(t *testing.T) {
	source := `+ - * / % == != < > <= >= && || ! = -> => . : :: ..`
	
	expected := []TokenType{
		TokenPlus, TokenMinus, TokenStar, TokenSlash, TokenMod,
		TokenEq, TokenNeq, TokenLt, TokenGt, TokenLe, TokenGe,
		TokenAnd, TokenOr, TokenNot, TokenAssign, TokenArrow, TokenFatArrow,
		TokenDot, TokenColon, TokenDoubleColon, TokenDotDot,
		TokenEOF,
	}
	
	l := New(source, "test.uad")
	for i, exp := range expected {
		tok := l.NextToken()
		if tok.Type != exp {
			t.Errorf("Token %d: expected %s, got %s", i, exp, tok.Type)
		}
	}
}

func TestLexer_Delimiters(t *testing.T) {
	source := `( ) { } [ ] , ;`
	
	expected := []TokenType{
		TokenLParen, TokenRParen, TokenLBrace, TokenRBrace,
		TokenLBracket, TokenRBracket, TokenComma, TokenSemicolon,
		TokenEOF,
	}
	
	l := New(source, "test.uad")
	for i, exp := range expected {
		tok := l.NextToken()
		if tok.Type != exp {
			t.Errorf("Token %d: expected %s, got %s", i, exp, tok.Type)
		}
	}
}

func TestLexer_Comments(t *testing.T) {
	source := `
		// This is a single-line comment
		let x = 10;
		/* This is a
		   multi-line comment */
		let y = 20;
	`
	
	expected := []TokenType{
		TokenComment, // //...
		TokenLet, TokenIdent, TokenAssign, TokenInt, TokenSemicolon,
		TokenComment, // /*...*/
		TokenLet, TokenIdent, TokenAssign, TokenInt, TokenSemicolon,
		TokenEOF,
	}
	
	l := New(source, "test.uad")
	for i, exp := range expected {
		tok := l.NextToken()
		if tok.Type != exp {
			t.Errorf("Token %d: expected %s, got %s", i, exp, tok.Type)
		}
	}
}

func TestLexer_CompleteFunction(t *testing.T) {
	source := `
		fn add(x: Int, y: Int) -> Int {
			return x + y;
		}
	`
	
	expected := []TokenType{
		TokenFn, TokenIdent, // fn add
		TokenLParen, TokenIdent, TokenColon, TokenIdent, TokenComma, // (x: Int,
		TokenIdent, TokenColon, TokenIdent, TokenRParen, // y: Int)
		TokenArrow, TokenIdent, // -> Int
		TokenLBrace, // {
		TokenReturn, TokenIdent, TokenPlus, TokenIdent, TokenSemicolon, // return x + y;
		TokenRBrace, // }
		TokenEOF,
	}
	
	l := New(source, "test.uad")
	for i, exp := range expected {
		tok := l.NextToken()
		if tok.Type != exp {
			t.Errorf("Token %d: expected %s, got %s (lexeme: %s)", i, exp, tok.Type, tok.Lexeme)
		}
	}
}

func TestLexer_Position(t *testing.T) {
	source := "let x = 10;\nlet y = 20;"
	
	l := New(source, "test.uad")
	
	// let (line 1, col 1)
	tok := l.NextToken()
	if tok.Span.Start.Line != 1 || tok.Span.Start.Column != 1 {
		t.Errorf("Expected position 1:1, got %d:%d", tok.Span.Start.Line, tok.Span.Start.Column)
	}
	
	// Skip to second line
	for i := 0; i < 4; i++ {
		l.NextToken()
	}
	
	// let (line 2, col 1)
	tok = l.NextToken()
	if tok.Span.Start.Line != 2 || tok.Span.Start.Column != 1 {
		t.Errorf("Expected position 2:1, got %d:%d", tok.Span.Start.Line, tok.Span.Start.Column)
	}
}

func TestLexer_UnterminatedString(t *testing.T) {
	source := `"unterminated`
	
	l := New(source, "test.uad")
	tok := l.NextToken()
	if tok.Type != TokenIllegal {
		t.Errorf("Expected TokenIllegal for unterminated string, got %s", tok.Type)
	}
}

func TestLexer_UnterminatedBlockComment(t *testing.T) {
	source := `/* unterminated comment`
	
	l := New(source, "test.uad")
	tok := l.NextToken()
	if tok.Type != TokenIllegal {
		t.Errorf("Expected TokenIllegal for unterminated block comment, got %s", tok.Type)
	}
}

func TestLexer_ModelDSL(t *testing.T) {
	source := `
		action_class MergeRequest {
			complexity = log(1 + lines_changed)
		}
		
		erh_profile "GitLab-DevSecOps" {
			actions from dataset "mr_security_logs"
		}
	`
	
	expected := []TokenType{
		TokenActionClass, TokenIdent, TokenLBrace, // action_class MergeRequest {
		TokenIdent, TokenAssign, TokenIdent, // complexity = log
		TokenLParen, TokenInt, TokenPlus, TokenIdent, TokenRParen, // (1 + lines_changed)
		TokenRBrace, // }
		TokenErhProfile, TokenString, TokenLBrace, // erh_profile "..." {
		TokenIdent, TokenFrom, TokenDataset, TokenString, // actions from dataset "..."
		TokenRBrace, // }
		TokenEOF,
	}
	
	l := New(source, "test.uadmodel")
	for i, exp := range expected {
		tok := l.NextToken()
		if tok.Type != exp {
			t.Errorf("Token %d: expected %s, got %s (lexeme: %s)", i, exp, tok.Type, tok.Lexeme)
		}
	}
}

func BenchmarkLexer_Keywords(b *testing.B) {
	source := "fn let if else return while for break continue struct enum"
	for i := 0; i < b.N; i++ {
		l := New(source, "bench.uad")
		for {
			tok := l.NextToken()
			if tok.Type == TokenEOF {
				break
			}
		}
	}
}

func BenchmarkLexer_LargeSource(b *testing.B) {
	// Generate a large source file
	source := ""
	for i := 0; i < 1000; i++ {
		source += "fn test" + string(rune('a'+i%26)) + "() { let x = 10; return x + 20; }\n"
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := New(source, "bench.uad")
		for {
			tok := l.NextToken()
			if tok.Type == TokenEOF {
				break
			}
		}
	}
}

