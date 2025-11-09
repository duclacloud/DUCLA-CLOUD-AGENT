# Contributing to Ducla Cloud Agent

Thank you for your interest in contributing to Ducla Cloud Agent! This document provides guidelines and information for contributors.

## 🚀 Getting Started

### Prerequisites

- Go 1.21 or later
- Git
- Linux/macOS/Windows development environment

### Development Setup

1. **Fork the repository**
   ```bash
   # Fork on GitHub, then clone your fork
   git clone https://github.com/YOUR_USERNAME/DUCLA-CLOUD-AGENT.git
   cd DUCLA-CLOUD-AGENT
   ```

2. **Set up upstream remote**
   ```bash
   git remote add upstream https://github.com/duclacloud/DUCLA-CLOUD-AGENT.git
   ```

3. **Install dependencies**
   ```bash
   go mod tidy
   ```

4. **Build and test**
   ```bash
   ./build-v1.sh
   ./demo-auto.sh
   ```

## 📋 How to Contribute

### Reporting Issues

- Use the [GitHub Issues](https://github.com/duclacloud/DUCLA-CLOUD-AGENT/issues) page
- Search existing issues before creating a new one
- Provide detailed information:
  - Operating system and version
  - Go version
  - Steps to reproduce
  - Expected vs actual behavior
  - Relevant logs or error messages

### Suggesting Features

- Open a [GitHub Issue](https://github.com/duclacloud/DUCLA-CLOUD-AGENT/issues) with the "enhancement" label
- Describe the feature and its use case
- Explain why it would be valuable
- Consider implementation complexity

### Code Contributions

1. **Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes**
   - Follow the coding standards (see below)
   - Add tests for new functionality
   - Update documentation as needed

3. **Test your changes**
   ```bash
   go test ./...
   ./build-v1.sh
   ./demo-auto.sh
   ```

4. **Commit your changes**
   ```bash
   git add .
   git commit -m "feat: add your feature description"
   ```

5. **Push and create a Pull Request**
   ```bash
   git push origin feature/your-feature-name
   ```

## 🎯 Coding Standards

### Go Code Style

- Follow standard Go formatting: `go fmt`
- Use meaningful variable and function names
- Add comments for exported functions and types
- Keep functions focused and small
- Handle errors appropriately

### Commit Messages

Use conventional commit format:
- `feat:` for new features
- `fix:` for bug fixes
- `docs:` for documentation changes
- `refactor:` for code refactoring
- `test:` for adding tests
- `chore:` for maintenance tasks

Examples:
```
feat: add CLI command for task management
fix: resolve memory leak in file operations
docs: update API documentation
```

### Code Organization

- Keep related functionality in appropriate packages
- Use interfaces for testability
- Follow the existing project structure
- Add unit tests for new code

## 🧪 Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/agent
```

### Writing Tests

- Write unit tests for new functions
- Use table-driven tests when appropriate
- Mock external dependencies
- Test both success and error cases

### Integration Testing

```bash
# Build and run demo
./build-v1.sh
./demo-auto.sh

# Test CLI commands
./demo-cli.sh

# Test packages
./build-packages-simple.sh
```

## 📚 Documentation

### Code Documentation

- Add godoc comments for exported functions
- Include usage examples in comments
- Document complex algorithms or business logic

### User Documentation

- Update README-VI.md for user-facing changes
- Update man page (docs/ducla-agent.1) for CLI changes
- Add examples to WORKSHOP.md for new features

## 🔍 Code Review Process

### Pull Request Guidelines

- Provide a clear description of changes
- Reference related issues
- Include screenshots for UI changes
- Ensure all tests pass
- Update documentation as needed

### Review Criteria

- Code quality and style
- Test coverage
- Documentation completeness
- Backward compatibility
- Performance impact

## 🏗️ Development Workflow

### Branch Strategy

- `main`: Stable release branch
- `develop`: Development integration branch
- `feature/*`: Feature development branches
- `hotfix/*`: Critical bug fix branches

### Release Process

1. Features merged to `develop`
2. Release candidate created from `develop`
3. Testing and bug fixes
4. Merge to `main` and tag release
5. Update documentation and changelog

## 🤝 Community Guidelines

### Code of Conduct

- Be respectful and inclusive
- Welcome newcomers and help them learn
- Focus on constructive feedback
- Respect different opinions and approaches

### Communication

- Use GitHub Issues for bug reports and feature requests
- Use GitHub Discussions for general questions
- Be patient and helpful in responses
- Provide context and examples

## 🛠️ Development Tips

### Useful Commands

```bash
# Format code
go fmt ./...

# Lint code (if golangci-lint is installed)
golangci-lint run

# Build for multiple platforms
GOOS=linux GOARCH=amd64 go build ./cmd/agent
GOOS=windows GOARCH=amd64 go build ./cmd/agent

# Generate documentation
go doc ./internal/agent
```

### Debugging

- Use `go run` for quick testing
- Add debug logging with `-debug` flag
- Use `go tool pprof` for performance profiling
- Test with different configurations

## 📊 Project Structure

```
DUCLA-CLOUD-AGENT/
├── cmd/agent/          # Main application
├── internal/           # Internal packages
│   ├── agent/         # Core agent logic
│   ├── api/           # HTTP/gRPC APIs
│   ├── config/        # Configuration
│   ├── executor/      # Task execution
│   ├── fileops/       # File operations
│   ├── health/        # Health checks
│   ├── metrics/       # Metrics collection
│   └── transport/     # Network transport
├── docs/              # Documentation
├── scripts/           # Build scripts
└── configs/           # Configuration examples
```

## 🎉 Recognition

Contributors will be recognized in:
- CHANGELOG.md for their contributions
- GitHub contributors page
- Release notes for significant contributions

## 📞 Getting Help

- Check existing documentation first
- Search GitHub Issues for similar problems
- Create a new issue with detailed information
- Join community discussions

Thank you for contributing to Ducla Cloud Agent! 🚀