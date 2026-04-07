# Contributing

Thank you for your interest in contributing to Twitch AI Bot! Here are some guidelines to follow.

## Commit Messages

This project uses [Conventional Commits](https://www.conventionalcommits.org/). Every commit must follow this format:

```
<type>: <description>
```

### Types

| Type | Description |
|---|---|
| `feat` | New feature |
| `fix` | Bug fix |
| `refactor` | Code change that neither fixes a bug nor adds a feature |
| `chore` | Maintenance tasks, config changes, dependencies |
| `docs` | Documentation changes |
| `style` | Formatting, missing semicolons, etc (no code change) |
| `test` | Adding or updating tests |
| `perf` | Performance improvements |

### Rules

- Use lowercase for the description
- Write in present tense ("add" not "added")
- Keep descriptions concise and clear
- No period at the end

### Examples

```
feat: add groq ai provider integration
fix: prevent bot from replying to its own messages
refactor: rename storage package for consistency
chore: update go dependencies
docs: add setup instructions to readme
```

## Branch Naming

Use descriptive branch names with a prefix that indicates the type of work:

```
feat/short-description
fix/short-description
refactor/short-description
chore/short-description
```

Examples:

```
feat/add-openai-provider
fix/websocket-reconnection
refactor/ai-interface
```

## How to Contribute

1. **Fork** the repository
2. **Create a branch** from `main` following the naming convention above
3. **Make your changes**
4. **Commit** using conventional commit format
5. **Push** to your fork
6. **Open a Pull Request** with a clear description of what was changed and why

## Pull Request Guidelines

- Fill out the PR template if one exists
- Keep changes focused — one concern per PR
- Reference related issues in the description
- Make sure your code builds and runs locally before submitting

## Code Style

- Follow idiomatic Go conventions
- Run `go fmt ./...` before committing
- Keep functions small and focused
- Add comments only when intent is not obvious from the code itself
