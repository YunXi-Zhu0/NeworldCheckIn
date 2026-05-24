# Login

## Endpoint

- Method: `POST`
- URL: `https://neworld.space/auth/login`

## Request headers

- `Content-Type: application/json`

## Request body

```json
{
  "email": "your-email@example.com",
  "passwd": "your-password"
}
```

## Observed responses

Successful login response:

```json
{
  "ret": 1,
  "msg": "登录成功，欢迎回来"
}
```

## Notes

- This repository no longer exposes a standalone login command.
- The login endpoint is still used internally by the check-in flow to establish the authenticated session.
