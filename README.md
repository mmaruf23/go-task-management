# Go Task Management

### migration command

- up

```pwsh
migrate -path db/migration -database "postgres://postgres:password@localhost:5432/go_task_management?sslmode=disable" up

migrate -path db/migration -database "postgres://postgres:password@localhost:5432/go_task_management?sslmode=disable" down

migrate -path db/migration -database "postgres://postgres:password@localhost:5432/go_task_management?sslmode=disable" version
```
