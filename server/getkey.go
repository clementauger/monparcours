package server

import (
	"context"
	"flag"
	"fmt"
	"log"

	myapp "github.com/clementauger/monparcours/server/app"
	appconf "github.com/clementauger/monparcours/server/config/app"
	"github.com/clementauger/monparcours/server/env"
	// pgsqlmodel "github.com/clementauger/monparcours/server/model/pgsql"
)

//GetKey returns the admin key
func Getkey(ctx context.Context) {
	flag.Parse()

	stage := env.Stage()

	var appConfig appconf.Environment
	{
		env, err := appconf.GetEnvironment("app.yml", stage)
		if err != nil {
			log.Fatal(err)
		}
		appConfig = *env
	}

	fmt.Print(myapp.GetKey(appConfig))
}
