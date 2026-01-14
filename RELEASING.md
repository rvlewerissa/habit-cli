# Releasing Guide

## Prerequisites

```bash
brew install goreleaser gh
gh auth login
```

## Release Steps

1. **Verify build**
   ```bash
   go build ./...
   goreleaser build --snapshot --clean
   ```

2. **Tag version**
   ```bash
   git tag v1.0.0
   git push origin main --tags
   ```

3. **Release**
   ```bash
   goreleaser release --clean
   ```

## Versioning

```
v{MAJOR}.{MINOR}.{PATCH}
```

| Change | Bump |
|--------|------|
| Bug fix | v1.0.0 → v1.0.1 |
| New feature | v1.0.0 → v1.1.0 |
| Breaking change | v1.0.0 → v2.0.0 |

## Output

GoReleaser creates:
- GitHub Release with changelog
- Binaries for macOS, Linux, Windows (amd64/arm64)
- Checksums file
