# learning-go

A start of my practicing Go, just a day into learning it.

Combines upstream JSON endpoints for a user into a single result API, using
purely stdlib and handler-level tests. Requests to upstream services are made
concurrently. Error conditions are at least mostly handled. All in a single file
because I haven't learned how Go projects are typically structured yet.

## Running Locally

```bash
go run main.go

# in other terminal
curl http://localhost:8081/v1/user-posts/1
```

## Tests

```bash
go test
```

