package service

import (
	"database/sql"
	"log"

	dbconf "github.com/clementauger/monparcours/server/config/db"
	"github.com/clementauger/monparcours/server/model"
	"github.com/clementauger/crud"
	cruded "github.com/clementauger/monparcours/server/model/cruded"
)

//Service gathers all data access services
type Service struct {
	DB             *sql.DB
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

	crud.Logger = nil // disable

	h.Protest = &cruded.ProtestService{
		Crud:    crud.MakeCrud(db, dialect, []model.Protest{}),
		Dialect: dialect,
	}
	h.Step = &cruded.StepService{
		Crud: crud.MakeCrud(db, dialect, []model.Step{}),
		// Dialect: dialect,
	}
	h.ContactMessage = &cruded.ContactMessageService{
		Crud: crud.MakeCrud(db, dialect, []model.ContactMessage{}),
		// Dialect:dialect,
	}

	// if dialect == "sqlite3" || dialect == "mysql" {
	// 	h.Protest = mysqlProtestService{DB: db}
	// 	h.Step = mysqlStepService{DB: db}
	// 	h.ContactMessage = mysqlContactMessageService{DB: db}
	// } else if dialect == "postgres" {
	// 	h.Protest = pgsqlProtestService{DB: db}
	// 	h.Step = pgsqlStepService{DB: db}
	// 	h.ContactMessage = pgsqlContactMessageService{DB: db}
	// 	h.Protest = &pgsqlProtestService{
	// 		Crud: crud.MakeCrud(db, dialect, []Protest{}),
	// 	}
	// 	h.Step = &pgsqlStepService{
	// 		Crud: crud.MakeCrud(db, dialect, []Step{}),
	// 	}
	// 	h.ContactMessage = &pgsqlContactMessageService{
	// 		Crud: crud.MakeCrud(db, dialect, []ContactMessage{}),
	// 	}
	// } else {
	// 	return fmt.Errorf("unknown dialect %q", dialect)
	// }

	return nil
}
