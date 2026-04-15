# vaultpull

> CLI tool to sync secrets from HashiCorp Vault into local `.env` files with profile support

---

## Installation

```bash
go install github.com/yourusername/vaultpull@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourusername/vaultpull/releases).

---

## Usage

Configure your Vault address and token, then pull secrets into a `.env` file:

```bash
export VAULT_ADDR="https://vault.example.com"
export VAULT_TOKEN="s.xxxxxxxxxxxxxxxx"

# Pull secrets from a path into .env
vaultpull pull --path secret/myapp --out .env

# Use a named profile
vaultpull pull --profile production --out .env.production

# List available profiles
vaultpull profiles list
```

A profile can be defined in `~/.vaultpull.yaml`:

```yaml
profiles:
  production:
    addr: https://vault.example.com
    path: secret/myapp/prod
  staging:
    addr: https://vault.example.com
    path: secret/myapp/staging
```

---

## Flags

| Flag | Description |
|------|-------------|
| `--path` | Vault secret path to pull from |
| `--out` | Output `.env` file path (default: `.env`) |
| `--profile` | Named profile from config file |
| `--overwrite` | Overwrite existing values in the output file |

---

## License

[MIT](LICENSE) © 2024 yourusername