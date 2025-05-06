// internal/migrations/migrator.go
package migrations

import (
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strings"

	"gorm.io/gorm"
)

//go:embed sql/*.sql
var migrationFS embed.FS

func Migrate(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		files, err := getMigrationFiles()
		if err != nil {
			return fmt.Errorf("failed to read migrations: %w", err)
		}

		for _, file := range files {
			if err := executeSQL(tx, file); err != nil {
				return fmt.Errorf("migration %s failed: %w", file.name, err)
			}
		}
		return nil
	})
}

type migrationFile struct {
	name    string
	content string
}

func getMigrationFiles() ([]migrationFile, error) {
	entries, err := fs.ReadDir(migrationFS, "sql")
	if err != nil {
		return nil, err
	}

	var files []migrationFile
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		content, err := fs.ReadFile(migrationFS, "sql/"+entry.Name())
		if err != nil {
			return nil, err
		}

		files = append(files, migrationFile{
			name:    entry.Name(),
			content: string(content),
		})
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].name < files[j].name
	})

	return files, nil
}

// internal/migrations/migrator.go (фрагмент)
func executeSQL(tx *gorm.DB, file migrationFile) error {
	commands := strings.Split(file.content, ";")
	for _, cmd := range commands {
		cmd = strings.TrimSpace(cmd)
		if cmd == "" {
			continue
		}
		if err := tx.Exec(cmd).Error; err != nil {
			return fmt.Errorf("ошибка в команде '%s': %w", cmd, err)
		}
	}
	return nil
}
