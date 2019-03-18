package server

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"github.com/clementauger/migrate"
	"github.com/clementauger/monparcours/server/env"

	dbconf "github.com/clementauger/monparcours/server/config/db"
	"github.com/gobuffalo/packr"
)

var (
	statikBox = packr.NewBox("../migrations")
)

func withDB(f func(db *sql.DB, stage, dialect string, env *dbconf.Environment)) {

	stage := env.Stage()

	env, err := dbconf.GetEnvironment("dbconfig.yml", stage)
	if err != nil {
		log.Fatalf("Could not parse config: %s", err)
	}
	conn, dialect, err := dbconf.GetConnection(env)
	if err != nil {
		log.Fatalf("Could not connect using %q: %s", stage, err)
	}
	defer conn.Close()

	if env.TableName != "" {
		migrate.SetTable(env.TableName)
	}

	if env.SchemaName != "" {
		migrate.SetSchema(env.SchemaName)
	}

	f(conn, stage, dialect, env)
}

func getBox(statik bool, dialect string) migrate.MigrationSource {
	var box migrate.MigrationSource = migrate.FileMigrationSource{
		Dir: filepath.Join("migrations", dialect),
	}
	if statik {
		box = &migrate.PackrMigrationSource{
			Box: statikBox,
			Dir: dialect,
		}
	}
	return box
}

//Hello tests the db connection
func Hello(ctx context.Context) {
	withDB(func(db *sql.DB, stage, dialect string, env *dbconf.Environment) {
		err := db.Ping()
		if err != nil {
			log.Fatalf("Could not ping using %q: %s", stage, err)
		}
		fmt.Println(fmt.Sprintf("Connection %q (%v, statik=%v) OK", stage, dialect, env.Statik))
	})
}

//MigrateNow creates new migration
func MigrateNow(ctx context.Context) {
	withDB(func(db *sql.DB, stage, dialect string, env *dbconf.Environment) {

		env.Statik = false

		var dryrun bool
		flag.BoolVar(&dryrun, "dryrun", false, "Don't apply migrations, just print them.")
		flag.BoolVar(&env.Statik, "statik", env.Statik, "Use statik assets or not.")
		flag.Parse()
		name := flag.Arg(flag.NArg() - 1)

		migrate.NewMigrator(db, dialect, getBox(env.Statik, dialect)).DryRun(dryrun).MigrateNow(ctx, name)
	})
}

//MigrateUp applies migration.
func MigrateUp(ctx context.Context) {
	withDB(func(db *sql.DB, stage, dialect string, env *dbconf.Environment) {

		var limit int
		var dryrun bool
		flag.IntVar(&limit, "limit", 1, "Max number of migrations to apply.")
		flag.BoolVar(&dryrun, "dryrun", false, "Don't apply migrations, just print them.")
		flag.BoolVar(&env.Statik, "statik", env.Statik, "Use statik assets or not.")
		flag.Parse()

		migrate.NewMigrator(db, dialect, getBox(env.Statik, dialect)).DryRun(dryrun).MigrateUp(ctx, limit)
	})
}

//MigrateDown uninstalls migrations
func MigrateDown(ctx context.Context) {
	withDB(func(db *sql.DB, stage, dialect string, env *dbconf.Environment) {

		var limit int
		var dryrun bool
		flag.IntVar(&limit, "limit", 1, "Max number of migrations to apply.")
		flag.BoolVar(&env.Statik, "statik", env.Statik, "Use statik assets or not.")
		flag.BoolVar(&dryrun, "dryrun", false, "Don't apply migrations, just print them.")
		flag.Parse()

		migrate.NewMigrator(db, dialect, getBox(env.Statik, dialect)).DryRun(dryrun).MigrateDown(ctx, limit)
	})
}

//MigrateSkip skips migrations
func MigrateSkip(ctx context.Context) {
	withDB(func(db *sql.DB, stage, dialect string, env *dbconf.Environment) {

		var limit int
		var dryrun bool
		flag.IntVar(&limit, "limit", 1, "Max number of migrations to apply.")
		flag.BoolVar(&dryrun, "dryrun", false, "Don't apply migrations, just print them.")
		flag.BoolVar(&env.Statik, "statik", env.Statik, "Use statik assets or not.")
		flag.Parse()

		migrate.NewMigrator(db, dialect, getBox(env.Statik, dialect)).DryRun(dryrun).MigrateSkip(ctx, limit)
	})
}

// MigrateRedo checks migrations
func MigrateRedo(ctx context.Context) {
	withDB(func(db *sql.DB, stage, dialect string, env *dbconf.Environment) {

		var dryrun bool
		flag.BoolVar(&dryrun, "dryrun", false, "Don't apply migrations, just print them.")
		flag.BoolVar(&env.Statik, "statik", env.Statik, "Use statik assets or not.")
		flag.Parse()

		migrate.NewMigrator(db, dialect, getBox(env.Statik, dialect)).DryRun(dryrun).MigrateRedo(ctx)
	})
}

//MigrateStatus displays migrations statuses
func MigrateStatus(ctx context.Context) {
	withDB(func(db *sql.DB, stage, dialect string, env *dbconf.Environment) {

		var dryrun bool
		flag.BoolVar(&dryrun, "dryrun", false, "Don't apply migrations, just print them.")
		flag.BoolVar(&env.Statik, "statik", env.Statik, "Use statik assets or not.")
		flag.Parse()

		migrate.NewMigrator(db, dialect, getBox(env.Statik, dialect)).DryRun(dryrun).MigrateStatus(ctx)
	})
}
