package pgsql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/clementauger/monparcours/server/model"
	"github.com/jmoiron/sqlx"
)

//ContactMessageService sqlite implementation.
type ContactMessageService struct {
	DB *sql.DB
}

//Insert a ContactMessage
func (s ContactMessageService) Insert(p model.ContactMessage) (model.ContactMessage, error) {
	stmt := `INSERT INTO contact_message
	(returnaddr,subject,body,created_at)
	VALUES
	(?,?,?,?)
	RETURNING oid `
	stmt = sqlx.Rebind(sqlx.BindType("postgres"), stmt)
	pstmt, err := s.DB.Prepare(stmt)
	if err != nil {
		return model.ContactMessage{}, err
	}
	defer pstmt.Close()
	createdAt := time.Now()
	row := pstmt.QueryRow(p.ReturnAddr, p.Subject, p.Body, createdAt)
	err = row.Scan(&p.ID)
	if err != nil {
		return p, err
	}
	p.CreatedAt = createdAt
	return p, nil
}

//Delete a ContactMessage
func (s ContactMessageService) Delete(p model.ContactMessage) error {
	stmt := `UPDATE contact_message SET deleted_at=? WHERE oid = ?`
	stmt = sqlx.Rebind(sqlx.BindType("postgres"), stmt)
	pstmt, err := s.DB.Prepare(stmt)
	if err != nil {
		return err
	}
	defer pstmt.Close()
	_, err = pstmt.Exec(time.Now(), p.ID)
	return err
}

//Get a ContactMessage by its ID.
func (s ContactMessageService) Get(id int64) (model.ContactMessage, error) {
	stmt := `SELECT
		returnaddr, subject, body, created_at, updated_at, deleted_at
	 FROM contact_message
   WHERE oid = ?
	 AND COALESCE(deleted_at::TEXT, '') = ''
	 `
	stmt = sqlx.Rebind(sqlx.BindType("postgres"), stmt)
	pstmt, err := s.DB.Prepare(stmt)
	if err != nil {
		return model.ContactMessage{}, err
	}
	defer pstmt.Close()
	rows, err := pstmt.Query(fmt.Sprint(id))
	if err != nil {
		return model.ContactMessage{}, err
	}
	var p model.ContactMessage
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(
			&p.ReturnAddr,
			&p.Subject,
			&p.Body,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.DeletedAt,
		)
		if err != nil {
			return p, err
		}
	}
	p.ID = id
	return p, rows.Err()
}

//GetAll ContactMessage.
func (s ContactMessageService) GetAll() ([]model.ContactMessage, error) {
	out := []model.ContactMessage{}
	stmt := `SELECT
		oid,returnaddr, subject, body, created_at, updated_at
	 FROM contact_message
   WHERE 1=1
	 AND COALESCE(deleted_at::TEXT, '') = ''
   ORDER BY created_at DESC
   LIMIT 50
   `
	stmt = sqlx.Rebind(sqlx.BindType("postgres"), stmt)
	pstmt, err := s.DB.Prepare(stmt)
	if err != nil {
		return out, err
	}
	defer pstmt.Close()
	rows, err := pstmt.Query()
	if err != nil {
		return out, err
	}
	defer rows.Close()
	for rows.Next() {
		var p model.ContactMessage
		err = rows.Scan(
			&p.ID,
			&p.ReturnAddr,
			&p.Subject,
			&p.Body,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return out, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}
