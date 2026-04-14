# Contributing

Thank you for your interest in contributing to imap-spam-cleaner!

Contributions of all kinds are welcome - from fixing bugs and adding features to improving documentation.

---

# Table of Contents

* [Ways to Contribute](#ways-to-contribute)
* [Before You Start](#before-you-start)
* [Development Setup](#development-setup)
* [Code Style & Linting](#code-style--linting)
* [Pull Request Guidelines](#pull-request-guidelines)
* [Commit Message Guidelines](#commit-message-guidelines)
* [Documentation](#documentation)
* [Testing](#testing)
* [Release Workflow](#release-workflow)
* [Security Issues](#security-issues)
* [Code of Conduct](#code-of-conduct)
* [Final Notes](#final-notes)

---

# Ways to Contribute

You can contribute in several ways:

* Fix bugs
* Add new features
* Improve documentation
* Add tests (where applicable)
* Improve performance or reliability

If you're unsure where to start, look for issues labeled:

* `bug`
* `enhancement`

---

# Before You Start

Before beginning work on a new feature or significant change:

1. Open an issue first to show what you will work on.
2. Check for feedback to avoid duplicate work.
3. Keep changes focused and small when possible.

This helps keep development organized and review cycles efficient.

---

# Development Setup

This project follows a standard Go workflow.

## Requirements

* Git
* Go 1.26+

## Clone the Repository

```bash
git clone https://github.com/dominicgisler/imap-spam-cleaner.git
cd imap-spam-cleaner
```

## Install Dependencies

```bash
go get ./...
```

## Run the Application

```bash
go run .
```

---

# Code Style & Linting

This project uses [golangci-lint](https://golangci-lint.run/).

## Run Linter

```bash
make lint
```

## Style Expectations

* Follow idiomatic Go conventions
* Keep code consistent with the existing codebase
* Prefer readability over cleverness
* Keep functions small and focused

Lint checks must pass before submitting a Pull Request.

---

# Pull Request Guidelines

All changes must be submitted via Pull Requests.

## General Rules

* Open an issue before implementing large changes
* Keep PRs small and focused
* One feature or fix per PR
* Multiple features should be split into multiple PRs

## Review Process

* PRs are reviewed by the maintainer
* Linting must pass
* Smaller PRs are reviewed faster

## AI Usage Policy

AI tools may be used as assistance, but:

* Fully AI-generated code submissions should be avoided
* You should understand and review all submitted code
* Code quality and correctness remain your responsibility

---

# Commit Message Guidelines

Write short and descriptive commit messages - prefixes may be used, but are not required.

Examples:

```
fix: handle nil pointer in message parser
feat: add support for multiple folders
docs: update configuration example
```

## Additional Notes

* Small PRs may be rebased
* Large PRs (with many commits) may be squash-merged

---

# Documentation

Documentation is maintained using [MkDocs](https://www.mkdocs.org/).

Location:

```
docs/
```

## Building and Serving Docs Locally

```
make docs
```

This will start a local documentation server so you can preview changes before submitting them.

When adding features or changing behavior:

* Update relevant documentation
* Add examples if useful
* Keep documentation clear and concise

Documentation improvements are always welcome.

## Documentation Deployment

The documentation is automatically built and deployed to [GitHub Pages](https://dominicgisler.github.io/imap-spam-cleaner) when changes are merged into the master branch.

---

# Testing

There are currently no automated tests in this project.

However, contributions adding tests are welcome.

If adding tests:

* Use Go's standard testing framework
* Keep tests simple and focused
* Ensure they are easy to run

---

# Release Workflow

Releases are created manually using GitHub Releases.

When a release is created:

* Docker image builds are triggered automatically

Contributors do not need to manage releases.

---

# Security Issues

If you discover a security vulnerability:

* Do not open a public issue for sensitive vulnerabilities
* Contact the maintainer privately
* Provide as much detail as possible

For non-sensitive issues, normal GitHub issues are acceptable.

---

# Code of Conduct

This project aims to be welcoming and respectful to all contributors.

Expected behavior:

* Be respectful and constructive
* Provide helpful feedback
* Assume good intentions
* Focus on improving the project

If a formal Code of Conduct is added in the future, this section will be updated.

---

# Final Notes

Thank you for contributing to imap-spam-cleaner!

Your time and effort help improve the project for everyone.

_This guide was created with the assistance of AI tools and adjusted by the project maintainer._
