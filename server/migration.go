package server

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/clementauger/monparcours/server/dbconnect"
	"github.com/clementauger/monparcours/server/env"

	"github.com/gobuffalo/packr"
	"github.com/olekukonko/tablewriter"
	migrate "github.com/rubenv/sql-migrate"
)

var (
	box = packr.NewBox("../migrations")
)

func getBox(env *dbconnect.Environment) migrate.MigrationSource {
	if env.Statik {
		return &migrate.PackrMigrationSource{
			Box: box,
			Dir: env.Dialect,
		}
	}
	return migrate.FileMigrationSource{
		Dir: filepath.Join("migrations", env.Dialect),
	}
}

func Hello(ctx context.Context) {

	flag.Parse()

	stage := env.Stage()

	env, err := dbconnect.GetEnvironment("dbconfig.yml", stage)
	if err != nil {
		log.Fatalf("Could not parse config: %s", err)
	}
	conn, dialect, err := dbconnect.GetConnection(env)
	if err != nil {
		log.Fatalf("Could not connect using %q: %s", stage, err)
	}
	defer conn.Close()

	fmt.Println(fmt.Sprintf("Connection %q (%v, statik=%v) OK", stage, dialect, env.Statik))
}

func MigrateNow(ctx context.Context) {

	flag.Parse()

	name := flag.Arg(flag.NArg() - 1)

	if name == "" {
		log.Fatal("name is required")
	}

	stage := env.Stage()

	var templateContent = `
-- +migrate Up
-- +migrate Down
`
	tpl := template.Must(template.New("new_migration").Parse(templateContent))

	env, err := dbconnect.GetEnvironment("dbconfig.yml", stage)
	if err != nil {
		log.Fatalf("Could not parse config: %s", err)
	}

	dir := filepath.Join("migrations", env.Dialect)
	os.MkdirAll(dir, os.ModePerm)

	fileName := fmt.Sprintf("%s-%s.sql", time.Now().Format("20060102150405"), strings.TrimSpace(name))
	pathName := path.Join(dir, fileName)

	f, err := os.Create(pathName)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = f.Close() }()

	if err := tpl.Execute(f, nil); err != nil {
		log.Fatal(err)
	}

	fmt.Println(fmt.Sprintf("Created migration %s", pathName))
}

func MigrateUp(ctx context.Context) {
	var limit int
	var dryrun bool
	flag.IntVar(&limit, "limit", 1, "Max number of migrations to apply.")
	flag.BoolVar(&dryrun, "dryrun", false, "Don't apply migrations, just print them.")
	flag.Parse()

	stage := env.Stage()

	err := applyMigrations(stage, migrate.Up, dryrun, limit)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func MigrateDown(ctx context.Context) {
	var limit int
	var dryrun bool
	flag.IntVar(&limit, "limit", 1, "Max number of migrations to apply.")
	flag.BoolVar(&dryrun, "dryrun", false, "Don't apply migrations, just print them.")
	flag.Parse()

	stage := env.Stage()

	err := applyMigrations(stage, migrate.Down, dryrun, limit)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func MigrateSkip(ctx context.Context) {
	var limit int
	flag.IntVar(&limit, "limit", 1, "Max number of migrations to apply.")
	flag.Parse()

	stage := env.Stage()

	env, err := dbconnect.GetEnvironment("dbconfig.yml", stage)
	if err != nil {
		log.Fatalf("Could not parse config: %s", err)
	}

	db, dialect, err := dbconnect.GetConnection(env)
	if err != nil {
		log.Fatal(err)
	}

	source := getBox(env)

	n, err := migrate.SkipMax(db, dialect, source, migrate.Up, limit)
	if err != nil {
		log.Fatalf("Migration failed: %s", err)
	}

	fmt.Println("Skipped 1 migration")

	if n == 1 {
		fmt.Println("Skipped 1 migration")
	} else {
		fmt.Println(fmt.Sprintf("Skipped %d migrations", n))
	}
}

func MigrateRedo(ctx context.Context) {
	var limit int
	var dryrun bool
	flag.IntVar(&limit, "limit", 1, "Max number of migrations to apply.")
	flag.BoolVar(&dryrun, "dryrun", false, "Don't apply migrations, just print them.")
	flag.Parse()

	stage := env.Stage()

	env, err := dbconnect.GetEnvironment("dbconfig.yml", stage)
	if err != nil {
		log.Fatalf("Could not parse config: %s", err)
	}

	db, dialect, err := dbconnect.GetConnection(env)
	if err != nil {
		log.Fatal(err)
	}

	source := getBox(env)

	migrations, _, err := migrate.PlanMigration(db, dialect, source, migrate.Down, 1)
	if len(migrations) == 0 {
		if err != nil {
			log.Println(err)
		}
		fmt.Println("Nothing to do!")
		return
	}
	if err != nil {
		log.Fatal(err)
	}

	if dryrun {
		printMigration(migrations[0], migrate.Down)
		printMigration(migrations[0], migrate.Up)
	} else {
		_, err := migrate.ExecMax(db, dialect, source, migrate.Down, 1)
		if err != nil {
			log.Fatal(fmt.Sprintf("Migration (down) failed: %s", err))
		}

		_, err = migrate.ExecMax(db, dialect, source, migrate.Up, 1)
		if err != nil {
			log.Fatal(fmt.Sprintf("Migration (up) failed: %s", err))
		}

		fmt.Println(fmt.Sprintf("Reapplied migration %s.", migrations[0].Id))
	}

}

func MigrateStatus(ctx context.Context) {

	flag.Parse()

	stage := env.Stage()

	env, err := dbconnect.GetEnvironment("dbconfig.yml", stage)
	if err != nil {
		log.Fatalf("Could not parse config: %s", err)
	}

	db, dialect, err := dbconnect.GetConnection(env)
	if err != nil {
		log.Fatal(err)
	}

	source := getBox(env)

	migrations, err := source.FindMigrations()
	if err != nil {
		log.Fatal(err)
	}

	records, err := migrate.GetMigrationRecords(db, dialect)
	if err != nil {
		log.Fatal(err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Migration", "Applied"})
	table.SetColWidth(60)

	rows := make(map[string]*statusRow)

	for _, m := range migrations {
		rows[m.Id] = &statusRow{
			ID:       m.Id,
			Migrated: false,
		}
	}

	for _, r := range records {
		if rows[r.Id] == nil {
			log.Println(fmt.Sprintf("Could not find migration file: %v", r.Id))
			continue
		}

		rows[r.Id].Migrated = true
		rows[r.Id].AppliedAt = r.AppliedAt
	}

	for _, m := range migrations {
		if rows[m.Id] != nil && rows[m.Id].Migrated {
			table.Append([]string{
				m.Id,
				rows[m.Id].AppliedAt.String(),
			})
		} else {
			table.Append([]string{
				m.Id,
				"no",
			})
		}
	}
	table.Render()

}

type statusRow struct {
	ID        string
	Migrated  bool
	AppliedAt time.Time
}

func applyMigrations(stage string, dir migrate.MigrationDirection, dryrun bool, limit int) error {
	env, err := dbconnect.GetEnvironment("dbconfig.yml", stage)
	if err != nil {
		return fmt.Errorf("Could not parse config: %s", err)
	}

	db, dialect, err := dbconnect.GetConnection(env)
	if err != nil {
		return err
	}

	source := getBox(env)

	if dryrun {
		migrations, _, err := migrate.PlanMigration(db, dialect, source, dir, limit)
		if err != nil {
			return fmt.Errorf("Cannot plan migration: %s", err)
		}
		for _, m := range migrations {
			printMigration(m, dir)
		}
	} else {
		n, err := migrate.ExecMax(db, dialect, source, dir, limit)
		if err != nil {
			return fmt.Errorf("Migration failed: %s", err)
		}

		if n == 1 {
			fmt.Println("Applied 1 migration")
		} else {
			fmt.Println(fmt.Sprintf("Applied %d migrations", n))
		}
	}

	return nil
}

func printMigration(m *migrate.PlannedMigration, dir migrate.MigrationDirection) {
	if dir == migrate.Up {
		fmt.Println(fmt.Sprintf("==> Would apply migration %s (up)", m.Id))
		for _, q := range m.Up {
			fmt.Println(q)
		}
	} else if dir == migrate.Down {
		fmt.Println(fmt.Sprintf("==> Would apply migration %s (down)", m.Id))
		for _, q := range m.Down {
			fmt.Println(q)
		}
	} else {
		panic("Not reached")
	}
}
