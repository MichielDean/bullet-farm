package cistern

import (
	"testing"
)

func TestLoadMigrations(t *testing.T) {
	migrations, err := loadMigrations()
	if err != nil {
		t.Fatalf("loadMigrations: %v", err)
	}
	if len(migrations) == 0 {
		t.Fatal("expected at least one migration, got 0")
	}
	t.Logf("Loaded %d migrations", len(migrations))
	for _, m := range migrations {
		stmts := splitStatements(m.SQL)
		t.Logf("  %s: %d statements", m.migrationID(), len(stmts))
	}
}

func TestSplitStatements(t *testing.T) {
	sql := `-- Migration 004: Normalize stored repo values to canonical casing.
UPDATE "droplets" SET repo = 'cistern' WHERE LOWER(repo) = LOWER('cistern') AND repo != 'cistern';
UPDATE "droplets" SET repo = 'ScaledTest' WHERE LOWER(repo) = LOWER('ScaledTest') AND repo != 'ScaledTest';`
	stmts := splitStatements(sql)
	if len(stmts) != 2 {
		t.Fatalf("expected 2 statements, got %d: %v", len(stmts), stmts)
	}
	for i, s := range stmts {
		t.Logf("  [%d] %s", i, s[:min(80, len(s))])
	}
}

func TestMigrationsAreApplied(t *testing.T) {
	c := testClient(t)
	defer c.Close()

	var count int
	if err := c.db.QueryRow(`SELECT COUNT(*) FROM "_schema_migrations"`).Scan(&count); err != nil {
		t.Fatalf("query _schema_migrations: %v", err)
	}
	if count == 0 {
		t.Errorf("expected migrations to be recorded, got %d rows", count)
	}

	// Check that the numbered migration IDs exist
	rows, err := c.db.Query(`SELECT "id" FROM "_schema_migrations" ORDER BY "id"`)
	if err != nil {
		t.Fatalf("query _schema_migrations: %v", err)
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			t.Fatal(err)
		}
		ids = append(ids, id)
	}
	t.Logf("Migration IDs: %v", ids)

	// All DDL migrations that succeeded should be tracked
	// At minimum the DML migrations (003, 004) should be tracked
	var count003 int
	c.db.QueryRow(`SELECT COUNT(*) FROM "_schema_migrations" WHERE "id" = '003_complexity_renumber'`).Scan(&count003)
	if count003 != 1 {
		t.Errorf("expected 003_complexity_renumber to be tracked, got count=%d", count003)
	}

	var count004 int
	c.db.QueryRow(`SELECT COUNT(*) FROM "_schema_migrations" WHERE "id" = '004_repo_case_normalize'`).Scan(&count004)
	if count004 != 1 {
		t.Errorf("expected 004_repo_case_normalize to be tracked, got count=%d", count004)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
