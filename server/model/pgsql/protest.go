package pgsql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/clementauger/monparcours/server/model"
	"github.com/jmoiron/sqlx"
)

//ProtestService sqlite implementation.
type ProtestService struct {
	DB *sql.DB
}

//Insert a Protest
func (s ProtestService) Insert(p model.Protest) (model.Protest, error) {
	stmt := `INSERT INTO protest
	(author_id,title,protest,description,organizer,public,password,gather_at,created_at)
	VALUES
	(?,?,?,?,?,?,?,?,?)
	RETURNING oid `
	stmt = sqlx.Rebind(sqlx.BindType("postgres"), stmt)
	pstmt, err := s.DB.Prepare(stmt)
	if err != nil {
		return model.Protest{}, err
	}
	defer pstmt.Close()
	createdAt := time.Now()
	row := pstmt.QueryRow(
		p.AuthorID,
		p.Title,
		p.Protest,
		p.Description,
		p.Organizer,
		p.Public,
		p.Password,
		p.GatherAt,
		createdAt,
	)
	err = row.Scan(&p.ID)
	if err != nil {
		return p, err
	}
	p.CreatedAt = createdAt
	return p, nil
}

//Delete a Protest
func (s ProtestService) Delete(p model.Protest) error {
	stmt := `UPDATE protest SET deleted_at=? WHERE oid = ?`
	stmt = sqlx.Rebind(sqlx.BindType("postgres"), stmt)
	pstmt, err := s.DB.Prepare(stmt)
	if err != nil {
		return err
	}
	defer pstmt.Close()
	_, err = pstmt.Exec(time.Now(), p.ID)
	if err != nil {
		return err
	}
	return nil
}

//Get a Protest by its ID.
func (s ProtestService) Get(id int64) (model.Protest, error) {
	stmt := `SELECT
		oid,author_id,title,protest,description,organizer,public,gather_at,created_at,updated_at
	 FROM protest
	 WHERE oid = ?
	 AND COALESCE(deleted_at::TEXT, '') = ''
	 AND password = ''
	 AND public = 1
	 `
	stmt = sqlx.Rebind(sqlx.BindType("postgres"), stmt)
	pstmt, err := s.DB.Prepare(stmt)
	if err != nil {
		return model.Protest{}, err
	}
	defer pstmt.Close()
	rows, err := pstmt.Query(fmt.Sprint(id))
	if err != nil {
		return model.Protest{}, err
	}
	var p model.Protest
	p.ID = -1
	defer rows.Close()
	var found bool
	for rows.Next() {
		err = rows.Scan(
			&p.ID,
			&p.AuthorID,
			&p.Title,
			&p.Protest,
			&p.Description,
			&p.Organizer,
			&p.Public,
			&p.GatherAt,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return p, err
		}
		found = true
	}
	if found {
		p.ID = id
	}
	return p, rows.Err()
}

//GetWithPassword a Protest by its ID and password.
func (s ProtestService) GetWithPassword(id int64, pwd string) (model.Protest, error) {
	stmt := `SELECT
		author_id,title,protest,description,organizer,public,gather_at,created_at,updated_at
	 FROM protest
	 WHERE oid = ?
	 AND password = ?
	 AND COALESCE(deleted_at::TEXT, '') = ''
	 `
	stmt = sqlx.Rebind(sqlx.BindType("postgres"), stmt)
	pstmt, err := s.DB.Prepare(stmt)
	if err != nil {
		return model.Protest{}, err
	}
	defer pstmt.Close()
	rows, err := pstmt.Query(fmt.Sprint(id), pwd)
	if err != nil {
		return model.Protest{}, err
	}
	var p model.Protest
	p.ID = -1
	defer rows.Close()
	var found bool
	for rows.Next() {
		err = rows.Scan(
			&p.AuthorID,
			&p.Title,
			&p.Protest,
			&p.Description,
			&p.Organizer,
			&p.Public,
			&p.GatherAt,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return p, err
		}
		found = true
	}
	if found {
		p.ID = id
	}
	return p, rows.Err()
}

//GetByAuthorID protests by theirs author IDs.
func (s ProtestService) GetByAuthorID(authorID string) ([]model.Protest, error) {
	out := []model.Protest{}
	stmt := `SELECT
			oid,title,protest,description,organizer,public,gather_at,created_at,updated_at
	 FROM protest
	 WHERE author_id = ?
	 AND COALESCE(deleted_at::TEXT, '') = ''
	 ORDER BY gather_at DESC
	 `
	stmt = sqlx.Rebind(sqlx.BindType("postgres"), stmt)
	pstmt, err := s.DB.Prepare(stmt)
	if err != nil {
		return out, err
	}
	defer pstmt.Close()
	rows, err := pstmt.Query(authorID)
	if err != nil {
		return out, err
	}
	defer rows.Close()
	for rows.Next() {
		var p model.Protest
		err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.Protest,
			&p.Description,
			&p.Organizer,
			&p.Public,
			&p.GatherAt,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return out, err
		}
		p.AuthorID = authorID
		if !p.Public {
			var s = "***"
			p.Title = p.Title[:1] + s
			p.Protest = s
			p.Description = s
			p.Organizer = s
			p.GatherAt = time.Now().Add(time.Hour * 24 * 365 * 10 * -1)
			p.CreatedAt = time.Now().Add(time.Hour * 24 * 365 * 10 * -1)
			t := time.Now().Add(time.Hour * 24 * 365 * 10 * -1)
			p.UpdatedAt = &t
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

//GetByIDs protests by theirs IDs.
func (s ProtestService) GetByIDs(ids ...int64) ([]model.Protest, error) {
	out := make([]model.Protest, 0, len(ids))
	var err error
	for _, id := range ids {
		var p model.Protest
		p, err = s.Get(id)
		if err != nil {
			return out, err
		}
		out = append(out, p)
	}
	return out, err
}

// SearchProtests given position and given date.
func (s ProtestService) SearchProtests(
	title, protest, organizer *string,
	startDate, endDate *time.Time,
	atLat, atLng *float64,
	ray float64,
) ([]model.Protest, error) {
	out := []model.Protest{}
	args := []interface{}{}
	stmt := `SELECT
			oid,title,protest,description,organizer,public,gather_at,created_at,updated_at
		 FROM protest
		 WHERE 1=1
		 AND COALESCE(deleted_at::TEXT, '') = ''
		 AND password = ''
		 AND public = 1
		`
	if atLat != nil && atLng != nil {
		stmt = `SELECT
			protest.oid,
			protest.title,
			protest.protest,
			protest.description,
			protest.organizer,
			protest.public,
			protest.gather_at,
			protest.created_at,
			protest.updated_at
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
		 AND COALESCE(protest.deleted_at::TEXT, '') = ''
		 AND protest.password = ''
		 AND protest.public = 1
		`
		stmt += fmt.Sprintf(`AND s.distance < ? `)
		args = append(args, *atLat, *atLng, *atLat, ray)
	}
	if title != nil && *title != "" {
		stmt += `AND protest.title LIKE ? `
		args = append(args, *title)
	}
	if protest != nil && *protest != "" {
		stmt += `AND protest.protest LIKE ? `
		args = append(args, *protest)
	}
	if organizer != nil && *organizer != "" {
		stmt += `AND protest.organizer LIKE ? `
		args = append(args, *organizer)
	}
	if startDate != nil {
		stmt += `AND protest.gather_at > ? `
		args = append(args, *startDate)
	}
	if endDate != nil {
		stmt += `AND protest.gather_at <= ? `
		args = append(args, *endDate)
	}
	if atLat != nil && atLng != nil {
		stmt += fmt.Sprintf(`
	GROUP BY protest.oid, s.distance
	ORDER BY s.distance ASC `)
	}
	stmt += `LIMIT 100;`
	stmt = sqlx.Rebind(sqlx.BindType("postgres"), stmt)
	// log.Println(stmt)
	// log.Println(args)
	pstmt, err := s.DB.Prepare(stmt)
	if err != nil {
		return out, err
	}
	defer pstmt.Close()
	rows, err := pstmt.Query(args...)
	if err != nil {
		return out, err
	}
	defer rows.Close()
	for rows.Next() {
		var p model.Protest
		if err = rows.Scan(
			&p.ID, &p.Title, &p.Protest, &p.Description,
			&p.Organizer, &p.Public,
			&p.GatherAt, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return out, err
		}
		out = append(out, p)
	}
	return out, rows.Err()

}

//StepService sqlite implementation
type StepService struct {
	DB *sql.DB
}

//InsertSteps of a Protest
func (s StepService) InsertSteps(p model.Protest) (model.Protest, error) {
	for i, v := range p.Steps {
		v.ProtestID = p.ID
		p.Steps[i] = v
	}
	var err error
	p.Steps, err = s.InsertAll(p.Steps)
	return p, err
}

//InsertAll steps
func (s StepService) InsertAll(steps []model.Step) ([]model.Step, error) {
	for i, v := range steps {
		t, err := s.Insert(v)
		if err != nil {
			return steps, err
		}
		steps[i] = t
	}
	return steps, nil
}

//Insert a Step
func (s StepService) Insert(step model.Step) (model.Step, error) {
	stmt := `INSERT INTO step
	(protest_id, place, details, gather_at, lat, lng)
	VALUES
	(?,?,?,?,?,?)
	RETURNING oid`
	stmt = sqlx.Rebind(sqlx.BindType("postgres"), stmt)
	pstmt, err := s.DB.Prepare(stmt)
	if err != nil {
		return model.Step{}, err
	}
	defer pstmt.Close()

	row := pstmt.QueryRow(
		step.ProtestID,
		step.Place,
		step.Details,
		step.GatherAt,
		step.Lat,
		step.Lng,
	)
	err = row.Scan(&step.ID)
	if err != nil {
		return step, err
	}
	return step, nil
}

//Delete a Step.
func (s StepService) Delete(p model.Step) error {
	stmt := `DELETE FROM step WHERE oid = ?`
	stmt = sqlx.Rebind(sqlx.BindType("postgres"), stmt)
	pstmt, err := s.DB.Prepare(stmt)
	if err != nil {
		return err
	}
	defer pstmt.Close()
	_, err = pstmt.Exec(p.ID)
	if err != nil {
		return err
	}
	return nil
}

//Get a Step by its ID.
func (s StepService) Get(id int64) (model.Step, error) {
	stmt := `SELECT
		oid,protest_id, place, details, gather_at::TIMESTAMP, lat, lng
	 FROM step WHERE oid = ?`
	stmt = sqlx.Rebind(sqlx.BindType("postgres"), stmt)
	pstmt, err := s.DB.Prepare(stmt)
	if err != nil {
		return model.Step{}, err
	}
	defer pstmt.Close()
	rows, err := pstmt.Query(id)
	if err != nil {
		return model.Step{}, err
	}
	var p model.Step
	p.ID = -1
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(
			&p.ID, &p.ProtestID, &p.Place, &p.Details, &p.GatherAt, &p.Lat, &p.Lng,
		); err != nil {
			return p, err
		}
		p.ID = id
	}
	return p, rows.Err()
}

//GetByProtestID steps by theirs protest IDs.
func (s StepService) GetByProtestID(protestID int64) ([]model.Step, error) {
	out := []model.Step{}
	stmt := `SELECT
		oid, place, details, gather_at, lat, lng
	 FROM step WHERE protest_id = ?`
	stmt = sqlx.Rebind(sqlx.BindType("postgres"), stmt)
	pstmt, err := s.DB.Prepare(stmt)
	if err != nil {
		return out, err
	}
	defer pstmt.Close()
	rows, err := pstmt.Query(protestID)
	if err != nil {
		return out, err
	}
	defer rows.Close()
	for rows.Next() {
		var p model.Step
		if err = rows.Scan(&p.ID, &p.Place, &p.Details, &p.GatherAt, &p.Lat, &p.Lng); err != nil {
			return out, err
		}
		p.ProtestID = protestID
		out = append(out, p)
	}
	return out, rows.Err()
}

// GetSteps of a Protest.
func (s StepService) GetSteps(p model.Protest) (model.Protest, error) {
	var err error
	if p.Public {
		p.Steps, err = s.GetByProtestID(p.ID)
	}
	return p, err
}

// GetProtectedSteps of a Protest.
func (s StepService) GetProtectedSteps(p model.Protest) (model.Protest, error) {
	var err error
	if !p.Public {
		p.Steps, err = s.GetByProtestID(p.ID)
	}
	return p, err
}

// FindStepsAround given position and given date.
func (s StepService) FindStepsAround(
	atDate time.Time,
	withinTime time.Duration,
	atLat, atLng float64,
) ([]model.Step, error) {
	out := []model.Step{}
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
		 AND gather_at <= ?
	 	ORDER BY distance_in_km ASC
	 	LIMIT 100
	`
	stmt = sqlx.Rebind(sqlx.BindType("postgres"), stmt)
	pstmt, err := s.DB.Prepare(stmt)
	if err != nil {
		return out, err
	}
	defer pstmt.Close()
	rows, err := pstmt.Query(atLat, atLng, atLat, atDate.Add(-1*withinTime), atDate.Add(withinTime))
	if err != nil {
		return out, err
	}
	defer rows.Close()
	for rows.Next() {
		var p model.Step
		if err = rows.Scan(&p.ID, &p.Place, &p.Details, &p.GatherAt, &p.Lat, &p.Lng); err != nil {
			return out, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}
