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

#### Results

|          | ns/op        | ops/min      |
|----------|--------------|--------------|
|          | 5.70E+04     | 1.05E+06     |
|          | 6.07E+04     | 9.88E+05     |
|          | 5.81E+04     | 1.03E+06     |
|          | 6.23E+04     | 9.63E+05     |
|          | 6.10E+04     | 9.83E+05     |
|          | 6.27E+04     | 9.57E+05     |
|          | 6.07E+04     | 9.88E+05     |
|          | 6.12E+04     | 9.80E+05     |
|          | 6.04E+04     | 9.94E+05     |
|          | 5.76E+04     | 1.04E+06     |
| ---      | ---          | ---          |
| **Avg.** | **6.02E+04** | **9.98E+05** |

````
goos: linux
goarch: amd64
pkg: github.com/bryanchriswhite/quick/integration
cpu: AMD Ryzen Threadripper 1920X 12-Core Processor 
BenchmarkClient_Increment-24               18451             57585 ns/op            8416 B/op          8 allocs/op
PASS
ok      github.com/bryanchriswhite/quick/integration    3.158s

````

## Future Work

- [ ] add client CLI
- [ ] add server CLI
- [ ] containerize everything (i.e. building CLI's, running tests, etc.)
- [ ] add support for HTTP transport & JSON serialization
- [ ] add support for gRPC transport & Protobuf serialization
- [ ] add benchmarks for HTTP & gRPC transport
