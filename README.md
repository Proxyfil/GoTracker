# GoTracker

## Requirements

- Go 1.24^
- Postgres 15^

## Setup

1. Start postgres container
```bash
docker run --name postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres:15
```

2. Get go dependencies
```bash
go get github.com/lib/pq
```

3. Run project
```bash
cd src
go run .
```