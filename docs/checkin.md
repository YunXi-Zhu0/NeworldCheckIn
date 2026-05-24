# Check-in

## CLI

```bash
go run ./cmd/checkin --email 'your-email@example.com' --passwd 'your-password'
```

Or:

```bash
NEWORD_EMAIL='your-email@example.com' NEWORD_PASSWD='your-password' go run ./cmd/checkin
```

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
