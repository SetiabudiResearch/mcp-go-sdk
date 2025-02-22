# Contributing

Thank you for your interest in contributing to the MCP Go SDK! This document provides guidelines and instructions for contributing.

## Development Setup

1. Make sure you have Go 1.23+ installed
2. Fork the repository
3. Clone your fork: `git clone https://github.com/YOUR-USERNAME/golang-sdk.git`
4. Install dependencies:
```bash
go mod download
```

## Development Workflow

1. Choose the correct branch for your changes:
   - For bug fixes to a released version: use the latest release branch (e.g. v1.1.x for 1.1.3)
   - For new features: use the main branch (which will become the next minor/major version)
   - If unsure, ask in an issue first

2. Create a new branch from your chosen base branch

3. Make your changes

4. Ensure tests pass:
```bash 
go test ./...
```

5. Run linting:
```bash
golangci-lint run
```

6. Format code:
```bash
go fmt ./...
```

7. Submit a pull request to the same branch you branched from

## Code Style

- We use `golangci-lint` for linting
- Follow standard Go style guidelines from [Effective Go](https://golang.org/doc/effective_go)
- Add comments for exported functions and types following Go conventions
- Use meaningful variable names and keep functions focused and concise

## Pull Request Process

1. Update documentation as needed
2. Add tests for new functionality
3. Ensure CI passes
4. Maintainers will review your code
5. Address review feedback

## Code of Conduct

Please note that this project is released with a [Code of Conduct](CODE_OF_CONDUCT.md). By participating in this project you agree to abide by its terms.

## License

By contributing, you agree that your contributions will be licensed under the MIT License. 