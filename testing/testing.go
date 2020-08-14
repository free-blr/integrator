package testing

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	migrate "github.com/rubenv/sql-migrate"
	"gopkg.in/ory-am/dockertest.v3"
)

type DB struct {
	*sqlx.DB
	resource *dockertest.Resource
}

func NewDB() (*DB, error) {
	var (
		connString string
		err        error
	)
	db := new(DB)
	if isCI() {
		connString = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("POSTGRES_PORT_5432_TCP_ADDR"),
			os.Getenv("POSTGRES_PORT_5432_TCP_PORT"),
			os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_PASSWORD"),
			os.Getenv("POSTGRES_DB"),
		)
	} else {
		db.resource, err = runPostgresDocker()
		if err != nil {
			return nil, errors.Wrap(err, "cant run pg docker")
		}
		connString = fmt.Sprintf("user=postgres dbname=postgres sslmode=disable port=%s", db.resource.GetPort("5432/tcp"))
	}
	db.DB, err = openDb(connString)
	if err != nil {
		return nil, errors.Wrap(err, "cant open db")
	}
	return db, nil
}

func (db *DB) Close() error {
	if err := db.DB.Close(); err != nil {
		return errors.Wrap(err, "cant close db")
	}
	if db.resource != nil {
		return db.resource.Close()
	}
	return nil
}

func runPostgresDocker() (*dockertest.Resource, error) {
	pool, err := dockertest.NewPool("")

	if err != nil {
		return nil, errors.Wrap(err, "Could not connect to docker pool")
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("mdillon/postgis", "11-alpine", []string{})
	if err != nil {
		return nil, errors.Wrap(err, "Could not start resource")
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		cs := fmt.Sprintf("user=postgres dbname=postgres sslmode=disable port=%s", resource.GetPort("5432/tcp"))

		if _, err := openDb(cs); err != nil {
			return errors.Wrap(err, "cant open db")
		}

		return nil

	}); err != nil {
		return nil, errors.Wrap(err, "Could not connect to docker")
	}

	return resource, nil
}

func openDb(connectionString string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		return nil, errors.Wrap(err, "cant open db")
	}

	return db, nil
}

func MustMigrateDb(db *sqlx.DB) {
	migrateUp(db, dbMigrationsDir())
}

func MigrateDownDb(db *sqlx.DB) {
	migrateDown(db, dbMigrationsDir())
}

func migrateUp(db *sqlx.DB, dir string) {
	migrations := &migrate.FileMigrationSource{Dir: dir}

	n, err := migrate.Exec(db.DB, "postgres", migrations, migrate.Up)

	if err != nil {
		panic(errors.Errorf("could not apply migrations: %+v", err))
	}

	fmt.Printf("Applied %d migrations! (%s)\n", n, dir)
}

func migrateDown(db *sqlx.DB, dir string) {
	migrations := &migrate.FileMigrationSource{Dir: dir}

	n, err := migrate.Exec(db.DB, "postgres", migrations, migrate.Down)

	if err != nil {
		panic(errors.Errorf("could not apply migrations: %+v", err))
	}

	fmt.Printf("Applied %d migrations! (%s)\n", n, dir)
}

func dbMigrationsDir() string {
	return os.Getenv("DB_MIGRATIONS_PATH")
}

func isCI() bool {
	return os.Getenv("CI_SERVER") == "yes"
}
