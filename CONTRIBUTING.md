# Contributing to tryssh

## Development Setup

1. Install Go 1.25 or later
2. Install golangci-lint: `go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest`
3. Fork and clone the repository
4. Create a feature branch from `dev`

## Development Workflow

1. Make your changes
2. Run tests: `make test`
3. Run linter: `make lint`
4. Ensure all tests pass and coverage is maintained
5. Commit with descriptive messages
6. Push and create a PR targeting `dev`

## Code Standards

- Follow standard Go formatting (`gofmt`)
- All exported symbols must have godoc comments
- Error handling: return errors, don't use `log.Fatalf` outside `cmd/`
- No global mutable state; use dependency injection
- Write unit tests for all new code

## Commit Messages

Use concise, descriptive commit messages. Prefix with type:

- `Add:` New features
- `Fix:` Bug fixes
- `Upd:` Updates and improvements
- `Refactor:` Code restructuring
- `Test:` Test additions or changes
- `Docs:` Documentation updates

## Pull Request Process

1. Ensure CI passes (lint + test)
2. Maintain test coverage
3. Update documentation if needed
4. Request review from maintainers
