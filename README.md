# Quick

Quick is a simple client/server implementation which allows the client to increment
 a `uint64` counter on the server and get the updated count in return.

## Dependencies

- [Go](https://golang.org/dl/)

## Setup Development Environment

```bash
# clone the repository (from github)
git clone https://github.com/bryanchriswhite/quick.git

# OR
# clone the repository (from a git bundle: ./quick.git)
git clone ./quick.git

# install go dependencies
go mod download
```

## Run Tests & Benchmarks

### All Tests

```bash
go test -race -tags integration,load,benchmark -bench=. -benchmem ./...
```

### Unit Tests

```bash
go test -race ./...
```

### Integration Tests

```bash
go test -race -tags integration ./integration/...
```

### Load Tests

```bash
go test -race -tags load ./integration/...
```

### Benchmarks

```bash
go test -race -tags benchmark -bench=. -benchmem ./integration/...
```

## Future Work

- [ ] add client CLI
- [ ] add server CLI
- [ ] add build tooling
- [ ] containerize everything (i.e. building CLI's, running tests, etc.)
- [ ] add support for HTTP transport & JSON serialization
- [ ] add support for gRPC transport & Protobuf serialization
- [ ] add benchmarks for HTTP & gRPC transport
