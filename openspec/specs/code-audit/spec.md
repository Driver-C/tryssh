## ADDED Requirements

### Requirement: golangci-lint integration
The project SHALL use `golangci-lint` with a comprehensive linter configuration. The configuration SHALL be stored in `.golangci.yml` at the project root.

#### Scenario: Lint execution
- **WHEN** running `golangci-lint run ./...`
- **THEN** all source files SHALL pass without errors or warnings

#### Scenario: Linter configuration
- **WHEN** `.golangci.yml` is present
- **THEN** it SHALL enable at minimum: errcheck, govet, staticcheck, unused, gosimple, ineffassign, typecheck, gocritic, gosec, revive

### Requirement: Vulnerability scanning
The project SHALL use `govulncheck` to scan for known vulnerabilities in dependencies.

#### Scenario: Vulnerability check
- **WHEN** running `govulncheck ./...`
- **THEN** no known vulnerabilities SHALL be found in the codebase or dependencies

#### Scenario: CI vulnerability gate
- **WHEN** a pull request is submitted
- **THEN** the CI pipeline SHALL run `govulncheck` and fail if vulnerabilities are detected

### Requirement: Security audit for sensitive operations
The following code areas SHALL be manually audited for security issues:
1. SSH key handling and storage in `pkg/launcher/base.go`
2. Password handling in configuration files `pkg/config/`
3. Known hosts management in `pkg/launcher/base.go`
4. File transfer operations in `pkg/launcher/scp.go`

#### Scenario: Password exposure check
- **WHEN** passwords are logged or displayed
- **THEN** they SHALL be masked or redacted

#### Scenario: SSH key file permissions
- **WHEN** SSH key files are read
- **THEN** the system SHALL warn if file permissions are too permissive

#### Scenario: Known hosts validation
- **WHEN** connecting to a server
- **THEN** the system SHALL validate the host key against known_hosts

### Requirement: Error handling audit
All error paths in the codebase SHALL be audited to ensure:
1. Errors are properly propagated (not silently ignored)
2. No `log.Fatalf` is used outside `main()` or command entry points
3. Error messages are descriptive and actionable

#### Scenario: Error propagation
- **WHEN** an error occurs in any package function
- **THEN** the error SHALL be returned to the caller, not logged and exited

#### Scenario: Descriptive error messages
- **WHEN** a configuration error occurs
- **THEN** the error message SHALL include the file path and specific issue

### Requirement: Code quality standards
The codebase SHALL comply with standard Go code quality practices:
1. All exported functions and types SHALL have doc comments
2. No global mutable state
3. Consistent naming following Go conventions
4. No dead code or unreachable paths

#### Scenario: Go vet compliance
- **WHEN** running `go vet ./...`
- **THEN** no issues SHALL be reported

#### Scenario: Exported symbol documentation
- **WHEN** an exported function or type exists
- **THEN** it SHALL have a godoc-compatible comment
