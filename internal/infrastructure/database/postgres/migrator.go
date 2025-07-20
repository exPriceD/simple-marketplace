package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/exPriceD/simple-marketplace/migrations"
	"github.com/jackc/pgx/v5"
	"io/fs"
	"sort"
	"strings"
)

type migration struct {
	Name string
	SQL  string
}

// RunMigrations читает встроенные SQL файлы из migrations.Files и применяет их по имени.
func RunMigrations(ctx context.Context, db *DB) error {
	files, err := fs.ReadDir(migrations.Files, ".")
	if err != nil {
		return fmt.Errorf("read migrations fs: %w", err)
	}

	var migs []migration
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		name := f.Name()
		if !strings.HasSuffix(name, ".sql") {
			continue
		}
		content, err := fs.ReadFile(migrations.Files, name)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", name, err)
		}
		migs = append(migs, migration{
			Name: name,
			SQL:  string(content),
		})
	}

	sort.Slice(migs, func(i, j int) bool {
		return migs[i].Name < migs[j].Name
	})

	for _, m := range migs {
		if err := applyMigration(ctx, db, m); err != nil {
			return fmt.Errorf("apply %s: %w", m.Name, err)
		}
	}
	return nil
}

func applyMigration(ctx context.Context, db *DB, m migration) error {
	const ddl = `
CREATE TABLE IF NOT EXISTS schema_migrations (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);`
	if _, err := db.Pool.Exec(ctx, ddl); err != nil {
		return fmt.Errorf("ensure schema_migrations: %w", err)
	}

	var exists bool
	err := db.Pool.QueryRow(ctx,
		`SELECT true FROM schema_migrations WHERE name=$1 LIMIT 1`, m.Name).Scan(&exists)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("check migration %s: %w", m.Name, err)
	}
	if exists {
		return nil
	}

	tx, err := db.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if _, err := tx.Exec(ctx, m.SQL); err != nil {
		return fmt.Errorf("exec migration %s: %w", m.Name, err)
	}
	if _, err := tx.Exec(ctx,
		`INSERT INTO schema_migrations (name) VALUES ($1)`, m.Name); err != nil {
		return fmt.Errorf("record migration %s: %w", m.Name, err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit migration %s: %w", m.Name, err)
	}
	return nil
}
