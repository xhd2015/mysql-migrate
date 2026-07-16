// Package plan is a pure migration plan status machine.
// Given inventory files and prior apply log rows, Build classifies each
// migration as skip / apply / blocked / deferred and sets HasBlock.
// No MySQL, no filesystem, no log writes.
package plan

import (
	"sort"

	"github.com/xhd2015/mysql-migrate/migrate/inventory"
)

// LogRow is a prior apply attempt snapshot for one migration_id.
type LogRow struct {
	MigrationID   string
	Status        string // running | success | failed | unknown | pending | ""
	ExactlyOnce   bool   // snapshot; file.ExactlyOnce is source of truth for rules
	ContentSHA256 string
	DurationMS    int
	ErrorMessage  string
	Note          string
}

// ItemAction is the plan action for one migration.
type ItemAction string

const (
	ActionSkip     ItemAction = "skip"
	ActionApply    ItemAction = "apply"
	ActionBlocked  ItemAction = "blocked"
	ActionDeferred ItemAction = "deferred"
)

// PlanItem is one classified migration in the plan.
type PlanItem struct {
	MigrationID  string
	ExactlyOnce  bool
	Action       ItemAction
	Reason       string
	HashMismatch bool
	// LogStatus is the effective log status after normalizing stale running → unknown.
	LogStatus string
}

// Plan is the ordered classification outcome for a set of migrations.
type Plan struct {
	Items    []PlanItem
	HasBlock bool
}

// Build classifies each migration file against log rows and applies stop-chain.
// Items are emitted in FileName ascending order. file.ExactlyOnce is the source
// of truth for ExactlyOnce on PlanItem (log ExactlyOnce is ignored for rules).
//
// Classification (per file):
//   - no log / status "" / "pending" → apply
//   - success + (log hash empty OR log hash == file hash) → skip
//   - success + non-empty log hash ≠ file hash → blocked, HashMismatch
//   - failed / unknown → blocked
//   - running → effective LogStatus unknown → blocked (stale)
//
// Stop-chain: after the first blocked item, later items that would apply become
// deferred; skip stays skip. HasBlock is true iff any item is blocked.
func Build(files []inventory.MigrationFile, logs []LogRow) Plan {
	if len(files) == 0 {
		return Plan{Items: []PlanItem{}, HasBlock: false}
	}

	// Stable order: FileName ascending (inventory list order).
	ordered := make([]inventory.MigrationFile, len(files))
	copy(ordered, files)
	sort.Slice(ordered, func(i, j int) bool {
		return ordered[i].FileName < ordered[j].FileName
	})

	// Index logs by MigrationID (first row wins if duplicates).
	byID := make(map[string]LogRow, len(logs))
	for _, row := range logs {
		if _, ok := byID[row.MigrationID]; ok {
			continue
		}
		byID[row.MigrationID] = row
	}

	items := make([]PlanItem, 0, len(ordered))
	seenBlock := false
	hasBlock := false

	for _, f := range ordered {
		item := classify(f, byID)
		if item.Action == ActionBlocked {
			hasBlock = true
			seenBlock = true
		} else if seenBlock && item.Action == ActionApply {
			item.Action = ActionDeferred
			if item.Reason == "" {
				item.Reason = "deferred_after_block"
			}
		}
		items = append(items, item)
	}

	return Plan{Items: items, HasBlock: hasBlock}
}

// classify maps one file + optional log row to a PlanItem (pre-stop-chain).
func classify(f inventory.MigrationFile, byID map[string]LogRow) PlanItem {
	item := PlanItem{
		MigrationID: f.ID,
		ExactlyOnce: f.ExactlyOnce,
	}

	row, hasLog := byID[f.ID]
	if !hasLog {
		item.Action = ActionApply
		item.LogStatus = ""
		item.Reason = "no_log"
		return item
	}

	status := row.Status
	// Stale running is normalized to unknown before classification.
	if status == "running" {
		item.LogStatus = "unknown"
		item.Action = ActionBlocked
		item.Reason = "stale_running"
		return item
	}
	item.LogStatus = status

	switch status {
	case "", "pending":
		item.Action = ActionApply
		if status == "pending" {
			item.Reason = "pending"
		} else {
			item.Reason = "empty_status"
		}
		return item

	case "success":
		// Empty log hash is treated as compatible (not a mismatch).
		if row.ContentSHA256 == "" || row.ContentSHA256 == f.ContentSHA256 {
			item.Action = ActionSkip
			item.Reason = "success"
			return item
		}
		item.Action = ActionBlocked
		item.HashMismatch = true
		item.Reason = "hash_mismatch"
		return item

	case "failed":
		item.Action = ActionBlocked
		if f.ExactlyOnce {
			item.Reason = "exactly_once_failed"
		} else {
			item.Reason = "failed"
		}
		return item

	case "unknown":
		item.Action = ActionBlocked
		if f.ExactlyOnce {
			item.Reason = "exactly_once_unknown"
		} else {
			item.Reason = "unknown"
		}
		return item

	default:
		// Unrecognized status is a human gate (same class as unknown).
		item.Action = ActionBlocked
		item.Reason = "unknown_status:" + status
		return item
	}
}
