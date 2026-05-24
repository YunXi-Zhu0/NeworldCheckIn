# Checkin

Go CLI for running the daily `neworld.space` check-in flow.

## Command

```bash
go run ./cmd/checkin --email 'your-email@example.com' --passwd 'your-password'
```

The CLI also supports environment variables:

```bash
NEWORD_EMAIL='your-email@example.com' NEWORD_PASSWD='your-password' go run ./cmd/checkin
```
