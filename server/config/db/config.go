package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"gopkg.in/gorp.v1"

	"github.com/clementauger/coders"
	_ "github.com/go-sql-driver/mysql" // they are here,
	_ "github.com/lib/pq"              // and maybe
	_ "github.com/mattn/go-sqlite3"    // only here,
	// deal with it.
)

// Environment defines connection string and execution options related to the databaase.
type Environment struct {
	//Dialect is oen of sqlite3, mysql, pgsql
	Dialect string `yaml:"dialect"`
	//DataSource is the dsn
	DataSource string `yaml:"datasource"`
	//TableName of the migrations.
	TableName string `yaml:"table"`
	//SchemaName (database) of the migrations.
	SchemaName string `yaml:"schema"`
	//ConnMaxLifetime runtime options.
	ConnMaxLifetime *time.Duration `yaml:"connmaxlifetime"`
	//MaxIdleConns runtime options.
	MaxIdleConns *int `yaml:"maxidleconns"`
	//MaxOpenConns runtime options.
	MaxOpenConns *int `yaml:"maxopenconns"`
	//Statik defines where the application should look for to find the migration files.
	Statik bool `yaml:"statik"`
}

var dialects = map[string]gorp.Dialect{
	"sqlite3":  gorp.SqliteDialect{},
	"postgres": gorp.PostgresDialect{},
	"mysql":    gorp.MySQLDialect{Engine: "InnoDB", Encoding: "UTF8"},
}

var afterOpen = map[string]func(*sql.DB) error{
	"sqlite3":  func(db *sql.DB) error { return nil },
	"postgres": func(db *sql.DB) error { return nil },
	"mysql": func(db *sql.DB) error {
		_, err := db.Exec("SET SQL_MODE='ALLOW_INVALID_DATES';")
		return err
	},
}

func GetEnvironment(filename, environment string) (*Environment, error) {
	conf := make(map[string]*Environment)
	err := coders.Decode(conf, filename)
	if err != nil {
		return nil, err
	}

	env := conf[environment]
	if env == nil {
		return nil, fmt.Errorf("environment %q does not have configuration", environment)
	}

	if env.Dialect == "" {
		return nil, errors.New("No dialect specified")
	}

	if env.DataSource == "" {
		return nil, errors.New("No data source specified")
	}
	env.DataSource = os.ExpandEnv(env.DataSource)

	if env.ConnMaxLifetime == nil {
		y := time.Hour
		env.ConnMaxLifetime = &y
	}
	if env.MaxIdleConns == nil {
		y := 10
		env.MaxIdleConns = &y
	}
	if env.MaxOpenConns == nil {
		y := 10
		env.MaxOpenConns = &y
	}

	return env, nil
}

func GetConnection(env *Environment) (*sql.DB, string, error) {

	// Make sure we only accept dialects that were compiled in.
	_, exists := dialects[env.Dialect]
	if !exists {
		return nil, "", fmt.Errorf("Unsupported dialect: %s", env.Dialect)
	}

	db, err := sql.Open(env.Dialect, env.DataSource)
	if err != nil {
		return nil, "", fmt.Errorf("Cannot connect to database: %s", err)
	}

	if env.ConnMaxLifetime != nil {
		db.SetConnMaxLifetime(*env.ConnMaxLifetime)
	}
	if env.MaxIdleConns != nil {
		db.SetMaxIdleConns(*env.MaxIdleConns)
	}
	if env.MaxOpenConns != nil {
		db.SetMaxOpenConns(*env.MaxOpenConns)
	}

	if fn, ok := afterOpen[env.Dialect]; ok {
		if err = fn(db); err != nil {
			return db, env.Dialect, err
		}
	}

	return db, env.Dialect, nil
}
