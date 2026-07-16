# Scenario

**Feature**: Config has DB field typed as sqlexec.DB and no DSN field

```
# reflect field surface
type Config struct {
  DB sqlexec.DB
  MigrationsDir, ProgramName, AppliedBy string
}
// FieldByName("DSN") not found
// zero Config.DB == nil
```

## Preconditions

- Op=`config_fields`.
- Offline.

## Steps

1. Set Op.
2. Assert ConfigHasDB, ConfigDBIsIface, !ConfigHasDSN, ConfigDBNilZero.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = "config_fields"
	t.Log("leaf config/db-field-no-dsn: reflect Config fields")
	return nil
}
```
