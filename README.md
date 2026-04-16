# envoy-cli

A lightweight CLI for managing and diffing `.env` files across environments.

---

## Installation

```bash
go install github.com/yourusername/envoy-cli@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envoy-cli.git
cd envoy-cli
go build -o envoy-cli .
```

---

## Usage

```bash
# Diff two .env files
envoy-cli diff .env.development .env.production

# Check for missing keys between environments
envoy-cli check --base .env.example --target .env.local

# Merge env files (target takes precedence)
envoy-cli merge .env.defaults .env.local -o .env.merged
```

Example output:

```
~ DB_HOST        dev-db.local → prod-db.internal
+ SENTRY_DSN     (missing in development)
- DEBUG          (missing in production)
```

---

## Features

- Diff `.env` files and highlight added, removed, and changed keys
- Detect missing keys against a base template
- Merge multiple env files with configurable precedence
- Zero external dependencies

---

## License

MIT © [yourusername](https://github.com/yourusername)