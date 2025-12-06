package lexer

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/dennislee928/uad-lang/internal/common"
)

// Lexer scans source code and produces tokens
type Lexer struct {
	source string // Source code
	file   string // Source file path
	offset int    // Current byte offset
	line   int    // Current line (1-based)
	column int    // Current column (1-based)
	
	// Lookahead
	ch     rune // Current character
	chSize int  // Size of current character in bytes
}

// New creates a new lexer for the given source code
func New(source, file string) *Lexer {
	l := &Lexer{
		source: source,
		file:   file,
		offset: 0,
		line:   1,
		column: 1,
	}
	l.readChar() // Initialize first character
	return l
}

// NextToken returns the next token from the source
func (l *Lexer) NextToken() Token {
	l.skipWhitespace()
	
	// Save start position
	startPos := l.currentPosition()
	
	// Handle EOF
	if l.ch == 0 {
		return l.makeToken(TokenEOF, "", startPos, l.currentPosition())
	}
	
	// Comments
	if l.ch == '/' && l.peek() == '/' {
		return l.lexLineComment(startPos)
	}
	if l.ch == '/' && l.peek() == '*' {
		return l.lexBlockComment(startPos)
	}
	
	// Identifiers and keywords
	if isLetter(l.ch) || l.ch == '_' {
		return l.lexIdentOrKeyword(startPos)
	}
	
	// Numbers
	if isDigit(l.ch) {
		return l.lexNumber(startPos)
	}
	
	// String literals
	if l.ch == '"' {
		return l.lexString(startPos)
	}
	
	// Operators and delimiters
	tok := l.lexOperator(startPos)
	if tok.Type != TokenIllegal {
		return tok
	}
	
	// Illegal character
	ch := l.ch
	l.readChar()
	return l.makeToken(TokenIllegal, string(ch), startPos, l.currentPosition())
}

// AllTokens returns all tokens from the source
func (l *Lexer) AllTokens() []Token {
	var tokens []Token
	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == TokenEOF {
			break
		}
	}
	return tokens
}

// readChar reads the next character from the source
func (l *Lexer) readChar() {
	if l.offset >= len(l.source) {
		l.ch = 0 // EOF
		l.chSize = 0
		return
	}
	
	r, size := utf8.DecodeRuneInString(l.source[l.offset:])
	l.ch = r
	l.chSize = size
	l.offset += size
	
	// Update line and column
	if l.ch == '\n' {
		l.line++
		l.column = 1
	} else {
		l.column++
	}
}

// peek returns the next character without advancing
func (l *Lexer) peek() rune {
	if l.offset >= len(l.source) {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(l.source[l.offset:])
	return r
}

// peekN returns the character n positions ahead
func (l *Lexer) peekN(n int) rune {
	offset := l.offset
	for i := 0; i < n && offset < len(l.source); i++ {
		_, size := utf8.DecodeRuneInString(l.source[offset:])
		offset += size
	}
	if offset >= len(l.source) {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(l.source[offset:])
	return r
}

// skipWhitespace skips whitespace characters
func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.ch) {
		l.readChar()
	}
}

// lexLineComment lexes a single-line comment
func (l *Lexer) lexLineComment(startPos common.Position) Token {
	l.readChar() // Skip first '/'
	l.readChar() // Skip second '/'
	
	start := l.offset - 2
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
	
	lexeme := l.source[start : l.offset-l.chSize]
	return l.makeToken(TokenComment, lexeme, startPos, l.currentPosition())
}

// lexBlockComment lexes a multi-line comment
func (l *Lexer) lexBlockComment(startPos common.Position) Token {
	l.readChar() // Skip '/'
	l.readChar() // Skip '*'
	
	start := l.offset - 2
	for {
		if l.ch == 0 {
			// Unterminated comment
			return l.makeToken(TokenIllegal, l.source[start:l.offset], startPos, l.currentPosition())
		}
		if l.ch == '*' && l.peek() == '/' {
			l.readChar() // Skip '*'
			l.readChar() // Skip '/'
			break
		}
		l.readChar()
	}
	
	lexeme := l.source[start:l.offset]
	return l.makeToken(TokenComment, lexeme, startPos, l.currentPosition())
}

// lexIdentOrKeyword lexes an identifier or keyword
func (l *Lexer) lexIdentOrKeyword(startPos common.Position) Token {
	start := l.offset - l.chSize
	
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}
	
	lexeme := l.source[start : l.offset-l.chSize]
	
	// Check for keywords
	if typ, isKeyword := LookupKeyword(lexeme); isKeyword {
		return l.makeToken(typ, lexeme, startPos, l.currentPosition())
	}
	
	return l.makeToken(TokenIdent, lexeme, startPos, l.currentPosition())
}

// lexNumber lexes a number (integer or float)
func (l *Lexer) lexNumber(startPos common.Position) Token {
	start := l.offset - l.chSize
	
	// Check for hex (0x) or binary (0b)
	if l.ch == '0' {
		next := l.peek()
		if next == 'x' || next == 'X' {
			return l.lexHexNumber(startPos)
		}
		if next == 'b' || next == 'B' {
			return l.lexBinaryNumber(startPos)
		}
	}
	
	// Read integer part
	for isDigit(l.ch) {
		l.readChar()
	}
	
	// Check for float
	if l.ch == '.' && isDigit(l.peek()) {
		l.readChar() // Skip '.'
		
		// Read fractional part
		for isDigit(l.ch) {
			l.readChar()
		}
		
		// Check for exponent
		if l.ch == 'e' || l.ch == 'E' {
			l.readChar()
			if l.ch == '+' || l.ch == '-' {
				l.readChar()
			}
			for isDigit(l.ch) {
				l.readChar()
			}
		}
		
		lexeme := l.source[start : l.offset-l.chSize]
		return l.makeToken(TokenFloat, lexeme, startPos, l.currentPosition())
	}
	
	// Check for duration suffix (10s, 5m, 2h, 3d)
	if l.ch == 's' || l.ch == 'm' || l.ch == 'h' || l.ch == 'd' {
		l.readChar()
		lexeme := l.source[start:l.offset]
		return l.makeToken(TokenDuration, lexeme, startPos, l.currentPosition())
	}
	
	// Integer
	lexeme := l.source[start : l.offset-l.chSize]
	return l.makeToken(TokenInt, lexeme, startPos, l.currentPosition())
}

// lexHexNumber lexes a hexadecimal number (0xABCD)
func (l *Lexer) lexHexNumber(startPos common.Position) Token {
	start := l.offset - l.chSize
	
	l.readChar() // Skip '0'
	l.readChar() // Skip 'x' or 'X'
	
	if !isHexDigit(l.ch) {
		return l.makeToken(TokenIllegal, l.source[start:l.offset], startPos, l.currentPosition())
	}
	
	for isHexDigit(l.ch) {
		l.readChar()
	}
	
	lexeme := l.source[start : l.offset-l.chSize]
	return l.makeToken(TokenInt, lexeme, startPos, l.currentPosition())
}

// lexBinaryNumber lexes a binary number (0b1010)
func (l *Lexer) lexBinaryNumber(startPos common.Position) Token {
	start := l.offset - l.chSize
	
	l.readChar() // Skip '0'
	l.readChar() // Skip 'b' or 'B'
	
	if l.ch != '0' && l.ch != '1' {
		return l.makeToken(TokenIllegal, l.source[start:l.offset], startPos, l.currentPosition())
	}
	
	for l.ch == '0' || l.ch == '1' {
		l.readChar()
	}
	
	lexeme := l.source[start : l.offset-l.chSize]
	return l.makeToken(TokenInt, lexeme, startPos, l.currentPosition())
}

// lexString lexes a string literal
func (l *Lexer) lexString(startPos common.Position) Token {
	start := l.offset - l.chSize
	l.readChar() // Skip opening quote
	
	var sb strings.Builder
	
	for l.ch != '"' && l.ch != 0 && l.ch != '\n' {
		if l.ch == '\\' {
			l.readChar()
			// Handle escape sequences
			switch l.ch {
			case 'n':
				sb.WriteRune('\n')
			case 't':
				sb.WriteRune('\t')
			case 'r':
				sb.WriteRune('\r')
			case '\\':
				sb.WriteRune('\\')
			case '"':
				sb.WriteRune('"')
			case '0':
				sb.WriteRune('\x00')
			case 'u':
				// Unicode escape: \u{XXXX}
				l.readChar()
				if l.ch != '{' {
					return l.makeToken(TokenIllegal, l.source[start:l.offset], startPos, l.currentPosition())
				}
				l.readChar()
				
				unicodeStr := ""
				for l.ch != '}' && l.ch != 0 {
					unicodeStr += string(l.ch)
					l.readChar()
				}
				
				if l.ch != '}' {
					return l.makeToken(TokenIllegal, l.source[start:l.offset], startPos, l.currentPosition())
				}
				
				// Parse hex value
				var code rune
				for _, ch := range unicodeStr {
					code = code*16 + hexValue(ch)
				}
				sb.WriteRune(code)
			default:
				// Unknown escape sequence
				sb.WriteRune(l.ch)
			}
			l.readChar()
		} else {
			sb.WriteRune(l.ch)
			l.readChar()
		}
	}
	
	if l.ch != '"' {
		// Unterminated string
		return l.makeToken(TokenIllegal, l.source[start:l.offset], startPos, l.currentPosition())
	}
	
	l.readChar() // Skip closing quote
	
	lexeme := l.source[start:l.offset]
	return l.makeToken(TokenString, lexeme, startPos, l.currentPosition())
}

// lexOperator lexes operators and delimiters
func (l *Lexer) lexOperator(startPos common.Position) Token {
	ch := l.ch
	next := l.peek()
	
	switch ch {
	case '+':
		l.readChar()
		return l.makeToken(TokenPlus, "+", startPos, l.currentPosition())
	case '-':
		l.readChar()
		if next == '>' {
			l.readChar()
			return l.makeToken(TokenArrow, "->", startPos, l.currentPosition())
		}
		return l.makeToken(TokenMinus, "-", startPos, l.currentPosition())
	case '*':
		l.readChar()
		return l.makeToken(TokenStar, "*", startPos, l.currentPosition())
	case '/':
		l.readChar()
		return l.makeToken(TokenSlash, "/", startPos, l.currentPosition())
	case '%':
		l.readChar()
		return l.makeToken(TokenMod, "%", startPos, l.currentPosition())
	case '=':
		l.readChar()
		if next == '=' {
			l.readChar()
			return l.makeToken(TokenEq, "==", startPos, l.currentPosition())
		}
		if next == '>' {
			l.readChar()
			return l.makeToken(TokenFatArrow, "=>", startPos, l.currentPosition())
		}
		return l.makeToken(TokenAssign, "=", startPos, l.currentPosition())
	case '!':
		l.readChar()
		if next == '=' {
			l.readChar()
			return l.makeToken(TokenNeq, "!=", startPos, l.currentPosition())
		}
		return l.makeToken(TokenNot, "!", startPos, l.currentPosition())
	case '<':
		l.readChar()
		if next == '=' {
			l.readChar()
			return l.makeToken(TokenLe, "<=", startPos, l.currentPosition())
		}
		return l.makeToken(TokenLt, "<", startPos, l.currentPosition())
	case '>':
		l.readChar()
		if next == '=' {
			l.readChar()
			return l.makeToken(TokenGe, ">=", startPos, l.currentPosition())
		}
		return l.makeToken(TokenGt, ">", startPos, l.currentPosition())
	case '&':
		l.readChar()
		if next == '&' {
			l.readChar()
			return l.makeToken(TokenAnd, "&&", startPos, l.currentPosition())
		}
		return l.makeToken(TokenIllegal, "&", startPos, l.currentPosition())
	case '|':
		l.readChar()
		if next == '|' {
			l.readChar()
			return l.makeToken(TokenOr, "||", startPos, l.currentPosition())
		}
		return l.makeToken(TokenIllegal, "|", startPos, l.currentPosition())
	case '.':
		l.readChar()
		if next == '.' {
			l.readChar()
			return l.makeToken(TokenDotDot, "..", startPos, l.currentPosition())
		}
		return l.makeToken(TokenDot, ".", startPos, l.currentPosition())
	case ':':
		l.readChar()
		if next == ':' {
			l.readChar()
			return l.makeToken(TokenDoubleColon, "::", startPos, l.currentPosition())
		}
		return l.makeToken(TokenColon, ":", startPos, l.currentPosition())
	case '(':
		l.readChar()
		return l.makeToken(TokenLParen, "(", startPos, l.currentPosition())
	case ')':
		l.readChar()
		return l.makeToken(TokenRParen, ")", startPos, l.currentPosition())
	case '{':
		l.readChar()
		return l.makeToken(TokenLBrace, "{", startPos, l.currentPosition())
	case '}':
		l.readChar()
		return l.makeToken(TokenRBrace, "}", startPos, l.currentPosition())
	case '[':
		l.readChar()
		return l.makeToken(TokenLBracket, "[", startPos, l.currentPosition())
	case ']':
		l.readChar()
		return l.makeToken(TokenRBracket, "]", startPos, l.currentPosition())
	case ',':
		l.readChar()
		return l.makeToken(TokenComma, ",", startPos, l.currentPosition())
	case ';':
		l.readChar()
		return l.makeToken(TokenSemicolon, ";", startPos, l.currentPosition())
	default:
		return Token{Type: TokenIllegal}
	}
}

// currentPosition returns the current position in the source
func (l *Lexer) currentPosition() common.Position {
	// Subtract 1 from column because we've already advanced
	col := l.column - 1
	if col < 1 {
		col = 1
	}
	return common.NewPosition(l.line, col, l.offset-l.chSize)
}

// makeToken creates a token with the given type and lexeme
func (l *Lexer) makeToken(typ TokenType, lexeme string, start, end common.Position) Token {
	span := common.NewSpan(start, end, l.file)
	return NewToken(typ, lexeme, span)
}

// Helper functions

// isLetter checks if a rune is a letter
func isLetter(ch rune) bool {
	return unicode.IsLetter(ch)
}

// isDigit checks if a rune is a digit
func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

// isHexDigit checks if a rune is a hexadecimal digit
func isHexDigit(ch rune) bool {
	return (ch >= '0' && ch <= '9') ||
		(ch >= 'a' && ch <= 'f') ||
		(ch >= 'A' && ch <= 'F')
}

// hexValue returns the numeric value of a hex digit
func hexValue(ch rune) rune {
	switch {
	case ch >= '0' && ch <= '9':
		return ch - '0'
	case ch >= 'a' && ch <= 'f':
		return ch - 'a' + 10
	case ch >= 'A' && ch <= 'F':
		return ch - 'A' + 10
	default:
		return 0
	}
}

// Error returns a formatted error message
func (l *Lexer) Error(message string) error {
	pos := l.currentPosition()
	return fmt.Errorf("%s at line %d, column %d", message, pos.Line, pos.Column)
}

