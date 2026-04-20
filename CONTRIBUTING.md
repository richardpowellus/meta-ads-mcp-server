# Contributing to meta-ads-mcp-server

Thank you for your interest in contributing! This guide will help you get started.

## Building from Source

```bash
git clone https://github.com/richardpowellus/meta-ads-mcp-server.git
cd meta-ads-mcp-server
go build ./cmd/meta-ads-mcp-server
```

Requirements: Go 1.26+

## Running Tests

```bash
go test ./...
```

## Code Style

- Format all code with `gofmt` (or `goimports`)
- Pass `go vet ./...` with no warnings
- Pass `golangci-lint run` with no errors

## Adding a New Tool

1. Create a new file in the appropriate `tools/` package (e.g. `metaads/tools/my_feature.go`)
2. Register the tool in the `RegisterAll()` function in `register.go`
3. Follow the existing tool pattern:
   - Accept `account` as a parameter
   - Parse input from `json.RawMessage`
   - Use the shared Meta Ads client for API calls
   - Use `paging.Emit()` for collection responses
4. Add tests for the new tool
5. Document the tool in the README

## Pull Request Process

1. Fork the repository and create a feature branch from `main`
2. Make your changes and ensure all tests pass
3. Run `go vet ./...` and `golangci-lint run`
4. Write a clear PR description explaining what and why
5. Link any related issues
6. PRs require at least one approval before merging

## Reporting Bugs

Use the [bug report template](https://github.com/richardpowellus/meta-ads-mcp-server/issues/new?template=bug_report.yml) to file issues.
