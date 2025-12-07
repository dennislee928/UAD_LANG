# Contributing to UAD Programming Language

Thank you for your interest in contributing to UAD! This document provides guidelines and instructions for contributors.

## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Coding Standards](#coding-standards)
- [Testing Guidelines](#testing-guidelines)
- [Documentation](#documentation)
- [Pull Request Process](#pull-request-process)
- [Release Process](#release-process)

---

## 🤝 Code of Conduct

### Our Pledge

We are committed to providing a welcoming and inclusive environment for all contributors, regardless of background or identity.

### Expected Behavior

- Be respectful and constructive
- Welcome newcomers and help them get started
- Focus on what is best for the project
- Show empathy towards other contributors

### Unacceptable Behavior

- Harassment, discrimination, or personal attacks
- Trolling or inflammatory comments
- Publishing others' private information
- Any conduct that would be inappropriate in a professional setting

---

## 🚀 Getting Started

### Prerequisites

- **Go 1.21 or later**
- **Git**
- **Make** (for build automation)
- **golangci-lint** (for linting)

### Setting Up Development Environment

1. **Fork and Clone**
   ```bash
   git clone https://github.com/YOUR_USERNAME/UAD_Programming.git
   cd UAD_Programming
   ```

2. **Install Dependencies**
   ```bash
   make deps
   ```

3. **Build Project**
   ```bash
   make build
   ```

4. **Run Tests**
   ```bash
   make test
   ```

5. **Dev Container (Optional)**
   - Open in VS Code with Dev Containers extension
   - Or use GitHub Codespaces

---

## 🔄 Development Workflow

### Branching Strategy

- `main`: Stable production branch
- `dev`: Development branch (default for PRs)
- `feature/*`: Feature branches
- `fix/*`: Bug fix branches
- `docs/*`: Documentation branches

### Workflow Steps

1. **Create a Branch**
   ```bash
   git checkout dev
   git pull origin dev
   git checkout -b feature/your-feature-name
   ```

2. **Make Changes**
   - Write code
   - Add tests
   - Update documentation

3. **Commit Changes**
   ```bash
   git add .
   git commit -m "feat: add new feature"
   ```
   
   **Commit Message Format** (following Conventional Commits):
   - `feat: ` - New feature
   - `fix: ` - Bug fix
   - `docs: ` - Documentation changes
   - `test: ` - Test additions/changes
   - `refactor: ` - Code refactoring
   - `chore: ` - Maintenance tasks

4. **Push and Create PR**
   ```bash
   git push origin feature/your-feature-name
   ```
   Then create a Pull Request on GitHub targeting the `dev` branch.

---

## 📝 Coding Standards

### Go Style Guide

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting (automated in CI)
- Follow Go naming conventions
  - Exported: `MyFunction`, `MyStruct`
  - Unexported: `myFunction`, `myStruct`

### Code Organization

```
internal/          # Internal implementation
├── ast/          # AST node definitions
├── lexer/        # Tokenization
├── parser/       # Syntax analysis
├── typer/        # Type checking
├── interpreter/  # Direct execution
├── vm/           # Bytecode execution
└── runtime/      # Runtime support
```

### Comments

- **All comments must be in English**
- Public APIs require documentation comments
- Complex algorithms need explanatory comments
- TODOs should include context: `// TODO(M2.3): Implement motif variations`

Example:
```go
// NewTemporalGrid creates a new temporal grid with default settings.
// The grid uses 480 ticks per beat (standard MIDI resolution) and
// a 4/4 time signature at 120 BPM.
func NewTemporalGrid() *TemporalGrid {
    // ...
}
```

### Error Handling

- Return errors explicitly (Go style)
- Provide context in error messages
- Use error wrapping when appropriate

```go
if err := someOperation(); err != nil {
    return fmt.Errorf("failed to perform operation: %w", err)
}
```

---

## 🧪 Testing Guidelines

### Test Structure

- Unit tests in `*_test.go` files alongside source
- Integration tests in `tests/` directory
- Table-driven tests for multiple cases

### Writing Tests

```go
func TestFeature(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
        wantErr  bool
    }{
        {"valid input", "foo", "FOO", false},
        {"empty input", "", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := MyFunction(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("got error %v, wantErr %v", err, tt.wantErr)
            }
            if result != tt.expected {
                t.Errorf("got %v, want %v", result, tt.expected)
            }
        })
    }
}
```

### Coverage Goals

- Core modules: **> 70% coverage**
- New features: **> 80% coverage**
- Critical paths: **100% coverage**

### Running Tests

```bash
# Run all tests
make test

# Run specific package tests
go test ./internal/lexer/...

# Run with coverage
make test-coverage

# Run specific test
go test -run TestFeature ./internal/parser/
```

---

## 📚 Documentation

### Required Documentation

For any new feature, you must provide:

1. **Code Comments**: Explain what and why
2. **README Updates**: If adding user-facing features
3. **Examples**: Add `.uad` example programs
4. **Specification**: Update relevant spec documents if changing language

### Documentation Style

- Use Markdown for all documentation
- Include code examples
- Provide both English and 繁體中文 when appropriate
- Keep line length under 100 characters

---

## 🔀 Pull Request Process

### Before Submitting

- [ ] All tests pass (`make test`)
- [ ] Code is formatted (`make fmt`)
- [ ] Linter passes (`make lint`)
- [ ] Documentation is updated
- [ ] Commit messages follow conventions
- [ ] Branch is up to date with `dev`

### PR Template

When creating a PR, include:

**Description:**
- What does this PR do?
- Why is this change needed?

**Type of Change:**
- [ ] Bug fix
- [ ] New feature
- [ ] Documentation update
- [ ] Refactoring
- [ ] Performance improvement

**Testing:**
- [ ] Unit tests added/updated
- [ ] Integration tests added/updated
- [ ] Manual testing performed

**Checklist:**
- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Comments added for complex code
- [ ] Documentation updated
- [ ] No breaking changes (or documented if necessary)

### Review Process

1. **Automated Checks**: CI must pass
2. **Code Review**: At least one maintainer approval required
3. **Testing**: Verify tests cover new code
4. **Documentation**: Check for completeness
5. **Merge**: Squash and merge to `dev`

### Getting Reviews

- Be patient - reviews may take a few days
- Address feedback constructively
- Ask questions if feedback is unclear
- Update PR based on feedback

---

## 🚢 Release Process

### Version Numbers

We follow [Semantic Versioning](https://semver.org/):
- **Major** (1.0.0): Breaking changes
- **Minor** (0.1.0): New features, backward compatible
- **Patch** (0.0.1): Bug fixes

### Release Steps (Maintainers Only)

1. Ensure `dev` branch is stable
2. Update version numbers
3. Update CHANGELOG.md
4. Create release branch: `release/v1.0.0`
5. Final testing
6. Merge to `main`
7. Tag release: `git tag v1.0.0`
8. Push tags: `git push --tags`
9. Create GitHub release with notes
10. Merge back to `dev`

---

## 🎯 Priority Areas for Contributors

### High Priority

1. **Fix Failing Tests**
   - 3 pre-existing test failures in type checker
   - See `internal/typer/type_checker_test.go`

2. **Complete IR/VM Implementation**
   - `internal/ir/builder.go` (partial)
   - `internal/vm/vm.go` (skeleton only)

3. **Add Example Programs**
   - More `.uad` examples in `examples/core/`
   - Real-world use cases

### Medium Priority

4. **Standard Library**
   - Implement functions in `runtime/stdlib/`
   - Math, string, array operations

5. **Documentation**
   - Tutorials
   - API documentation
   - Language guides

6. **Testing**
   - Integration tests
   - Benchmark suite

### Future/Advanced

7. **LSP Implementation**
   - Language server protocol
   - VS Code extension

8. **Optimization**
   - IR optimization passes
   - Performance improvements

---

## 💡 Ideas and Proposals

### Proposing New Features

1. Open an issue with label `proposal`
2. Describe the feature and motivation
3. Provide examples of usage
4. Discuss implementation approach
5. Wait for maintainer feedback

### RFC (Request for Comments)

For major changes:
1. Write an RFC document
2. Submit as PR to `docs/rfcs/`
3. Community discussion period
4. Approval or rejection by maintainers

---

## ❓ Getting Help

- **Documentation**: Check `docs/` directory
- **Issues**: Search existing issues first
- **Discussions**: Use GitHub Discussions
- **Email**: Contact maintainers (coming soon)

---

## 🙏 Acknowledgments

All contributors will be credited in the project README and release notes.

---

**Happy Contributing! 🎉**

