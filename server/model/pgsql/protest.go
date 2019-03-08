package model

//
// import (
// 	"database/sql"
// 	"fmt"
// 	"time"
//
// 	"github.com/clementauger/monparcours/server/model"
// )
//
// //ProtestService sqlite implementation.
// type ProtestService struct {
// 	DB *sql.DB
// }
//
// //Insert a Protest
// func (s ProtestService) Insert(p model.Protest) (model.Protest, error) {
// 	stmt := `INSERT INTO protest
// 	(author_id,title,protest,description,organizer,public,password,gather_at,created_at)
// 	VALUES
// 	(?,?,?,?,?,?,?,?,?)`
// 	pstmt, err := s.DB.Prepare(stmt)
// 	if err != nil {
// 		return model.Protest{}, err
// 	}
// 	defer pstmt.Close()
// 	createdAt := time.Now()
// 	res, err := pstmt.Exec(
// 		p.AuthorID,
// 		p.Title,
// 		p.Protest,
// 		p.Description,
// 		p.Organizer,
// 		p.Public,
// 		p.Password,
// 		p.GatherAt,
// 		createdAt,
// 	)
// 	if err != nil {
// 		return p, err
// 	}
// 	p.CreatedAt = createdAt
// 	if p.ID, err = res.LastInsertId(); err != nil {
// 		return p, err
// 	}
// 	return p, nil
// }
//
// //Delete a Protest
// func (s ProtestService) Delete(p model.Protest) error {
// 	stmt := `DELETE FROM protest WHERE oid = ?`
// 	pstmt, err := s.DB.Prepare(stmt)
// 	if err != nil {
// 		return err
// 	}
// 	defer pstmt.Close()
// 	_, err = pstmt.Exec(p.ID)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
//
// //Get a Protest by its ID.
// func (s ProtestService) Get(id int64) (model.Protest, error) {
// 	stmt := `SELECT
// 		author_id,title,protest,description,organizer,public,gather_at,created_at,updated_at
// 	 FROM protest
// 	 WHERE oid = ?
// 	 AND IFNULL(deleted_at, '') = ''
// 	 AND password = ""
// 	 AND public = 1
// 	 `
// 	pstmt, err := s.DB.Prepare(stmt)
// 	if err != nil {
// 		return model.Protest{}, err
// 	}
// 	defer pstmt.Close()
// 	rows, err := pstmt.Query(fmt.Sprint(id))
// 	if err != nil {
// 		return model.Protest{}, err
// 	}
// 	var p model.Protest
// 	p.ID = -1
// 	defer rows.Close()
// 	var found bool
// 	for rows.Next() {
// 		err = rows.Scan(
// 			&p.AuthorID,
// 			&p.Title,
// 			&p.Protest,
// 			&p.Description,
// 			&p.Organizer,
// 			&p.Public,
// 			&p.GatherAt,
// 			&p.CreatedAt,
// 			&p.UpdatedAt,
// 		)
// 		if err != nil {
// 			return p, err
// 		}
// 		found = true
// 	}
// 	if found {
// 		p.ID = id
// 	}
// 	return p, rows.Err()
// }
//
// //GetWithPassword a Protest by its ID and password.
// func (s ProtestService) GetWithPassword(id int64, pwd string) (model.Protest, error) {
// 	stmt := `SELECT
// 		author_id,title,protest,description,organizer,gather_at,created_at,updated_at
// 	 FROM protest
// 	 WHERE oid = ?
// 	 AND password = ?
// 	 AND IFNULL(deleted_at, '') = ''
// 	 `
// 	pstmt, err := s.DB.Prepare(stmt)
// 	if err != nil {
// 		return model.Protest{}, err
// 	}
// 	defer pstmt.Close()
// 	rows, err := pstmt.Query(fmt.Sprint(id), pwd)
// 	if err != nil {
// 		return model.Protest{}, err
// 	}
// 	var p model.Protest
// 	p.ID = -1
// 	defer rows.Close()
// 	var found bool
// 	for rows.Next() {
// 		err = rows.Scan(
// 			&p.AuthorID,
// 			&p.Title,
// 			&p.Protest,
// 			&p.Description,
// 			&p.Organizer,
// 			&p.GatherAt,
// 			&p.CreatedAt,
// 			&p.UpdatedAt,
// 		)
// 		if err != nil {
// 			return p, err
// 		}
// 		found = true
// 	}
// 	if found {
// 		p.ID = id
// 	}
// 	return p, rows.Err()
// }
//
// //GetByAuthorID protests by theirs author IDs.
// func (s ProtestService) GetByAuthorID(authorID string) ([]model.Protest, error) {
// 	out := []model.Protest{}
// 	stmt := `SELECT
// 			oid,title,protest,description,organizer,public,gather_at,created_at,updated_at
// 	 FROM protest
// 	 WHERE author_id = ?
// 	 AND IFNULL(deleted_at, '') = ''
// 	 ORDER BY gather_at DESC
// 	 `
// 	pstmt, err := s.DB.Prepare(stmt)
// 	if err != nil {
// 		return out, err
// 	}
// 	defer pstmt.Close()
// 	rows, err := pstmt.Query(authorID)
// 	if err != nil {
// 		return out, err
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		var p model.Protest
// 		err = rows.Scan(
// 			&p.ID,
// 			&p.Title,
// 			&p.Protest,
// 			&p.Description,
// 			&p.Organizer,
// 			&p.Public,
// 			&p.GatherAt,
// 			&p.CreatedAt,
// 			&p.UpdatedAt,
// 		)
// 		if err != nil {
// 			return out, err
// 		}
// 		p.AuthorID = authorID
// 		if !p.Public {
// 			p.Title = p.Title[:3] + "***"
// 			p.Protest = "***"
// 			p.Description = "***"
// 			p.Organizer = "***"
// 			p.GatherAt = time.Now().Add(time.Hour * 24 * 365 * 10 * -1)
// 			p.CreatedAt = time.Now().Add(time.Hour * 24 * 365 * 10 * -1)
// 			t := time.Now().Add(time.Hour * 24 * 365 * 10 * -1)
// 			p.UpdatedAt = &t
// 		}
// 		out = append(out, p)
// 	}
// 	return out, rows.Err()
// }
//
// //GetByIDs protests by theirs IDs.
// func (s ProtestService) GetByIDs(ids ...int64) ([]model.Protest, error) {
// 	out := make([]model.Protest, len(ids))
// 	var err error
// 	for _, id := range ids {
// 		var p model.Protest
// 		p, err = s.Get(id)
// 		if err != nil {
// 			return out, err
// 		}
// 		out = append(out, p)
// 	}
// 	return out, err
// }
//
// // SearchProtests given position and given date.
// func (s ProtestService) SearchProtests(
// 	title, protest, organizer *string,
// 	startDate, endDate *time.Time,
// 	atLat, atLng *float64,
// 	ray float64,
// ) ([]model.Protest, error) {
// 	out := []model.Protest{}
// 	args := []interface{}{}
// 	stmt := `SELECT
// 			oid,title,protest,description,organizer,gather_at,created_at,updated_at
// 		 FROM protest
// 		 WHERE 1=1
// 		 AND IFNULL(deleted_at, '') = ''
// 		 AND password = ""
// 		 AND public = 1
// 		`
// 	if atLat != nil && atLng != nil {
// 		stmt = `SELECT
// 			protest.oid,protest.title,protest.protest,protest.description,protest.organizer,protest.gather_at,protest.created_at,protest.updated_at,
// 			( ACOS( COS( RADIANS( ?  ) )
// 							* COS( RADIANS( lat ) )
// 							* COS( RADIANS( lng ) - RADIANS( ? ) )
// 							+ SIN( RADIANS( ?   ) )
// 							* SIN( RADIANS( lat ) )
// 					)
// 				* 6371
// 				) AS distance_in_km
// 		 FROM protest
// 		 INNER JOIN step ON (protest.oid=step.protest_id)
// 		 WHERE 1=1
// 		 AND IFNULL(protest.deleted_at, '') = ''
// 		 AND protest.password = ""
// 		 AND protest.public = 1
// 		`
// 		// qstmt := `SELECT
// 		// 	 protest_id,
// 		// 	 ( ACOS( COS( RADIANS( ?  ) )
// 		// 					 * COS( RADIANS( lat ) )
// 		// 					 * COS( RADIANS( lng ) - RADIANS( ? ) )
// 		// 					 + SIN( RADIANS( ?   ) )
// 		// 					 * SIN( RADIANS( lat ) )
// 		// 			 )
// 		// 		 * 6371
// 		// 		 ) AS distance_in_km
// 		// 	FROM step
// 		// 	WHERE 1=1
// 		// `
// 		args = append(args, *atLat, *atLng, *atLat)
// 		// if startDate != nil {
// 		// 	stmt += `AND step.gather_at > ?
// 		// 		`
// 		// 	args = append(args, *startDate)
// 		// }
// 		// if endDate != nil {
// 		// 	stmt += `AND step.gather_at < ?
// 		// 		`
// 		// 	args = append(args, *endDate)
// 		// }
// 		// qstmt += `
// 		// 	HAVING distance_in_km < 5
// 		// 	ORDER BY distance_in_km ASC
// 		// `
// 		// stmt += `AND oid IN (` + qstmt + `)
// 		// `
// 	}
// 	if title != nil && *title != "" {
// 		stmt += `AND protest.title LIKE ?
// 			`
// 		args = append(args, *title)
// 	}
// 	if protest != nil && *protest != "" {
// 		stmt += `AND protest.protest LIKE ?
// 			`
// 		args = append(args, *protest)
// 	}
// 	if organizer != nil && *organizer != "" {
// 		stmt += `AND protest.organizer LIKE ?
// 			`
// 		args = append(args, *organizer)
// 	}
// 	if startDate != nil {
// 		stmt += `AND protest.gather_at > ?
// 			`
// 		args = append(args, *startDate)
// 	}
// 	if endDate != nil {
// 		stmt += `AND protest.gather_at < ?
// 			`
// 		args = append(args, *endDate)
// 	}
// 	if atLat != nil && atLng != nil {
// 		stmt += fmt.Sprintf(`HAVING distance_in_km < %v
// 	ORDER BY distance_in_km ASC
// `, ray)
// 	}
// 	stmt += `LIMIT 100
// 			`
// 	pstmt, err := s.DB.Prepare(stmt)
// 	if err != nil {
// 		return out, err
// 	}
// 	defer pstmt.Close()
// 	rows, err := pstmt.Query(args...)
// 	if err != nil {
// 		return out, err
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		var p model.Protest
// 		if atLat != nil && atLng != nil {
// 			var dist float64
// 			if err = rows.Scan(&p.ID, &p.Title, &p.Protest, &p.Description, &p.Organizer, &p.GatherAt, &p.CreatedAt, &p.UpdatedAt, &dist); err != nil {
// 				return out, err
// 			}
// 		} else {
// 			if err = rows.Scan(&p.ID, &p.Title, &p.Protest, &p.Description, &p.Organizer, &p.GatherAt, &p.CreatedAt, &p.UpdatedAt); err != nil {
// 				return out, err
// 			}
// 		}
//
// 		out = append(out, p)
// 	}
// 	return out, rows.Err()
//
// }
//
// //StepService sqlite implementation
// type StepService struct {
// 	DB *sql.DB
// }
//
// //InsertSteps of a Protest
// func (s StepService) InsertSteps(p model.Protest) (model.Protest, error) {
// 	for i, v := range p.Steps {
// 		v.ProtestID = p.ID
// 		p.Steps[i] = v
// 	}
// 	var err error
// 	p.Steps, err = s.InsertAll(p.Steps)
// 	return p, err
// }
//
// //InsertAll steps
// func (s StepService) InsertAll(steps []model.Step) ([]model.Step, error) {
// 	for i, v := range steps {
// 		t, err := s.Insert(v)
// 		if err != nil {
// 			return steps, err
// 		}
// 		steps[i] = t
// 	}
// 	return steps, nil
// }
//
// //Insert a Step
// func (s StepService) Insert(step model.Step) (model.Step, error) {
// 	stmt := `INSERT INTO step
// 	(protest_id, place, details, gather_at, lat, lng)
// 	VALUES
// 	(?,?,?,?,?,?)`
// 	pstmt, err := s.DB.Prepare(stmt)
// 	if err != nil {
// 		return model.Step{}, err
// 	}
// 	defer pstmt.Close()
// 	res, err := pstmt.Exec(
// 		step.ProtestID,
// 		step.Place,
// 		step.Details,
// 		step.GatherAt,
// 		step.Lat,
// 		step.Lng,
// 	)
// 	if err != nil {
// 		return step, err
// 	}
// 	if step.ID, err = res.LastInsertId(); err != nil {
// 		return step, err
// 	}
// 	return step, nil
// }
//
// //Delete a Step.
// func (s StepService) Delete(p model.Step) error {
// 	stmt := `DELETE FROM step WHERE oid = ?`
// 	pstmt, err := s.DB.Prepare(stmt)
// 	if err != nil {
// 		return err
// 	}
// 	defer pstmt.Close()
// 	_, err = pstmt.Exec(p.ID)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
//
// //Get a Step by its ID.
// func (s StepService) Get(id int64) (model.Step, error) {
// 	stmt := `SELECT
// 		protest_id, place, details, gather_at, lat, lng
// 	 FROM step WHERE oid = ?`
// 	pstmt, err := s.DB.Prepare(stmt)
// 	if err != nil {
// 		return model.Step{}, err
// 	}
// 	defer pstmt.Close()
// 	rows, err := pstmt.Query(id)
// 	if err != nil {
// 		return model.Step{}, err
// 	}
// 	var p model.Step
// 	defer rows.Close()
// 	for rows.Next() {
// 		if err = rows.Scan(&p.ProtestID, &p.Place, &p.Details, &p.GatherAt, &p.Lat, &p.Lng); err != nil {
// 			return p, err
// 		}
// 		p.ID = id
// 	}
// 	return p, rows.Err()
// }
//
// //GetAll steps by theirs protest IDs.
// func (s StepService) GetAll(protestID int64) ([]model.Step, error) {
// 	out := []model.Step{}
// 	stmt := `SELECT
// 		oid, place, details, gather_at, lat, lng
// 	 FROM step WHERE protest_id = ?`
// 	pstmt, err := s.DB.Prepare(stmt)
// 	if err != nil {
// 		return out, err
// 	}
// 	defer pstmt.Close()
// 	rows, err := pstmt.Query(protestID)
// 	if err != nil {
// 		return out, err
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		var p model.Step
// 		if err = rows.Scan(&p.ID, &p.Place, &p.Details, &p.GatherAt, &p.Lat, &p.Lng); err != nil {
// 			return out, err
// 		}
// 		p.ProtestID = protestID
// 		out = append(out, p)
// 	}
// 	return out, rows.Err()
// }
//
// // GetSteps of a Protest.
// func (s StepService) GetSteps(p model.Protest) (model.Protest, error) {
// 	var err error
// 	if p.Public {
// 		p.Steps, err = s.GetAll(p.ID)
// 	}
// 	return p, err
// }
//
// // GetProtectedSteps of a Protest.
// func (s StepService) GetProtectedSteps(p model.Protest) (model.Protest, error) {
// 	var err error
// 	if !p.Public {
// 		p.Steps, err = s.GetAll(p.ID)
// 	}
// 	return p, err
// }
//
// // FindStepsAround given position and given date.
// func (s StepService) FindStepsAround(
// 	atDate time.Time,
// 	withinTime time.Duration,
// 	atLat, atLng float64,
// ) ([]model.Step, error) {
// 	out := []model.Step{}
// 	stmt := `SELECT
// 			oid, place, details, gather_at, lat, lng
// 			, ( ACOS( COS( RADIANS( ?  ) )
// 							* COS( RADIANS( lat ) )
// 							* COS( RADIANS( lng ) - RADIANS( ? ) )
// 							+ SIN( RADIANS( ?   ) )
// 							* SIN( RADIANS( lat ) )
// 					)
// 				* 6371
// 				) AS distance_in_km
// 		 FROM step
// 		 WHERE gather_at > ?
// 		 AND gather_at < ?
// 	 	ORDER BY distance_in_km ASC
// 	 	LIMIT 100
// 		`
// 	pstmt, err := s.DB.Prepare(stmt)
// 	if err != nil {
// 		return out, err
// 	}
// 	defer pstmt.Close()
// 	rows, err := pstmt.Query(atLat, atLng, atLat, atDate.Add(-1*withinTime), atDate.Add(withinTime))
// 	if err != nil {
// 		return out, err
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		var p model.Step
// 		if err = rows.Scan(&p.ID, &p.Place, &p.Details, &p.GatherAt, &p.Lat, &p.Lng); err != nil {
// 			return out, err
// 		}
// 		out = append(out, p)
// 	}
// 	return out, rows.Err()
// 	/*
// 		SELECT m.school_id
// 				, m.location_id
// 				, m.school_name
// 				, m.lat
// 				, m.lng
//
// 				, ( ACOS( COS( RADIANS( @lat  ) )
// 								* COS( RADIANS( m.lat ) )
// 								* COS( RADIANS( m.lng ) - RADIANS( @lng ) )
// 								+ SIN( RADIANS( @lat  ) )
// 								* SIN( RADIANS( m.lat ) )
// 						)
// 					* 6371
// 					) AS distance_in_km
//
// 		FROM mytable m
// 		ORDER BY distance_in_km ASC
// 		LIMIT 100
// 	*/
// }
