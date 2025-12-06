package lexer

import (
	"github.com/dennislee928/uad-lang/internal/common"
)

// TokenType represents the type of a token
type TokenType int

const (
	// Special tokens
	TokenEOF TokenType = iota
	TokenIllegal
	TokenComment

	// Identifiers and literals
	TokenIdent     // myVar
	TokenInt       // 42
	TokenFloat     // 3.14
	TokenString    // "hello"
	TokenDuration  // 10s, 5m, 2h, 3d

	// Keywords - Control flow
	TokenFn       // fn
	TokenReturn   // return
	TokenIf       // if
	TokenElse     // else
	TokenMatch    // match
	TokenWhile    // while
	TokenFor      // for
	TokenBreak    // break
	TokenContinue // continue

	// Keywords - Declarations
	TokenLet       // let
	TokenStruct    // struct
	TokenEnum      // enum
	TokenTypeDecl  // type
	TokenImport    // import
	TokenModule    // module
	TokenPub       // pub

	// Keywords - Literals
	TokenTrue  // true
	TokenFalse // false
	TokenNil   // nil

	// Keywords - Pattern matching
	TokenCase // case
	TokenWhen // when
	TokenIn   // in

	// Keywords - Domain-specific
	TokenAction // action
	TokenJudge  // judge
	TokenAgent  // agent
	TokenEvent  // event
	TokenEmit   // emit

	// Keywords - Model DSL
	TokenActionClass    // action_class
	TokenErhProfile     // erh_profile
	TokenScenario       // scenario
	TokenCognitiveSiem  // cognitive_siem
	TokenFrom           // from
	TokenDataset        // dataset
	TokenPrimeThreshold // prime_threshold
	TokenFitAlpha       // fit_alpha

	// Operators - Arithmetic
	TokenPlus  // +
	TokenMinus // -
	TokenStar  // *
	TokenSlash // /
	TokenMod   // %

	// Operators - Comparison
	TokenEq  // ==
	TokenNeq // !=
	TokenLt  // <
	TokenGt  // >
	TokenLe  // <=
	TokenGe  // >=

	// Operators - Logical
	TokenAnd // &&
	TokenOr  // ||
	TokenNot // !

	// Operators - Assignment
	TokenAssign // =

	// Operators - Other
	TokenArrow      // ->
	TokenFatArrow   // =>
	TokenDot        // .
	TokenColon      // :
	TokenDoubleColon // ::
	TokenDotDot     // ..

	// Delimiters
	TokenLParen    // (
	TokenRParen    // )
	TokenLBrace    // {
	TokenRBrace    // }
	TokenLBracket  // [
	TokenRBracket  // ]
	TokenComma     // ,
	TokenSemicolon // ;
)

// Token represents a lexical token
type Token struct {
	Type   TokenType
	Lexeme string
	Span   common.Span
}

// NewToken creates a new token
func NewToken(typ TokenType, lexeme string, span common.Span) Token {
	return Token{
		Type:   typ,
		Lexeme: lexeme,
		Span:   span,
	}
}

// String returns the string representation of the token type
func (t TokenType) String() string {
	switch t {
	case TokenEOF:
		return "EOF"
	case TokenIllegal:
		return "ILLEGAL"
	case TokenComment:
		return "COMMENT"
	case TokenIdent:
		return "IDENT"
	case TokenInt:
		return "INT"
	case TokenFloat:
		return "FLOAT"
	case TokenString:
		return "STRING"
	case TokenDuration:
		return "DURATION"
	case TokenFn:
		return "fn"
	case TokenReturn:
		return "return"
	case TokenIf:
		return "if"
	case TokenElse:
		return "else"
	case TokenMatch:
		return "match"
	case TokenWhile:
		return "while"
	case TokenFor:
		return "for"
	case TokenBreak:
		return "break"
	case TokenContinue:
		return "continue"
	case TokenLet:
		return "let"
	case TokenStruct:
		return "struct"
	case TokenEnum:
		return "enum"
	case TokenTypeDecl:
		return "type"
	case TokenImport:
		return "import"
	case TokenModule:
		return "module"
	case TokenPub:
		return "pub"
	case TokenTrue:
		return "true"
	case TokenFalse:
		return "false"
	case TokenNil:
		return "nil"
	case TokenCase:
		return "case"
	case TokenWhen:
		return "when"
	case TokenIn:
		return "in"
	case TokenAction:
		return "action"
	case TokenJudge:
		return "judge"
	case TokenAgent:
		return "agent"
	case TokenEvent:
		return "event"
	case TokenEmit:
		return "emit"
	case TokenActionClass:
		return "action_class"
	case TokenErhProfile:
		return "erh_profile"
	case TokenScenario:
		return "scenario"
	case TokenCognitiveSiem:
		return "cognitive_siem"
	case TokenFrom:
		return "from"
	case TokenDataset:
		return "dataset"
	case TokenPrimeThreshold:
		return "prime_threshold"
	case TokenFitAlpha:
		return "fit_alpha"
	case TokenPlus:
		return "+"
	case TokenMinus:
		return "-"
	case TokenStar:
		return "*"
	case TokenSlash:
		return "/"
	case TokenMod:
		return "%"
	case TokenEq:
		return "=="
	case TokenNeq:
		return "!="
	case TokenLt:
		return "<"
	case TokenGt:
		return ">"
	case TokenLe:
		return "<="
	case TokenGe:
		return ">="
	case TokenAnd:
		return "&&"
	case TokenOr:
		return "||"
	case TokenNot:
		return "!"
	case TokenAssign:
		return "="
	case TokenArrow:
		return "->"
	case TokenFatArrow:
		return "=>"
	case TokenDot:
		return "."
	case TokenColon:
		return ":"
	case TokenDoubleColon:
		return "::"
	case TokenDotDot:
		return ".."
	case TokenLParen:
		return "("
	case TokenRParen:
		return ")"
	case TokenLBrace:
		return "{"
	case TokenRBrace:
		return "}"
	case TokenLBracket:
		return "["
	case TokenRBracket:
		return "]"
	case TokenComma:
		return ","
	case TokenSemicolon:
		return ";"
	default:
		return "UNKNOWN"
	}
}

// String returns the string representation of the token
func (t Token) String() string {
	if t.Lexeme != "" {
		return t.Type.String() + "(" + t.Lexeme + ")"
	}
	return t.Type.String()
}

// keywords maps keyword strings to token types
var keywords = map[string]TokenType{
	"fn":               TokenFn,
	"return":           TokenReturn,
	"if":               TokenIf,
	"else":             TokenElse,
	"match":            TokenMatch,
	"while":            TokenWhile,
	"for":              TokenFor,
	"break":            TokenBreak,
	"continue":         TokenContinue,
	"let":              TokenLet,
	"struct":           TokenStruct,
	"enum":             TokenEnum,
	"type":             TokenTypeDecl,
	"import":           TokenImport,
	"module":           TokenModule,
	"pub":              TokenPub,
	"true":             TokenTrue,
	"false":            TokenFalse,
	"nil":              TokenNil,
	"case":             TokenCase,
	"when":             TokenWhen,
	"in":               TokenIn,
	"action":           TokenAction,
	"judge":            TokenJudge,
	"agent":            TokenAgent,
	"event":            TokenEvent,
	"emit":             TokenEmit,
	"action_class":     TokenActionClass,
	"erh_profile":      TokenErhProfile,
	"scenario":         TokenScenario,
	"cognitive_siem":   TokenCognitiveSiem,
	"from":             TokenFrom,
	"dataset":          TokenDataset,
	"prime_threshold":  TokenPrimeThreshold,
	"fit_alpha":        TokenFitAlpha,
}

// LookupKeyword checks if an identifier is a keyword
func LookupKeyword(ident string) (TokenType, bool) {
	if typ, ok := keywords[ident]; ok {
		return typ, true
	}
	return TokenIdent, false
}

// IsKeyword checks if a token type is a keyword
func IsKeyword(typ TokenType) bool {
	return typ >= TokenFn && typ <= TokenFitAlpha
}

// IsOperator checks if a token type is an operator
func IsOperator(typ TokenType) bool {
	return (typ >= TokenPlus && typ <= TokenMod) ||
		(typ >= TokenEq && typ <= TokenGe) ||
		(typ >= TokenAnd && typ <= TokenNot) ||
		typ == TokenAssign
}

// IsLiteral checks if a token type is a literal
func IsLiteral(typ TokenType) bool {
	return typ == TokenInt || typ == TokenFloat || typ == TokenString ||
		typ == TokenTrue || typ == TokenFalse || typ == TokenDuration
}

