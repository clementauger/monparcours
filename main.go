package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strings"

	_ "expvar"
	_ "net/http/pprof"

	"github.com/clementauger/monparcours/server" // TODO: Replace with the absolute import path
	// TODO: Replace with the absolute import path
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// build info
var (
	Tag     = "0.0.0"
	Version = "dev"
	Appname = "monparcours"
	actions = map[string]func(context.Context){
		"hello":         server.Hello,
		"migrate":       server.MigrateNow,
		"migrateup":     server.MigrateUp,
		"migratedown":   server.MigrateDown,
		"migrateredo":   server.MigrateRedo,
		"migrateskip":   server.MigrateSkip,
		"migratestatus": server.MigrateStatus,
		"getkey":        server.Getkey,
		"serve":         server.ServeHTTP,
		"help":          help,
	}
	descriptions = map[string]string{
		"hello":         `test db connection`,
		"migrate":       `create a new named migration.`,
		"migrateup":     `apply migrations`,
		"migratedown":   `revert migrations`,
		"migrateredo":   `redo last migration`,
		"migrateskip":   `skip next migration`,
		"migratestatus": `show migrations status`,
		"getkey":        `show admin key`,
		"serve":         `serve http application`,
		"help":          `show help`,
	}
	aliases = map[string]string{
		"migrate":       `m`,
		"migrateup":     `mu`,
		"migratedown":   `md`,
		"migrateredo":   `mr`,
		"migrateskip":   `mk`,
		"migratestatus": `ms`,
		"getkey":        `g`,
		"serve":         `s`,
		"help":          `-h --help`,
	}
	opts = map[string]string{
		"migrate": `[name]`,
	}
	action = "serve"
	banner = fmt.Sprintf(`%v %v %v`, Appname, Version, Tag)
)

func main() {

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go onSignal(os.Interrupt, cancel)

	a := action
	if len(os.Args) > 1 {
		a = os.Args[1]
	}

	for c, alias := range aliases {
		alias = " " + alias + " "
		if strings.Index(alias, " "+a+" ") > -1 {
			a = c
			break
		}
	}
	if _, ok := actions[a]; !ok {
		if !strings.HasPrefix(a, "-") {
			help(nil)
			fmt.Println()
			fmt.Println("unknown command ", a)
			fmt.Println()
			os.Exit(2)
		}
		a = action
	} else if len(os.Args) > 1 {
		os.Args = append(os.Args[:1], os.Args[2:]...)
	}

	flag.Usage = func() {
		d, _ := descriptions[a]
		t, _ := aliases[a]
		c, _ := opts[a]
		t = strings.Replace(t, " ", "|", -1)
		if t != "" {
			t = "|" + t
		}
		w := flag.CommandLine.Output()
		fmt.Fprintln(w, banner)
		fmt.Fprintln(w, "")
		fmt.Fprintf(w, `  %v
`, a+t+" "+c)
		fmt.Fprintln(w, "")
		fmt.Fprintln(w, "  "+d)
		fmt.Fprintln(w, "")
		flag.PrintDefaults()
	}

	if fn, ok := actions[a]; ok {
		fn(ctx)
	}
}

var h string

func init() {
	h = fmt.Sprintf(`%v

 usage
    %v [action] args...
    %v args...

 actions
`, banner, Appname, Appname)
	keys := []string{}
	for a := range actions {
		keys = append(keys, a)
	}
	sort.Strings(keys)
	for _, a := range keys {
		d, _ := descriptions[a]
		t, _ := aliases[a]
		c, _ := opts[a]
		t = strings.Replace(t, " ", "|", -1)
		if t != "" {
			t = "|" + t
		}
		h += fmt.Sprintf(`    %v
	%v
`, a+t+" "+c, d)
	}
}
func help(ctx context.Context) {
	fmt.Print(h)
}

//todo: report to golang. want improved error message that includes the expected format.
//./main.go:31: struct field tag `json:id` not compatible with reflect.StructTag.Get: bad syntax for struct tag value

func onSignal(s os.Signal, h func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, s)
	<-c
	h()
}
