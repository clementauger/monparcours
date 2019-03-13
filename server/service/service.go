package service

import (
	"database/sql"
	"fmt"
	"log"

	dbconf "github.com/clementauger/monparcours/server/config/db"
	"github.com/clementauger/monparcours/server/model"
	mysqlmodel "github.com/clementauger/monparcours/server/model/mysql"
	// pgsqlmodel "github.com/clementauger/monparcours/server/model/pgsql"
)

//Service gathers all data access services
type Service struct {
	DB *sql.DB
	// Dialect        string
	Protest        model.ProtestService
	Step           model.StepService
	ContactMessage model.ContactMessageService
}

//Close the service
func (h *Service) Close() error {
	return h.DB.Close()
}

//Init service model
func (h *Service) Init(stage string) error {

	var db *sql.DB
	var dialect string
	{
		env, err := dbconf.GetEnvironment("dbconfig.yml", stage)
		if err != nil {
			log.Fatal(err)
		}
		conn, x, err := dbconf.GetConnection(env)
		if err != nil {
			log.Fatal(err)
		}
		db = conn
		dialect = x
	}

	h.DB = db
	// h.Dialect = dialect

	if dialect == "sqlite3" || dialect == "mysql" {
		h.Step = mysqlmodel.StepService{DB: db}
		h.Protest = mysqlmodel.ProtestService{DB: db}
		h.ContactMessage = mysqlmodel.ContactMessageService{DB: db}
		// } else if dialect == "postgres" {
		// 	app.StepService = pgsqlmodel.StepService{DB: db}
		// 	app.ProtestService = pgsqlmodel.ProtestService{DB: db}
		// 	app.ContactMessageService = pgsqlmodel.ContactMessageService{DB: db}
	} else {
		return fmt.Errorf("unknown dialect %q", dialect)
	}

	return nil
}
