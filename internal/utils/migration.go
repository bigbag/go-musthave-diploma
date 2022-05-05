package utils

import (
	"errors"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigration(databaseDSN string) (bool, error) {
	m, err := migrate.New("file://migration", databaseDSN)
	if err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return false, err
		}
	}
	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return false, err
		}
	}
	return true, nil
}
