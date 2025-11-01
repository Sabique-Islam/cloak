
---

# Contributing

Follow these guidelines to help maintain code quality.

---

## Getting Started

- Fork the repository and clone your fork.
- Set up your Go environment and install dependencies. "[Check SETUP.md for assistance](https://github.com/Sabique-Islam/cloak/blob/master/.github/SETUP.md)"
- Create a new branch for your feature or fix.
- Make changes, test locally, and commit.

---

## Commit Messages

- Follow the conventions defined in [`COMMIT_CONVENTION.md`](https://github.com/Sabique-Islam/cloak/blob/master/.github/COMMIT_CONVENTION.md).
- If you need a new commit type or convention, propose it by adding to "TYPES" section of `COMMIT_CONVENTION.md` in your PR.

---

## Pull Requests

- Use the PR template ([`PULL_REQUEST_TEMPLATE.md`](https://github.com/Sabique-Islam/cloak/blob/master/.github/PULL_REQUEST_TEMPLATE.md)) when opening a pull request.

---

## Testing

- Write unit tests for new features and bug fixes.
- Ensure existing tests pass before submitting.
- Test database interactions and API calls where applicable.

---

## Linting & Checks

Before submitting a PR, run the following checks locally:

### Linting :-

```
golangci-lint run
```

Ensure that no linting errors remain.

### Run Tests :-

```
go test ./... -v
```

All tests should pass.

### Formatting :-

```
go fmt ./...
go vet ./...
```

Ensure code is properly formatted.

---
