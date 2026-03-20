<!--
Copyright (c) 2026 dox authors.

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

.github/CONTRIBUTING.md

- Author   : Frost Leo <frostleo.dev@gmail.com>
- Created  : 2026-03-20
- Modified : 2026-03-20
-->

# Contributing to dox

Thank you for your interest in contributing to dox. This document provides guidelines and information for contributors.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Commit Guidelines](#commit-guidelines)
- [Pull Request Process](#pull-request-process)
- [Code Style](#code-style)
- [Testing](#testing)

## Code of Conduct

By participating in this project, you agree to abide by our [Code of Conduct](CODE_OF_CONDUCT.md). Please read it before contributing.

## Getting Started

1. Fork the repository
2. Clone your fork locally
3. Set up the development environment
4. Create a new branch for your changes
5. Make your changes
6. Submit a pull request

## Development Setup

### Prerequisites

- Go 1.22 or later
- PostgreSQL 15 or later
- Node.js 20 or later (for frontend)
- Docker (optional, for containerized development)

### Installation

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/dox.git
cd dox

# Add upstream remote
git remote add upstream https://github.com/opendox/dox.git

# Install dependencies
go mod download

# Copy environment configuration
cp .env.example .env

# Run database migrations
make migrate

# Start development server
make dev
```

## Making Changes

### Branch Naming

Use descriptive branch names with the following prefixes:

- `feat/` - New features
- `fix/` - Bug fixes
- `docs/` - Documentation changes
- `refactor/` - Code refactoring
- `test/` - Test additions or modifications
- `chore/` - Maintenance tasks

Example: `feat/add-product-analytics`

### Before Submitting

1. Ensure your code compiles without errors
2. Run the test suite
3. Update documentation if needed
4. Rebase your branch on the latest main

## Commit Guidelines

We follow a structured commit message format that links commits to GitHub issues and PRs.

### Commit Message Format

```
gh-<issue> <type>: <subject> (#<pr>)

<body>

Date: <YYYY-MM-DD>
Closes #<issue>
Signed-off-by: Name <email>
```

### Components

- `gh-<issue>` - GitHub issue number being addressed
- `<type>` - Type of change (see below)
- `<subject>` - Brief description of the change
- `(#<pr>)` - Pull request number
- `<body>` - Detailed description of changes (use bullet points for multiple items)
- `Date` - Date of the commit
- `Closes #<issue>` - Links and closes the related issue
- `Signed-off-by` - Developer Certificate of Origin signature

### Types

- `feat` - New feature
- `fix` - Bug fix
- `docs` - Documentation only changes
- `style` - Code style changes (formatting, semicolons, etc.)
- `refactor` - Code refactoring
- `perf` - Performance improvements
- `test` - Adding or updating tests
- `chore` - Maintenance tasks
- `ci` - CI/CD changes

### Examples

```
gh-42 feat: Add product analytics dashboard (#45)

Implement product performance analytics with the following features:

- ProductAnalyticsService: Core analytics calculation engine
- DashboardController: REST endpoints for dashboard data
- ChartComponents: Reusable chart components for visualization

Date: 2026-03-20
Closes #42
Signed-off-by: Your Name <your.email@example.com>
```

```
gh-18 fix: Resolve data sync timeout issue (#21)

Fixed timeout errors during large dataset synchronization by implementing
batch processing with configurable chunk sizes.

Date: 2026-03-19
Closes #18
Signed-off-by: Your Name <your.email@example.com>
```

### Signing Off

All commits must include a `Signed-off-by` line. This certifies that you wrote the code or have the right to submit it under the project's license.

Add it automatically with:

```bash
git commit -s -m "your message"
```

Or configure Git to always sign off:

```bash
git config --global user.name "Your Name"
git config --global user.email "your.email@example.com"
```

## Pull Request Process

1. Ensure your PR addresses a specific issue or clearly describes the problem it solves
2. Update the README.md or relevant documentation if needed
3. Add tests for new functionality
4. Ensure all tests pass
5. Request review from maintainers
6. Address review feedback promptly

### PR Title Format

Follow the same format as commit message subjects:

```
gh-<issue> <type>: <subject>
```

Example: `gh-42 feat: Add product analytics dashboard`

## Code Style

### Go

- Follow [Effective Go](https://golang.org/doc/effective_go) guidelines
- Use `gofmt` for formatting
- Run `golangci-lint` before committing
- Keep functions focused and concise
- Write meaningful comments for exported functions

### Frontend

- Follow the project's ESLint configuration
- Use TypeScript for type safety
- Write meaningful component and function names

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific package tests
go test ./internal/api/...
```

### Writing Tests

- Write unit tests for new functionality
- Include both positive and negative test cases
- Use table-driven tests where appropriate
- Mock external dependencies

## Questions?

If you have questions about contributing, feel free to:

- Open a [Discussion](https://github.com/opendox/dox/discussions)
- Check existing issues and discussions

Thank you for contributing to dox.
