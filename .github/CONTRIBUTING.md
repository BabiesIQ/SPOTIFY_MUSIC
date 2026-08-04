# Contributing to SPOTIFY_MUSIC

Thanks for taking the time to contribute! 🎉

## How to Contribute

### Reporting Bugs
- Check if the issue already exists before opening a new one.
- Include your Go version, OS, and relevant logs.
- Provide steps to reproduce the problem.

### Suggesting Features
- Open an issue with the label `enhancement`.
- Describe the use case clearly — why is it useful?

### Submitting Code

1. **Fork** the repository.
2. **Create a branch** from `main`:
   ```bash
   git checkout -b feat/your-feature-name
   ```
3. **Write your code** — keep it clean and consistent with the existing style.
4. **Test** your changes locally before submitting.
5. **Commit** using conventional commit format:
   ```
   feat: add spotify playlist support
   fix: resolve queue skip race condition
   docs: update deployment guide
   ```
6. **Push** your branch and open a **Pull Request** against `main`.

## Code Style

- Follow standard Go formatting — run `gofmt` before committing.
- Keep functions small and focused.
- Add comments for non-obvious logic.
- Avoid unnecessary external dependencies.

## Pull Request Checklist

- [ ] Code compiles without errors (`go build ./...`)
- [ ] No new linting warnings
- [ ] Relevant documentation updated (if applicable)
- [ ] PR description explains *what* and *why*

## Questions?

Open an issue or reach out via the support group linked in the [README](../README.md).
