# CI workflow note

## Status

**updated** (pushed)

## Branch / remote / push

| Field | Value |
|-------|--------|
| Branch | `master-2026-08-06-use-go-best-practice-to-review-current-project` |
| Remote | `origin` → `ssh://git@github.com/xhd2015/mysql-migrate.git` |
| Upstream | `origin/master-2026-08-06-use-go-best-practice-to-review-current-project` |
| Push result | success — CI workflow at `ab63caf`; subsequent docs commits for this note also pushed |
| CI commit SHA | `ab63cafc2305f459d21acbb0df8c1b1fe22bfcec` (short: `ab63caf`) |
| Branch tip (after note commits) | see `git rev-parse origin/master-2026-08-06-use-go-best-practice-to-review-current-project` |

## Paths changed (CI commit)

- `.github/workflows/test.yml` — full doctest-pattern CI (coverage, stages, merge, artifacts)
- `script/ci/coverage-package-table.py` — package coverage markdown for step summary

## How to view Actions for this push

- Repo: https://github.com/xhd2015/mysql-migrate
- Actions (branch filter): https://github.com/xhd2015/mysql-migrate/actions?query=branch%3Amaster-2026-08-06-use-go-best-practice-to-review-current-project
- Workflow file on branch: https://github.com/xhd2015/mysql-migrate/blob/master-2026-08-06-use-go-best-practice-to-review-current-project/.github/workflows/test.yml
- CI commit: https://github.com/xhd2015/mysql-migrate/commit/ab63cafc2305f459d21acbb0df8c1b1fe22bfcec

## How this differs from doctest’s workflow

| Aspect | doctest reference | this repo |
|--------|-------------------|-----------|
| Triggers | `push` + `pull_request` | same |
| Runner | `ubuntu-latest` + `setup-go` from `go.mod` | same |
| `COVERPKG` | `github.com/xhd2015/doctest/...` | `github.com/xhd2015/mysql-migrate/...` |
| Go test + coverage | yes | yes |
| Doctest install | `go install ./cmd/doctest` (this checkout) | `go install github.com/xhd2015/doctest/cmd/doctest@latest` |
| Discovery | `--label '!e2e'` + cover | same |
| e2e stage | `--label e2e` + cover | same (no e2e labels yet; stage kept for parity) |
| xgo merge | yes | yes |
| Package table | `script/ci/coverage-package-table.py` (doctest paths) | same helper adapted to mysql-migrate module; skips `cmd/` |
| Artifacts | coverage profiles | same |
| MySQL | N/A | no service container; MySQL leaves skip when DSN unreachable |

## Notes

- Prior workflow used a golang container, no coverage, and a single `doctest test --label-all` step.
- Local pre-commit may warn that the file “differs from recommended workflow”; this intentional richer pattern matches the doctest project reference.
