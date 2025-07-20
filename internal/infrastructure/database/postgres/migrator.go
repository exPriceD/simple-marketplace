package postgres

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strings"
)

//go:embed ../../../migrations/*.sql
var migrationFiles embed.FS

// Примечание: Путь в //go:embed relatif к каталогу файла. Данный файл находится в internal/infrastructure/database/postgres,
// поэтому к корневой migrations папке путь ../../../migrations/*.sql. (Go модуль корень — там же go.mod).

// RunMigrations применяет SQL файлы из корневой папки migrations.
func RunMigrations(ctx context.Context, db *DB) error {
	if db == nil || db.Pool == nil {
		return fmt.Errorf("nil db")
	}

	if _, err := db.Pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations(
			name text PRIMARY KEY,
			applied_at timestamptz NOT NULL DEFAULT now()
		)`); err != nil {
		return fmt.Errorf("create schema_migrations: %w", err)
	}

	entries, err := fs.ReadDir(migrationFiles, "../../../migrations")
	if err != nil {
		return fmt.Errorf("read migrations: %w", err)
	}
	var files []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".sql") {
			continue
		}
		files = append(files, e.Name())
	}
	sort.Strings(files)

	for _, fname := range files {
		var exists bool
		if err := db.Pool.QueryRow(ctx,
			`SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE name=$1)`, fname,
		).Scan(&exists); err != nil {
			return fmt.Errorf("check migration %s: %w", fname, err)
		}
		if exists {
			continue
		}
		content, err := migrationFiles.ReadFile("../../../migrations/" + fname)
		if err != nil {
			return fmt.Errorf("read file %s: %w", fname, err)
		}
		if _, err := db.Pool.Exec(ctx, string(content)); err != nil {
			return fmt.Errorf("apply %s: %w", fname, err)
		}
		if _, err := db.Pool.Exec(ctx, `INSERT INTO schema_migrations(name) VALUES($1)`, fname); err != nil {
			return fmt.Errorf("record %s: %w", fname, err)
		}
	}
	return nil
}
