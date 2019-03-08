package model

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/clementauger/monparcours/server/model"
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
	(?,?,?,?)`
	pstmt, err := s.DB.Prepare(stmt)
	if err != nil {
		return model.ContactMessage{}, err
	}
	defer pstmt.Close()
	createdAt := time.Now()
	res, err := pstmt.Exec(
		p.ReturnAddr, p.Subject, p.Body, createdAt)
	if err != nil {
		return p, err
	}
	p.CreatedAt = createdAt
	if p.ID, err = res.LastInsertId(); err != nil {
		return p, err
	}
	return p, nil
}

//Delete a ContactMessage
func (s ContactMessageService) Delete(p model.ContactMessage) error {
	stmt := `UPDATE contact_message SET deleted_at=? WHERE oid = ?`
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
	 AND IFNULL(deleted_at, '') = ''
	 `
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
	 AND IFNULL(deleted_at, '') = ''
   ORDER BY created_at DESC
   LIMIT 50
   `
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
