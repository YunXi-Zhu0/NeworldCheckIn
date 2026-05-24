# Checkin

Go CLI for running the daily `neworld.space` check-in flow.

## Command

```bash
go run ./cmd/checkin --email 'your-email@example.com' --passwd 'your-password'
```

You can also put credentials in `config.yaml` at the repo root:

```yaml
email: "your-email@example.com"
passwd: "your-password"
```

The CLI also supports environment variables:

```bash
NEWORD_EMAIL='your-email@example.com' NEWORD_PASSWD='your-password' go run ./cmd/checkin
```

Priority order is `CLI > ENV > config.yaml`. Use `config.yaml.example` as the template, and keep the real `config.yaml` out of Git.
