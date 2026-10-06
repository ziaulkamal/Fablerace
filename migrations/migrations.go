// Package migrations meng-embed file SQL agar binary tunggal membawa skemanya
// sendiri (docs/09-devops.md §1). Dibaca oleh internal/platform/database.Migrate.
package migrations

import "embed"

//go:embed *.sql
var FS embed.FS
