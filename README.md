# polymer
My playground to implement everything I find it interesting

---

## Supabase 2.109.1

## Node.js 24.16.0

## npm 11.13.0

## pnpm 11.5.2

## Nest 11.0.21

---

## Go 1.26.4

## Atlas 1.2.0

### Sqlfluff 4.2.2

### Mockery v3.7.1

### sqlc v1.31.1

### golangci-lint 2.12.2

### fieldalignment 0.46.0

---

## How to Run

### Supabase

- Go to `/supabase` directory
- Run `supabase start`
- Useful command
    ```
    supabase db reset
    supabase stop
    ```

### Account

- Go to `/backend/services/account` directory
- Run `pnpm run db:migrate`
- Run `pnpm run start`

### Wallet

- Go to `/backend/services/wallet` directory
- Run `make migrate`
- Run `go run cmd/server/main.go`