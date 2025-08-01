package migrate

import (
	"errors"
	"fmt"

	gomigrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // postgresql engine for migrate
	_ "github.com/golang-migrate/migrate/v4/source/file"       // engine for migration files
)

var _ Migrate = (*pgMigrate)(nil)

// Migrate implementation for PostgreSQL.
type pgMigrate struct {
	mgrt *gomigrate.Migrate
}

func NewPostgreSQLMigrate(sourceURL, databaseURL string) (Migrate, error) {
	mgrt, err := gomigrate.New(sourceURL, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("create migrate manager: %w", err)
	}
	return &pgMigrate{mgrt: mgrt}, nil
}

func (c *pgMigrate) Status() (version uint, isDirty bool, err error) {
	v, d, err := c.mgrt.Version()
	if err != nil && !errors.Is(err, gomigrate.ErrNilVersion) {
		return 0, false, fmt.Errorf("migrate status: %w", err)
	}
	return v, d, nil
}

func (c *pgMigrate) Up() error {
	err := c.mgrt.Up()
	if err != nil && !errors.Is(err, gomigrate.ErrNoChange) {
		return fmt.Errorf("migrate up: %w", err)
	}
	return nil
}

func (c *pgMigrate) Down() error {
	err := c.mgrt.Down()
	if err != nil && !errors.Is(err, gomigrate.ErrNoChange) {
		return fmt.Errorf("migrate down: %w", err)
	}
	return nil
}

func (c *pgMigrate) Step(n int) error {
	err := c.mgrt.Steps(n)
	if err != nil {
		return fmt.Errorf("migrate step: %w", err)
	}
	return nil
}

func (c *pgMigrate) Force(n int) error {
	err := c.mgrt.Force(n)
	if err != nil {
		return fmt.Errorf("migrate force: %w", err)
	}
	return nil
}

func (c *pgMigrate) Close() error {
	sourceCloseErr, dbCloseErr := c.mgrt.Close()
	if sourceCloseErr != nil && dbCloseErr != nil {
		return fmt.Errorf("close database: %s && close source: %w", dbCloseErr.Error(), sourceCloseErr)
	}
	if sourceCloseErr != nil {
		return fmt.Errorf("close source: %w", sourceCloseErr)
	}
	if dbCloseErr != nil {
		return fmt.Errorf("close database: %w", dbCloseErr)
	}
	return nil
}
