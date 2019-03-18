package cruded

import (
	"fmt"
	"time"

	"github.com/clementauger/monparcours/server/model"
	"github.com/clementauger/crud"
)

//ProtestService sqlite implementation.
type ProtestService struct {
	Crud    *crud.Crudable
	Dialect string
}

//Insert a Protest
func (s *ProtestService) Insert(p model.Protest) (model.Protest, error) {
	err := s.Crud.Insert(&p)
	return p, err
}

//Delete a Protest
func (s *ProtestService) Delete(p model.Protest) error {
	return s.Crud.Delete(&p)
}

//Get a ContactMessage by its ID.
func (s *ProtestService) Get(id int64) (p model.Protest, err error) {
	p.ID = -1
	stmt := `AND oid = ? AND password = '' AND public = 1 `
	err = s.Crud.Fetch(&p, stmt, id)
	return p, err
}

//GetWithPassword a Protest by its ID and password.
func (s *ProtestService) GetWithPassword(id int64, pwd string) (p model.Protest, err error) {
	p.ID = -1
	stmt := `AND oid = ? AND password = ?  `
	err = s.Crud.Fetch(&p, stmt, id, pwd)
	return p, err
}

//GetByAuthorID protests by theirs author IDs.
func (s *ProtestService) GetByAuthorID(authorID string) (out []model.Protest, err error) {
	stmt := `AND author_id = ? ORDER BY gather_at DESC `
	err = s.Crud.Fetch(&out, stmt, authorID)
	for i, p := range out {
		if !p.Public {
			p.Title = p.Title[:3] + "***"
			p.Protest = "***"
			p.Description = "***"
			p.Organizer = "***"
			p.GatherAt = time.Now().Add(time.Hour * 24 * 365 * 10 * -1)
			p.CreatedAt = time.Now().Add(time.Hour * 24 * 365 * 10 * -1)
			t := time.Now().Add(time.Hour * 24 * 365 * 10 * -1)
			p.UpdatedAt = &t
			out[i] = p
		}
	}
	return out, err
}

//GetAll ContactMessage.
func (s *ProtestService) GetAll() (out []model.Protest, err error) {
	stmt := `ORDER BY created_at DESC
   LIMIT 50
   `
	err = s.Crud.Fetch(&out, stmt)
	return out, err
}

//GetByIDs protests by theirs IDs.
func (s *ProtestService) GetByIDs(ids ...int64) (out []model.Protest, err error) {
	err = s.Crud.FetchByIDs(&out, ids)
	return out, err
}

// SearchProtests given position and given date.
func (s *ProtestService) SearchProtests(
	title, protest, organizer *string,
	startDate, endDate *time.Time,
	atLat, atLng *float64,
	ray float64,
) (out []model.Protest, err error) {

	args := []interface{}{}
	stmt := `SELECT
			oid,title,protest,description,organizer,gather_at,created_at,updated_at
		 FROM protest
		 WHERE 1=1
		 AND password = '' `
	if s.Dialect == "postgres" {
		stmt += `AND public = 1 `
		stmt += `AND COALESCE(deleted_at::TEXT, '') = '' `
	} else {
		stmt += `AND public = 1 `
		stmt += `AND IFNULL(deleted_at, '') = '' `
	}
	if atLat != nil && atLng != nil {
		stmt = `SELECT
			protest.*
		 FROM protest
		 INNER JOIN (
			 SELECT step.*,
		 			( ACOS( COS( RADIANS( ?  ) )
		 							* COS( RADIANS( step.lat ) )
		 							* COS( RADIANS( step.lng ) - RADIANS( ? ) )
		 							+ SIN( RADIANS( ?   ) )
		 							* SIN( RADIANS( step.lat ) )
		 					)
		 				* 6371
		 			)
		 			AS distance
			 FROM step
			 ) s ON (protest.oid=s.protest_id)
		 WHERE 1=1
		 AND protest.password = '' `
		stmt += `AND protest.public = 1 `
		if s.Dialect == "postgres" {
			stmt += `AND COALESCE(deleted_at::TEXT, '') = '' `
		} else {
			stmt += `AND IFNULL(deleted_at, '') = '' `
		}
		stmt += fmt.Sprintf(`AND s.distance < ? `)
		args = append(args, *atLat, *atLng, *atLat, ray)
	}
	if title != nil && *title != "" {
		stmt += `AND protest.title LIKE ?
			`
		args = append(args, *title)
	}
	if protest != nil && *protest != "" {
		stmt += `AND protest.protest LIKE ?
			`
		args = append(args, *protest)
	}
	if organizer != nil && *organizer != "" {
		stmt += `AND protest.organizer LIKE ?
			`
		args = append(args, *organizer)
	}
	if startDate != nil {
		stmt += `AND protest.gather_at > ?
			`
		args = append(args, *startDate)
	}
	if endDate != nil {
		stmt += `AND protest.gather_at < ?
			`
		args = append(args, *endDate)
	}
	if atLat != nil && atLng != nil {
		stmt += fmt.Sprintf(`
			GROUP BY protest.oid, s.distance
			ORDER BY s.distance ASC `)
	}
	stmt += `LIMIT 100;`

	err = s.Crud.Read(&out, stmt, args...)
	return
}

//StepService sqlite implementation
type StepService struct {
	Crud *crud.Crudable
}

//InsertSteps of a Protest
func (s *StepService) InsertSteps(p model.Protest) (_ model.Protest, err error) {
	for i, v := range p.Steps {
		v.ProtestID = p.ID
		p.Steps[i] = v
	}
	p.Steps, err = s.InsertAll(p.Steps)
	return p, err
}

//InsertAll steps
func (s *StepService) InsertAll(steps []model.Step) ([]model.Step, error) {
	var err error
	err = s.Crud.Insert(&steps)
	return steps, err
}

//Insert a Step
func (s *StepService) Insert(step model.Step) (model.Step, error) {
	var err error
	err = s.Crud.Insert(&step)
	return step, err
}

//Delete a Step.
func (s *StepService) Delete(p model.Step) error {
	return s.Crud.Delete(p)
}

//Get a Step by its ID.
func (s *StepService) Get(id int64) (out model.Step, err error) {
	out.ID = -1
	err = s.Crud.FetchByID(&out, id)
	return out, err
}

//GetByProtestID steps by theirs protest IDs.
func (s *StepService) GetByProtestID(protestID int64) (out []model.Step, err error) {
	err = s.Crud.Fetch(&out, `AND protest_id = ?`, protestID)
	return out, err
}

// GetSteps of a Protest.
func (s *StepService) GetSteps(p model.Protest) (model.Protest, error) {
	var err error
	if p.Public {
		p.Steps, err = s.GetByProtestID(p.ID)
	}
	return p, err
}

// GetProtectedSteps of a Protest.
func (s *StepService) GetProtectedSteps(p model.Protest) (model.Protest, error) {
	var err error
	if !p.Public {
		p.Steps, err = s.GetByProtestID(p.ID)
	}
	return p, err
}

// FindStepsAround given position and given date.
func (s *StepService) FindStepsAround(
	atDate time.Time,
	withinTime time.Duration,
	atLat, atLng float64,
) (out []model.Step, err error) {

	stmt := `SELECT
      oid, place, details, gather_at, lat, lng
      , ( ACOS( COS( RADIANS( ?  ) )
              * COS( RADIANS( lat ) )
              * COS( RADIANS( lng ) - RADIANS( ? ) )
              + SIN( RADIANS( ?   ) )
              * SIN( RADIANS( lat ) )
          )
        * 6371
        ) AS distance_in_km
     FROM step
     WHERE gather_at > ?
     AND gather_at < ?
    ORDER BY distance_in_km ASC
    LIMIT 100;`

	err = s.Crud.Read(&out, stmt, atLat, atLng, atLat, atDate.Add(-1*withinTime), atDate.Add(withinTime))
	return out, err
}
