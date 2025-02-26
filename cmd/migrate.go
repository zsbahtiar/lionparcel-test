package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
)

var (
	migrationDir string
	migrationExt string
)

func migrateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Database migration commands",
		Long:  `Commands to run database migrations: up, down, version, create`,
	}

	cmd.PersistentFlags().StringVarP(&migrationDir, "dir", "d", "internal/migration", "Migration files directory")
	cmd.PersistentFlags().StringVarP(&migrationExt, "ext", "e", "sql", "Migration files extension")

	cmd.AddCommand(migrateUpCmd())
	cmd.AddCommand(migrateDownCmd())
	cmd.AddCommand(migrateVersionCmd())
	cmd.AddCommand(migrateCreateCmd())

	return cmd
}

func migrateUpCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "Run all up migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := getMigrate()
			if err != nil {
				return err
			}

			if err := m.Up(); err != nil && err != migrate.ErrNoChange {
				return fmt.Errorf("error running up migrations: %v", err)
			}

			log.Println("Successfully ran up migrations")
			return nil
		},
	}
}

func migrateDownCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "down",
		Short: "Run all down migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := getMigrate()
			if err != nil {
				return err
			}

			if err := m.Down(); err != nil && err != migrate.ErrNoChange {
				return fmt.Errorf("error running down migrations: %v", err)
			}

			log.Println("Successfully ran down migrations")
			return nil
		},
	}
}

func migrateVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show current migration version",
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := getMigrate()
			if err != nil {
				return err
			}

			version, dirty, err := m.Version()
			if err != nil {
				return fmt.Errorf("error getting migration version: %v", err)
			}

			log.Printf("Current migration version: %d, Dirty: %v", version, dirty)
			return nil
		},
	}
}

func getMigrate() (*migrate.Migrate, error) {

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)

	m, err := migrate.New("file://internal/migration", dbURL)
	if err != nil {
		return nil, fmt.Errorf("error creating migrate instance: %v", err)
	}

	return m, nil
}

func migrateCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create [migration name]",
		Short: "Create a new migration",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires a migration name argument")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			return createMigrationFiles(name)
		},
	}
}

func createMigrationFiles(name string) error {
	if err := os.MkdirAll(migrationDir, 0755); err != nil {
		return fmt.Errorf("failed to create migrations directory: %v", err)
	}

	timestamp := time.Now().Format("20060102150405")

	upFile := filepath.Join(migrationDir, fmt.Sprintf("%s_%s.up.%s", timestamp, name, migrationExt))
	if err := createFile(upFile, getUpTemplate(name)); err != nil {
		return fmt.Errorf("failed to create up migration: %v", err)
	}

	downFile := filepath.Join(migrationDir, fmt.Sprintf("%s_%s.down.%s", timestamp, name, migrationExt))
	if err := createFile(downFile, getDownTemplate(name)); err != nil {
		return fmt.Errorf("failed to create down migration: %v", err)
	}

	log.Printf("Created migration files:\n  %s\n  %s", upFile, downFile)
	return nil
}

func createFile(filename string, content string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(content)
	return err
}

func getUpTemplate(name string) string {
	return fmt.Sprintf(`-- Migration Up: %s
BEGIN;


COMMIT;
`, name)
}

func getDownTemplate(name string) string {
	return fmt.Sprintf(`-- Migration Down: %s

BEGIN;

COMMIT;
`, name)
}
