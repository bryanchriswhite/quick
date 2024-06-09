# Quick

Quick is a simple client/server implementation which allows the client to increment
 a `uint64` counter on the server and get the updated count in return.

## Dependencies

- [Go](https://golang.org/dl/)
- [make](https://www.gnu.org/software/make/) _(optional)_

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
make all-tests
```

### Unit Tests

```bash
make unit-tests
```

### Integration Tests

```bash
make integration-tests
```

### Load Tests

```bash
make load-tests
```

### Benchmarks

```bash
make benchmarks
```

## Future Work

- [ ] add client CLI
- [ ] add server CLI
- [ ] containerize everything (i.e. building CLI's, running tests, etc.)
- [ ] add support for HTTP transport & JSON serialization
- [ ] add support for gRPC transport & Protobuf serialization
- [ ] add benchmarks for HTTP & gRPC transport
