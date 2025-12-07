package lsp

import (
	"github.com/dennislee928/uad-lang/internal/lexer"
	"github.com/dennislee928/uad-lang/internal/parser"
	"github.com/dennislee928/uad-lang/internal/typer"
)

// Analyzer performs code analysis
type Analyzer struct {
	// Configuration
	config AnalyzerConfig
}

// AnalyzerConfig holds analyzer configuration
type AnalyzerConfig struct {
	EnableTypeChecking bool
	EnableLinting      bool
}

// NewAnalyzer creates a new analyzer
func NewAnalyzer() *Analyzer {
	return &Analyzer{
		config: AnalyzerConfig{
			EnableTypeChecking: true,
			EnableLinting:      true,
		},
	}
}

// Analyze performs full analysis on a document
func (a *Analyzer) Analyze(doc *Document) []Diagnostic {
	diagnostics := []Diagnostic{}
	
	// Tokenize
	l := lexer.New(doc.Content, doc.URI)
	tokens := []lexer.Token{}
	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == lexer.TokenEOF {
			break
		}
	}
	
	// Parse
	p := parser.New(tokens, doc.URI)
	module, err := p.ParseModule()
	if err != nil {
		// Parser error - type: *common.ErrorList
		// For now, create a simple diagnostic
		diagnostics = append(diagnostics, Diagnostic{
			Range: Range{
				Start: Position{Line: 0, Character: 0},
				End:   Position{Line: 0, Character: 1},
			},
			Severity: DiagnosticSeverityError,
			Source:   "parser",
			Message:  err.Error(),
		})
		
		return diagnostics
	}
	
	// Cache AST
	doc.AST = module
	
	// Type check
	if a.config.EnableTypeChecking {
		tc := typer.NewTypeChecker()
		if err := tc.Check(module); err != nil {
			// Type checking errors - for now, simple error message
			diagnostics = append(diagnostics, Diagnostic{
				Range: Range{
					Start: Position{Line: 0, Character: 0},
					End:   Position{Line: 0, Character: 1},
				},
				Severity: DiagnosticSeverityError,
				Source:   "type-checker",
				Message:  err.Error(),
			})
		}
	}
	
	// TODO: Additional linting checks
	// TODO: Extract precise error positions from common.ErrorList
	
	return diagnostics
}

// Diagnostic represents a diagnostic message
type Diagnostic struct {
	Range    Range
	Severity DiagnosticSeverity
	Code     string
	Source   string
	Message  string
}

// DiagnosticSeverity levels
type DiagnosticSeverity int

const (
	DiagnosticSeverityError       DiagnosticSeverity = 1
	DiagnosticSeverityWarning     DiagnosticSeverity = 2
	DiagnosticSeverityInformation DiagnosticSeverity = 3
	DiagnosticSeverityHint        DiagnosticSeverity = 4
)

// Range in a document
type Range struct {
	Start Position
	End   Position
}

// Position in a document
type Position struct {
	Line      int
	Character int
}

