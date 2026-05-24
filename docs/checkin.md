# Check-in

## CLI

```bash
go run ./cmd/checkin --email 'your-email@example.com' --passwd 'your-password'
```

Or put credentials in `config.yaml`:

```yaml
email: "your-email@example.com"
passwd: "your-password"
```

Or:

```bash
NEWORD_EMAIL='your-email@example.com' NEWORD_PASSWD='your-password' go run ./cmd/checkin
```

Priority order is `CLI > ENV > config.yaml`.

## Ubuntu deployment

1. Build the binary:

```bash
go build -o checkin ./cmd/checkin
```

2. Fill in `config.yaml`.
3. Install the timer:

```bash
sudo bash ./scripts/install-ubuntu-checkin.sh
```

This installs the binary and config to `/opt/neworld_check-in` and enables a `systemd` timer that runs once per day at a random time between `06:00` and `09:00`.

## Endpoint

- Method: `POST`
- URL: `https://neworld.space/user/checkin`

## Request headers

- `Content-Type: application/x-www-form-urlencoded`

## Request body

```text
checkin_type=time
```

## Notes

- This endpoint requires an authenticated session from a successful login.
- The Go CLI logs in internally, then submits the check-in request with the same cookie jar.

## Observed responses

Successful check-in response:

```json
{
  "ret": 1,
  "msg": "签到获得了 4 小时时长"
}
```

Already checked-in response:

```json
{
  "ret": 0,
  "msg": "今天已经签到过了"
}
```
