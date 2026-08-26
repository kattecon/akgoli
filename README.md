# AkGoLi

A Go utility library with small, focused packages that I use across different
private projects.

## Packages

| Package | Description |
|---|---|
| `absos` | Interfaces for DNS lookups and the system clock, with mock implementations for tests. |
| `appinfo` | Application name, version, and Go version metadata injected at build time via linker flags. |
| `logging` | Zap logger factory with Prometheus log-event counters by level. |
| `metrics` | Private Prometheus registry wrapper with application-aware metric naming. |
| `sbox` | Authenticated encryption (NaCl secretbox) for small JSON-serializable payloads. |
| `testutils` | Test helpers: buffered Zap logger, panic capture, stdout/stderr capture, HTTP transport fakes. |
| `typesafe` | Generic type-safe `sync.Map` wrapper. |
| `utils` | Constant-time string comparison, string masking, secure random IDs, buffer pools, sentinel errors, rune slicing, ASCII validation, in-place slice helpers. |

## Installation

```sh
go get github.com/kattecon/akgoli
```

## Usage

```go
import "github.com/kattecon/akgoli/sbox"

svc := sbox.NewSBoxSvc()
encrypted, err := svc.Encode(map[string]any{"user": "alice"})
if err != nil {
    log.Fatal(err)
}
// encrypted is a URL-safe Base64 string

var result map[string]any
if err := svc.Decode(encrypted, &result); err != nil {
    log.Fatal(err)
}
// result["user"] == "alice"
```

See each package's Go doc comments for the full API.

## Testing

Use the Makefile to run tests. It supplies linker flags required by the
`appinfo` package:

```sh
make test
```

The target also runs staticcheck and go vet after the tests.

Running `go test ./...` directly fails the `appinfo` tests because the
linker flags are missing. To run tests without the Makefile, pass equivalent
flags:

```sh
go test -ldflags="-X 'github.com/kattecon/akgoli/appinfo.version=3.2.1' -X 'github.com/kattecon/akgoli/appinfo.idName=test'" ./...
```
