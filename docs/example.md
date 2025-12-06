cmd/uadc/main.go（compiler CLI 骨架）
package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	inputPath  string
	outputPath string
	mode       string // "core", "model"
)

func init() {
	flag.StringVar(&inputPath, "i", "", "input .uad or .uadmodel file")
	flag.StringVar(&outputPath, "o", "out.uadir", "output .uad-IR file")
	flag.StringVar(&mode, "mode", "auto", "parse mode: auto|core|model")
}

func main() {
	flag.Parse()

	if inputPath == "" {
		fmt.Fprintln(os.Stderr, "uadc: missing -i <input> file")
		os.Exit(1)
	}

	if err := compile(inputPath, outputPath, mode); err != nil {
		fmt.Fprintf(os.Stderr, "uadc: %v\n", err)
		os.Exit(1)
	}
}

func compile(inPath, outPath, mode string) error {
	// TODO:
	// 1. 讀檔
	// 2. 判斷是 .uad-core 還是 .uad-model (auto 模式時看副檔名或檔頭)
	// 3. lexer → parser → AST
	// 4. 若是 model：model.Desugar(AST_model) → AST_core
	// 5. typer.TypeCheck(AST_core)
	// 6. ir.Build(AST_core) → IR module
	// 7. ir.Encode(IR, outPath)
	return fmt.Errorf("compile not implemented yet")
}

2.9 cmd/uadvm/main.go（VM runner 骨架）
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourname/uad-lang/internal/ir"
	"github.com/yourname/uad-lang/internal/vm"
)

var irPath string

func init() {
	flag.StringVar(&irPath, "i", "", "input .uad-IR file")
}

func main() {
	flag.Parse()

	if irPath == "" {
		fmt.Fprintln(os.Stderr, "uadvm: missing -i <input> .uad-IR file")
		os.Exit(1)
	}

	mod, err := ir.LoadFromFile(irPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "uadvm: failed to load IR: %v\n", err)
		os.Exit(1)
	}

	v := vm.New()
	if err := v.Run(mod); err != nil {
		fmt.Fprintf(os.Stderr, "uadvm: execution error: %v\n", err)
		os.Exit(1)
	}
}

2.10 internal/common/position.go
package common

type Position struct {
	Line   int
	Column int
}

type Span struct {
	Start Position
	End   Position
}

2.11 internal/lexer/tokens.go
package lexer

import "github.com/yourname/uad-lang/internal/common"

type TokenType int

const (
	// Special
	TokenEOF TokenType = iota
	TokenIllegal

	// Identifiers & literals
	TokenIdent
	TokenInt
	TokenFloat
	TokenString

	// Keywords
	TokenFn
	TokenStruct
	TokenEnum
	TokenLet
	TokenIf
	TokenElse
	TokenTrue
	TokenFalse
	// ... 之後補齊

	// Operators & delimiters
	TokenPlus
	TokenMinus
	TokenStar
	TokenSlash
	TokenLParen
	TokenRParen
	TokenLBrace
	TokenRBrace
	TokenComma
	TokenColon
	TokenSemicolon
	TokenArrow // ->
)

type Token struct {
	Type    TokenType
	Lexeme  string
	Span    common.Span
}

2.12 internal/lexer/lexer.go（掃描器骨架）
package lexer

import (
	"unicode"
	"unicode/utf8"

	"github.com/yourname/uad-lang/internal/common"
)

type Lexer struct {
	src    string
	offset int
	line   int
	col    int
}

func New(src string) *Lexer {
	return &Lexer{
		src:  src,
		line: 1,
		col:  1,
	}
}

func (l *Lexer) NextToken() Token {
	// TODO: skip whitespace & comments
	// read rune by rune, build tokens
	r, size := l.peek()
	if size == 0 {
		return l.token(TokenEOF, "")
	}

	if isLetter(r) || r == '_' {
		return l.lexIdentOrKeyword()
	}

	if unicode.IsDigit(r) {
		return l.lexNumber()
	}

	// TODO: operators, strings, etc.
	return l.token(TokenIllegal, string(r))
}

func (l *Lexer) peek() (rune, int) {
	if l.offset >= len(l.src) {
		return 0, 0
	}
	r, size := utf8.DecodeRuneInString(l.src[l.offset:])
	return r, size
}

func (l *Lexer) advance(size int) {
	if size == 0 {
		return
	}
	for i := 0; i < size; i++ {
		if l.src[l.offset+i] == '\n' {
			l.line++
			l.col = 1
		} else {
			l.col++
		}
	}
	l.offset += size
}

func (l *Lexer) token(t TokenType, lexeme string) Token {
	// NOTE: 這裡可以更準確地記錄 Span，先簡化
	return Token{
		Type:   t,
		Lexeme: lexeme,
		Span: common.Span{
			Start: common.Position{Line: l.line, Column: l.col},
			End:   common.Position{Line: l.line, Column: l.col + len(lexeme)},
		},
	}
}

func isLetter(r rune) bool {
	return unicode.IsLetter(r)
}


（後續你可在 Cursor 裡補完 lexIdentOrKeyword、lexNumber…）

2.13 internal/ast/core_nodes.go（.uad-core AST）
package ast

import "github.com/yourname/uad-lang/internal/common"

type Node interface {
	Span() common.Span
}

type Expr interface {
	Node
	exprNode()
}

type Stmt interface {
	Node
	stmtNode()
}

type Ident struct {
	Name string
	Loc  common.Span
}

func (i *Ident) Span() common.Span { return i.Loc }
func (i *Ident) exprNode()         {}

type LiteralKind int

const (
	LitInt LiteralKind = iota
	LitFloat
	LitString
	LitBool
)

type Literal struct {
	Kind LiteralKind
	Text string
	Loc  common.Span
}

func (l *Literal) Span() common.Span { return l.Loc }
func (l *Literal) exprNode()         {}

type BinaryOp int

const (
	OpAdd BinaryOp = iota
	OpSub
	OpMul
	OpDiv
	OpEq
	OpNeq
	OpLt
	OpGt
)

type BinaryExpr struct {
	Op    BinaryOp
	Left  Expr
	Right Expr
	Loc   common.Span
}

func (b *BinaryExpr) Span() common.Span { return b.Loc }
func (b *BinaryExpr) exprNode()         {}

// TODO: CallExpr, IfExpr, WhileStmt, FnDecl, StructDecl, EnumDecl, Module 等

2.14 examples/core/hello_world.uad
fn main() {
  print("Hello, .uad core!")
}


（你可以先在 VM runtime 裡硬寫一個 print builtin 讓這個跑起來）

2.15 examples/model/devsecops_erh.uadmodel
action_class MergeRequest {
  complexity = log(1 + lines_changed) + 0.5 * files_changed
  true_value = if has_incident_within(90d) then -1.0 else +1.0
  importance = asset_criticality * (1 + internet_exposed)
}

judge pipeline_judge for MergeRequest {
  decision = case {
    pipeline_passed && !overridden -> +1.0
    !pipeline_passed               -> -1.0
    else                           -> 0.0
  }
}

erh_profile "GitLab-DevSecOps" {
  actions from dataset "mr_security_logs"
  judge   = pipeline_judge

  prime_threshold {
    mistake_delta > 0.5
    importance_quantile >= 0.9
    complexity >= 40.0
  }

  fit_alpha {
    range  = [10.0, 80.0]
    method = "loglog_regression"
  }
}

3. 建議實作順序（很重要）

以你現在的工作量與手上的 ERH/Security 專案，我會建議這樣排：

Phase 1 – 跑得動的 .uad-core + Interpreter

完成：lexer（subset）、ast、parser/core_parser.go（支援 fn / let / basic expr）、
簡單的 interpreter（你可以先在 internal/vm 實作 AST interpreter，而不是 IR VM）。

目標：examples/core/hello_world.uad 可以印字串；
再做一個簡單 erh_basic.uad 做 Σ/平均。

Phase 2 – IR + VM

把 interpreter 往 IR 移：
AST → IR → IR VM 執行。

先只支援 Int/Float/Bool + 基本算術 / if / while。

Phase 3 – .uad-model desugar

實作 internal/model/desugar.go，
把 .uadmodel 的 action_class / judge / erh_profile 轉成 .uad-core 程式。

這時你就可以用自己設計的語言描述 ERH profile，然後實際跑出 α 與 primes。

Phase 4 – Security / Cyber Range / SIEM 擴展

把你的 ERH_ON_SECURITY_POC 裡的 schema 映射到 .uad-model。

在 runtime/runtime/stdlib/ 補 erh.uad, security.uad 等 library。
