package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	_ "expvar"
	_ "net/http/pprof"

	"github.com/clementauger/commander"
	"github.com/clementauger/monparcours/server"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// build info
var (
	Tag     = "0.0.0"
	Version = "dev"
	Appname = "monparcours"
	banner  = fmt.Sprintf(`%v %v %v`, Appname, Version, Tag)
)

func main() {
	cmder := commander.New(banner, Appname)

	cmder.Add("serve", server.ServeHTTP).
		Alias("s").Description(`serve http application`)

	cmder.Add("hello", server.Hello).
		Alias("").Description(`print and test configuration`)

	cmder.Add("getkey", server.Getkey).
		Alias("g").Description(`show admin key`)

	cmder.Add("migrate", server.MigrateNow).
		Alias("m").Description(`create a new named migration`)

	cmder.Add("migratestatus", server.MigrateStatus).
		Alias("ms").Description(`show migrations status`)

	cmder.Add("migrateskip", server.MigrateSkip).
		Alias("mk").Description(`skip next migration`)

	cmder.Add("migrateup", server.MigrateUp).
		Alias("mu").Description(`apply migrations`)

	cmder.Add("migratedown", server.MigrateDown).
		Alias("md").Description(`revert migrations`)

	cmder.Add("migrateredo", server.MigrateRedo).
		Alias("mr").Description(`redo last migration`)

	cmder.Add("help", nil).
		Alias("-h --help -help").Description(`show help`)

	cmder.MustRun(context.Background())

}

//todo: report to golang. want improved error message that includes the expected format.
//./main.go:31: struct field tag `json:id` not compatible with reflect.StructTag.Get: bad syntax for struct tag value

func onSignal(s os.Signal, h func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, s)
	<-c
	h()
}
