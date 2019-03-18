package model_test

import (
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/clementauger/crud"
	"github.com/clementauger/monparcours/server/model"
	crudmodel "github.com/clementauger/monparcours/server/model/cruded"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var (
	db      *sql.DB
	dialect string
)

var (
	Contact model.ContactMessageService
	Protest model.ProtestService
	Step    model.StepService
)

func TestMain(m *testing.M) {

	dialect = "mysql"
	dialect = "postgres"
	dialect = "sqlite3"

	if dialect == "postgres" {
		x, err := sql.Open(dialect, "dbname=monparcours user=test password=test sslmode=disable")
		if err != nil {
			log.Fatal(err)
		}
		defer x.Close()
		db = x
	} else if dialect == "sqlite3" {
		x, err := sql.Open(dialect, "../../data/monparcours.db")
		if err != nil {
			log.Fatal(err)
		}
		defer x.Close()
		db = x
	} else if dialect == "mysql" {
		x, err := sql.Open(dialect, "test:test@/monparcours?parseTime=true")
		if err != nil {
			log.Fatal(err)
		}
		defer x.Close()
		db = x
	}

	// {
	// 	if dialect == "sqlite3" || dialect == "mysql" {
	// 		Protest = mysqlmodel.ProtestService{DB: db}
	// 		Step = mysqlmodel.StepService{DB: db}
	// 		Contact = mysqlmodel.ContactMessageService{DB: db}
	// 	} else if dialect == "postgres" {
	// 		Protest = pgsqlmodel.ProtestService{DB: db}
	// 		Step = pgsqlmodel.StepService{DB: db}
	// 		Contact = pgsqlmodel.ContactMessageService{DB: db}
	// 	}
	// }
	Protest = &crudmodel.ProtestService{Dialect: dialect, Crud: crud.MakeCrud(db, dialect, []model.Protest{})}
	Step = &crudmodel.StepService{Crud: crud.MakeCrud(db, dialect, []model.Step{})}
	Contact = &crudmodel.ContactMessageService{Crud: crud.MakeCrud(db, dialect, []model.ContactMessage{})}

	retCode := m.Run()
	os.Exit(retCode)
}

func TestProtestInsert(t *testing.T) {
	notCreatedAt := time.Now().Add(time.Hour * 48)
	d := model.Protest{
		ID:        -1,
		Title:     "test",
		CreatedAt: notCreatedAt,
		Public:    true,
	}
	var err error
	d, err = Protest.Insert(d)
	if err != nil {
		t.Fatal(err)
	}
	if d.ID == -1 {
		t.Fatal("wanted OID to be set, current value = -1")
	}
	if d.CreatedAt == notCreatedAt {
		t.Fatal("wanted CreatedAt to be set, current value = ", notCreatedAt)
	}
}

func TestProtestGetById(t *testing.T) {
	notCreatedAt := time.Now().Add(time.Hour * 48)
	d := model.Protest{
		ID:        -1,
		Title:     "test",
		CreatedAt: notCreatedAt,
		Public:    true,
	}
	var err error
	d, err = Protest.Insert(d)
	if err != nil {
		t.Fatal(err)
	}
	d2, err := Protest.Get(d.ID)
	if err != nil {
		t.Fatal(err)
	}
	if d2.ID == -1 {
		t.Fatal("wanted OID to be set, current value = -1")
	}
	if d2.CreatedAt == notCreatedAt {
		t.Fatal("wanted CreatedAt to be set, current value = ", notCreatedAt)
	}
	if d2.UpdatedAt != nil {
		t.Fatal("wanted UpdatedAt to be nil, current value = ", d2.UpdatedAt)
	}
	if d2.DeletedAt != nil {
		t.Fatal("wanted DeletedAt to be nil, current value = ", d2.DeletedAt)
	}
}

func TestProtestSelectWithPwd(t *testing.T) {
	notCreatedAt := time.Now().Add(time.Hour * 48)
	d := model.Protest{
		ID:        -1,
		Title:     "test",
		CreatedAt: notCreatedAt,
		Password:  "testpwd",
		Public:    false,
	}
	var err error
	d, err = Protest.Insert(d)
	if err != nil {
		t.Fatal(err)
	}
	var d2 model.Protest
	d2, err = Protest.Get(d.ID)
	if err != nil {
		t.Fatal(err)
	}
	if d2.ID != -1 {
		t.Fatal("wanted OID to be eq -1, current value = ", d2.ID)
	}
	var d3 model.Protest
	d3, err = Protest.GetWithPassword(d.ID, "testpwd")
	if err != nil {
		t.Fatal(err)
	}
	if d3.ID != d.ID {
		t.Fatal("wanted ID to eq", d.ID, ", current value = ", d3.ID)
	}
	if d3.CreatedAt == notCreatedAt {
		t.Fatal("wanted CreatedAt to be set, current value = ", notCreatedAt)
	}
	if d3.UpdatedAt != nil {
		t.Fatal("wanted UpdatedAt to be nil, current value = ", d3.UpdatedAt)
	}
	if d3.DeletedAt != nil {
		t.Fatal("wanted DeletedAt to be nil, current value = ", d3.DeletedAt)
	}
}

func TestProtestSelectByAuthorID(t *testing.T) {
	authorID := "author-id"
	notCreatedAt := time.Now().Add(time.Hour * 48)
	d := model.Protest{
		ID:        -1,
		Title:     "test",
		AuthorID:  authorID,
		CreatedAt: notCreatedAt,
		Public:    true,
	}
	var err error
	d, err = Protest.Insert(d)
	if err != nil {
		t.Fatal(err)
	}
	var d2 []model.Protest
	d2, err = Protest.GetByAuthorID(d.AuthorID)
	if err != nil {
		t.Fatal(err)
	}
	if len(d2) == 0 {
		t.Fatal("wanted len(results) to eq 1, current value = ", len(d2))
	}
	if d2[0].ID == -1 {
		t.Fatal("wanted OID to not eq -1, current value = ", d2[0].ID)
	}
	if d2[0].AuthorID != authorID {
		t.Fatal("wanted AuthoID to eq ", authorID, ", current value = ", d2[0].AuthorID)
	}
	if d2[0].CreatedAt == notCreatedAt {
		t.Fatal("wanted created_at to be set, current value = ", d2[0].CreatedAt)
	}
}

func TestProtestSelectByIDs(t *testing.T) {
	notCreatedAt := time.Now().Add(time.Hour * 48)
	d := model.Protest{
		ID:        -1,
		Title:     "test",
		CreatedAt: notCreatedAt,
		Public:    true,
	}
	var err error
	d, err = Protest.Insert(d)
	if err != nil {
		t.Fatal(err)
	}
	var d2 []model.Protest
	d2, err = Protest.GetByIDs(d.ID, d.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(d2) == 0 {
		t.Fatal("wanted len(results) to eq 1, current value = ", len(d2))
	}
	if d2[0].ID == -1 {
		t.Fatal("wanted OID to not eq -1, current value = ", d2[0].ID)
	}
	if d2[0].ID != d.ID {
		t.Fatal("wanted ID to eq ", d.ID, ", current value = ", d2[0].ID)
	}
	if d2[0].CreatedAt == notCreatedAt {
		t.Fatal("wanted created_at to be set, current value = ", d2[0].CreatedAt)
	}
}

func TestProtestSearch(t *testing.T) {
	var d2 []model.Protest
	var atLat, atLng *float64
	var startDate, endDate *time.Time
	var title, protest, organizer string
	title = "test"
	var err error
	d2, err = Protest.SearchProtests(
		&title, &protest, &organizer, startDate, endDate, atLat, atLng, 0.0,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(d2) == 0 {
		t.Fatal("wanted len(results) to eq 1, current value = ", len(d2))
	}
	if d2[0].ID == -1 {
		t.Fatal("wanted OID to not eq -1, current value = ", d2[0].ID)
	}
	if d2[0].Title != "test" {
		t.Fatal("wanted ID to eq test, current value = ", d2[0].Title)
	}
}

func TestProtestDelete(t *testing.T) {
	d := model.Protest{
		ID:     -1,
		Public: true,
	}
	var err error
	d, err = Protest.Insert(d)
	if err != nil {
		t.Fatal(err)
	}
	if d.ID == -1 {
		t.Fatal("wanted OID to be set, current value = -1")
	}
	err = Protest.Delete(d)
	if err != nil {
		t.Fatal(err)
	}
	d, err = Protest.Get(d.ID)
	if err != nil {
		t.Fatal(err)
	}
	if d.ID != -1 {
		t.Fatal("wanted OID to eq -1, current value = ", d.ID)
	}
}

func TestStepInsert(t *testing.T) {
	gatherAt := time.Now().Add(time.Hour)
	d := model.Protest{
		ID:       -1,
		Title:    "test",
		Public:   true,
		GatherAt: gatherAt,
		Steps:    []model.Step{model.Step{Place: "place", GatherAt: gatherAt, Details: "details"}},
	}
	var err error
	d, err = Protest.Insert(d)
	if err != nil {
		t.Fatal(err)
	}
	d.Steps[0], err = Step.Insert(d.Steps[0])
	if err != nil {
		t.Fatal(err)
	}
	if d.Steps[0].ID == -1 {
		t.Fatal("wanted d.Steps[0].ID to be set, current value = -1")
	}
	if d.Steps[0].Place != "place" {
		t.Fatal("wanted d.Steps[0].Place to eq place, current value = ", d.Steps[0].Place)
	}
	if d.Steps[0].Details != "details" {
		t.Fatal("wanted d.Steps[0].Details to eq details, current value = ", d.Steps[0].Details)
	}

	var shortANSIC = "Mon Jan _2 15 2006"
	if d.Steps[0].GatherAt.Format(shortANSIC) != d.GatherAt.Format(shortANSIC) {
		t.Fatal("wanted d.Steps[0].GatherAt to eq ", d.GatherAt, ", current value = ", d.Steps[0].GatherAt)
	}
}

func TestStepInsertAll(t *testing.T) {
	gatherAt := time.Now().Add(time.Hour)
	d := model.Protest{
		ID:       -1,
		Title:    "test",
		Public:   true,
		GatherAt: gatherAt,
		Steps:    []model.Step{model.Step{Place: "place", GatherAt: gatherAt, Details: "details"}},
	}
	var err error
	d, err = Protest.Insert(d)
	if err != nil {
		t.Fatal(err)
	}
	d.Steps, err = Step.InsertAll(d.Steps)
	if err != nil {
		t.Fatal(err)
	}
	if d.Steps[0].ID == -1 {
		t.Fatal("wanted d.Steps[0].ID to be set, current value = -1")
	}
	if d.Steps[0].Place != "place" {
		t.Fatal("wanted d.Steps[0].Place to eq place, current value = ", d.Steps[0].Place)
	}
	if d.Steps[0].Details != "details" {
		t.Fatal("wanted d.Steps[0].Details to eq details, current value = ", d.Steps[0].Details)
	}
	var shortANSIC = "Mon Jan _2 15 2006"
	if d.Steps[0].GatherAt.Format(shortANSIC) != d.GatherAt.Format(shortANSIC) {
		t.Fatal("wanted d.Steps[0].GatherAt to eq ", d.GatherAt, ", current value = ", d.Steps[0].GatherAt)
	}
}

func TestStepInsertSteps(t *testing.T) {
	gatherAt := time.Now().Add(time.Hour)
	d := model.Protest{
		ID:       -1,
		Title:    "test",
		Public:   true,
		GatherAt: gatherAt,
		Steps:    []model.Step{model.Step{Place: "place", GatherAt: gatherAt, Details: "details"}},
	}
	var err error
	d, err = Protest.Insert(d)
	if err != nil {
		t.Fatal(err)
	}
	d, err = Step.InsertSteps(d)
	if err != nil {
		t.Fatal(err)
	}
	if d.Steps[0].ID == -1 {
		t.Fatal("wanted d.Steps[0].ID to be set, current value = -1")
	}
	if d.Steps[0].ProtestID != d.ID {
		t.Fatal("wanted d.Steps[0].ID to eq", d.ID, ", current value = ", d.Steps[0].ID)
	}
	if d.Steps[0].Place != "place" {
		t.Fatal("wanted d.Steps[0].Place to eq place, current value = ", d.Steps[0].Place)
	}
	if d.Steps[0].Details != "details" {
		t.Fatal("wanted d.Steps[0].Details to eq details, current value = ", d.Steps[0].Details)
	}
	var shortANSIC = "Mon Jan _2 15 2006"
	if d.Steps[0].GatherAt.Format(shortANSIC) != d.GatherAt.Format(shortANSIC) {
		t.Fatal("wanted d.Steps[0].GatherAt to eq ", d.GatherAt, ", current value = ", d.Steps[0].GatherAt)
	}
}

func TestStepGetById(t *testing.T) {
	gatherAt := time.Now().Add(time.Hour * 48)
	d := model.Step{
		GatherAt: gatherAt,
	}
	var err error
	d, err = Step.Insert(d)
	if err != nil {
		t.Fatal(err)
	}
	d2, err := Step.Get(d.ID)
	if err != nil {
		t.Fatal(err)
	}
	if d2.ID == -1 {
		t.Fatal("wanted ID to be set, current value = -1")
	}
	var shortANSIC = "Mon Jan _2 15 2006"
	if d2.GatherAt.Format(shortANSIC) != d.GatherAt.Format(shortANSIC) {
		t.Fatal("wanted GatherAt to eq ", d.GatherAt.Format(shortANSIC), ", current value = ", d2.GatherAt.Format(shortANSIC))
	}
}

func TestStepDelete(t *testing.T) {
	d := model.Step{}
	var err error
	d, err = Step.Insert(d)
	if err != nil {
		t.Fatal(err)
	}
	err = Step.Delete(d)
	if err != nil {
		t.Fatal(err)
	}
	d, err = Step.Get(d.ID)
	if err != nil {
		t.Fatal(err)
	}
	if d.ID != -1 {
		t.Fatal("wanted OID to eq -1, current value = ", d.ID)
	}
}

func TestStepGetByProtestID(t *testing.T) {
	gatherAt := time.Now().Add(time.Hour)
	d := model.Protest{
		ID:     -1,
		Title:  "test",
		Public: true,
		Steps:  []model.Step{model.Step{Place: "place", GatherAt: gatherAt, Details: "details"}},
	}
	var err error
	d, err = Protest.Insert(d)
	if err != nil {
		t.Fatal(err)
	}
	d, err = Step.InsertSteps(d)
	if err != nil {
		t.Fatal(err)
	}
	d, err = Step.GetSteps(d)
	if err != nil {
		t.Fatal(err)
	}

	var steps []model.Step
	steps, err = Step.GetByProtestID(d.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(steps) != 1 {
		t.Fatal("wanted ID to eq 1, current value = ", len(steps))
	}
	if steps[0].ID == -1 {
		t.Fatal("wanted ID to not eq -1, current value = ", steps[0].ID)
	}
	if steps[0].ID != d.Steps[0].ID {
		t.Fatal("wanted ID to eq ", d.Steps[0].ID, ", current value = ", steps[0].ID)
	}
}

func TestStepGetSteps(t *testing.T) {
	gatherAt := time.Now().Add(time.Hour)
	d := model.Protest{
		ID:     -1,
		Title:  "test",
		Public: true,
		Steps:  []model.Step{model.Step{Place: "place", GatherAt: gatherAt, Details: "details"}},
	}
	var err error
	d, err = Protest.Insert(d)
	if err != nil {
		t.Fatal(err)
	}
	d, err = Step.InsertSteps(d)
	if err != nil {
		t.Fatal(err)
	}

	var d2 model.Protest
	d2, err = Step.GetSteps(d)
	if err != nil {
		t.Fatal(err)
	}
	if len(d2.Steps) != 1 {
		t.Fatal("wanted ID to eq 1, current value = ", len(d2.Steps))
	}
	if d2.Steps[0].ID == -1 {
		t.Fatal("wanted ID to not eq -1, current value = ", d2.Steps[0].ID)
	}
}

func TestStepGetProtectedSteps(t *testing.T) {
	gatherAt := time.Now().Add(time.Hour)
	d := model.Protest{
		ID:       -1,
		Title:    "test",
		Public:   false,
		Password: "testpwd",
		Steps:    []model.Step{model.Step{Place: "place", GatherAt: gatherAt, Details: "details"}},
	}
	var err error
	d, err = Protest.Insert(d)
	if err != nil {
		t.Fatal(err)
	}
	d, err = Step.InsertSteps(d)
	if err != nil {
		t.Fatal(err)
	}

	var d2 model.Protest
	d2, err = Step.GetProtectedSteps(d)
	if err != nil {
		t.Fatal(err)
	}
	if len(d2.Steps) != 1 {
		t.Fatal("wanted ID to eq 1, current value = ", len(d2.Steps))
	}
	if d2.Steps[0].ID == -1 {
		t.Fatal("wanted ID to not eq -1, current value = ", d2.Steps[0].ID)
	}
}

func TestStepGetProtectedStepsNotFound(t *testing.T) {
	gatherAt := time.Now().Add(time.Hour)
	d := model.Protest{
		ID:     -1,
		Title:  "test",
		Public: true,
		Steps:  []model.Step{model.Step{Place: "place", GatherAt: gatherAt, Details: "details"}},
	}
	var err error
	d, err = Protest.Insert(d)
	if err != nil {
		t.Fatal(err)
	}
	d, err = Step.InsertSteps(d)
	if err != nil {
		t.Fatal(err)
	}

	var d2 model.Protest
	d2.Public = d.Public
	d2.ID = d.ID
	d2, err = Step.GetProtectedSteps(d2)
	if err != nil {
		t.Fatal(err)
	}
	if len(d2.Steps) != 0 {
		t.Fatal("wanted ID to eq 0, current value = ", len(d2.Steps))
	}
}

//
// func TestProtestSearchLocation(t *testing.T) {
// 	var err error
// 	d := model.Protest{
// 		ID:     -1,
// 		Title:  "testloc",
// 		Public: true,
// 		Steps: []model.Step{
// 			model.Step{Lat: 5.0, Lng: 5.0},
// 		},
// 	}
// 	d, err = Protest.Insert(d)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	d, err = Step.InsertSteps(d)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	var d2 []model.Protest
// 	var atLat, atLng float64
// 	var startDate, endDate *time.Time
// 	var title, protest, organizer string
// 	title = "testloc"
// 	atLat = 1.0
// 	atLng = 1.0
// 	d2, err = Protest.SearchProtests(
// 		&title, &protest, &organizer, startDate, endDate, &atLat, &atLng, 10.0,
// 	)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if len(d2) == 0 {
// 		t.Fatal("wanted len(results) to eq 1, current value = ", len(d2))
// 	}
// 	if d2[0].ID == -1 {
// 		t.Fatal("wanted OID to not eq -1, current value = ", d2[0].ID)
// 	}
// 	if d2[0].Title != "testloc" {
// 		t.Fatal("wanted Title to eq test, current value = ", d2[0].Title)
// 	}
// }
//
// func TestProtestSearchLocationNotFound(t *testing.T) {
// 	var err error
// 	d := model.Protest{
// 		ID:     -1,
// 		Title:  "testloc2",
// 		Public: true,
// 		Steps: []model.Step{
// 			model.Step{Lat: 40.0, Lng: 40.0},
// 		},
// 	}
// 	d, err = Protest.Insert(d)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	d, err = Step.InsertSteps(d)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	log.Println(d)
//
// 	var d2 []model.Protest
// 	var atLat, atLng float64
// 	var startDate, endDate *time.Time
// 	var title, protest, organizer string
// 	title = "testloc2"
// 	atLat = 5.0
// 	atLng = 5.0
// 	d2, err = Protest.SearchProtests(
// 		&title, &protest, &organizer, startDate, endDate, &atLat, &atLng, 20.0,
// 	)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if len(d2) > 0 {
// 		t.Fatal("wanted len(results) to eq 0, current value = ", len(d2))
// 	}
// }
