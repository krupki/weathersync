# Contributing to WeatherSync

Thank you for your interest in contributing to WeatherSync! We welcome contributions from everyone, regardless of skill level. This document will guide you through the contribution process.

## Code of Conduct

Please note that this project is released with a [Contributor Code of Conduct](CODE_OF_CONDUCT.md). By participating in this project, you agree to abide by its terms.

---

## Getting Started

### Prerequisites

- Go 1.18 or later
- Git
- A GitHub account

### Setting Up Your Development Environment

1. **Fork the repository**

   ```bash
   # Click "Fork" on the GitHub repository
   ```

2. **Clone your fork**

   ```bash
   git clone https://github.com/YOUR-USERNAME/weathersync.git
   cd weathersync
   ```

3. **Add upstream remote**

   ```bash
   git remote add upstream https://github.com/krupki/weathersync.git
   ```

4. **Create a development branch**

   ```bash
   git checkout -b feature/your-feature-name
   # or
   git checkout -b fix/your-bug-fix
   ```

5. **Install dependencies**

   ```bash
   cd weather-compare
   go mod download
   ```

---

## Making Changes

### Code Style

- Follow [Effective Go](https://golang.org/doc/effective_go) guidelines
- Run `go fmt` on all Go files before committing
- Use `golangci-lint` for linting:

  ```bash
  golangci-lint run ./...
  ```

### Commit Messages

Write clear, descriptive commit messages:

```text
feat: add caching layer for API responses
fix: handle nil pointer in weather comparison
docs: update README with new examples
test: add unit tests for fetcher package
```

Use these prefixes:

- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation
- `test:` - Tests
- `refactor:` - Code refactoring
- `perf:` - Performance improvements
- `chore:` - Maintenance

### Testing

Before submitting a PR, ensure all tests pass:

```bash
go test ./...
```

Write tests for new features:

```bash
# Add tests to the same package
# Example: internal/fetcher/fetcher_test.go
```

### Documentation

- Add doc comments to all public functions:

  ```go
  // FetchWeatherData fetches weather information for a given city.
  // It returns a WeatherResult with current temperature and fetch duration.
  func FetchWeatherData(city string, lat, lon float64) (*WeatherResult, error) {
  ```

- Update the README if you add new features
- Add examples for complex functionality

---

## Submitting a Pull Request

1. **Update your branch**

   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **Push to your fork**

   ```bash
   git push origin feature/your-feature-name
   ```

3. **Create a Pull Request**
   - Go to GitHub and click "New Pull Request"
   - Fill in the PR template with:
     - What does this PR do?
     - Why are these changes needed?
     - Any related issues? (e.g., Closes #123)
     - Screenshots/examples if applicable

4. **PR Description Template**

   ```markdown
   ## Description
   Brief description of the changes.

   ## Type of Change
   - [ ] Bug fix
   - [ ] New feature
   - [ ] Breaking change
   - [ ] Documentation update

   ## Testing
   How was this tested?

   ## Checklist
   - [ ] Code follows style guidelines
   - [ ] Tests added/updated
   - [ ] Documentation updated
   - [ ] No breaking changes
   ```

---

## Review Process

- A maintainer will review your PR
- Be open to feedback and suggestions
- Make requested changes in new commits (don't force push)
- Once approved, your PR will be merged!

---

## Areas for Contribution

### Easy (Great for Beginners)

- [ ] Add more cities to `config/cities.yaml`
- [ ] Improve documentation
- [ ] Fix typos in comments/docs
- [ ] Add examples

### Medium

- [ ] Add unit tests for existing functions
- [ ] Improve error handling
- [ ] Add configuration validation
- [ ] Create a Makefile for common tasks

### Advanced

- [ ] Implement retry logic with exponential backoff
- [ ] Add database storage layer
- [ ] Create REST API endpoint
- [ ] Implement caching mechanism
- [ ] Add Prometheus metrics export

---

## Questions?

- Open an issue for bugs or feature requests
- Use GitHub Discussions for questions
- Check existing issues before opening a new one

---

## License

By contributing to WeatherSync, you agree that your contributions will be licensed under its MIT License.

Thank you for contributing! 🚀
