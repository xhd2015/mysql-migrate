# Scenario

**Feature**: migrate package import path is github.com/xhd2015/mysql-migrate/migrate

```
# import path identity for Config
reflect.TypeOf(migrate.Config{})
  -> PkgPath == github.com/xhd2015/mysql-migrate/migrate
  -> Name == Config
```

## Preconditions

- Package `github.com/xhd2015/mysql-migrate/migrate` is importable from tests.
- Type `Config` is defined in that package (not an alias in another package).

## Steps

1. Set `req.Mode` to `module-importable`.
2. Root `Run` reflects on `migrate.Config{}` and records `PkgPath` and type name.

## Context

- Proves the locked import path, not only that some Config type exists.
- Classic RED: missing package → compile failure of the generated harness.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Mode = "module-importable"
	return nil
}
```
