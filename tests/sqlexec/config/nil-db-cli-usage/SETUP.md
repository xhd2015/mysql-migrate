# Scenario

**Feature**: cli.Run status with nil cfg.DB is usage error exit 2

```
# usage gate before any sql.Open
cli.Run(cfg{DB:nil, MigrationsDir:tmp}, ["status"])
  -> exit 2
  -> Error mentions db / missing / config (not a hang)
```

## Preconditions

- Op=`cli_nil_db`.
- MigrationsDir is a non-empty temp dir so the failure is about DB, not dir.
- Offline — no MySQL.

## Steps

1. Create temp migrations dir; set MigrationsDir.
2. Set Op=`cli_nil_db`.
3. Assert exit 2 and combined output mentions db/missing/config.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = "cli_nil_db"
	req.MigrationsDir = t.TempDir()
	t.Logf("leaf config/nil-db-cli-usage: MigrationsDir=%s DB=nil", req.MigrationsDir)
	return nil
}
```
