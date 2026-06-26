package database

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Migrate applies the Up section of every goose-style .sql file in dir, in
// filename order, idempotently. The full template swaps this for `goose`; the
// skeleton keeps it dependency-free so `docker compose up` just works.
func (d *DB) Migrate(ctx context.Context, dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)

	if _, err := d.Pool.Exec(ctx,
		`CREATE TABLE IF NOT EXISTS schema_migrations (name TEXT PRIMARY KEY)`); err != nil {
		return err
	}

	for _, name := range files {
		var applied bool
		if err := d.Pool.QueryRow(ctx,
			`SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE name=$1)`, name).
			Scan(&applied); err != nil {
			return err
		}
		if applied {
			continue
		}

		b, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return err
		}
		up := extractUp(string(b))
		if strings.TrimSpace(up) != "" {
			if _, err := d.Pool.Exec(ctx, up); err != nil {
				return err
			}
		}
		if _, err := d.Pool.Exec(ctx,
			`INSERT INTO schema_migrations (name) VALUES ($1)`, name); err != nil {
			return err
		}
	}
	return nil
}

// extractUp returns the SQL between `-- +goose Up` and `-- +goose Down`.
func extractUp(content string) string {
	const upMarker = "-- +goose Up"
	const downMarker = "-- +goose Down"
	start := strings.Index(content, upMarker)
	if start == -1 {
		return content
	}
	body := content[start+len(upMarker):]
	if end := strings.Index(body, downMarker); end != -1 {
		body = body[:end]
	}
	return body
}
