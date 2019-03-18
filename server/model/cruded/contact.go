package cruded

import (
	"github.com/clementauger/monparcours/server/model"
	"github.com/clementauger/crud"
)

//ContactMessageService sqlite implementation.
type ContactMessageService struct {
	Crud *crud.Crudable
}

//Insert a ContactMessage
func (s *ContactMessageService) Insert(p model.ContactMessage) (model.ContactMessage, error) {
	err := s.Crud.Insert(&p)
	return p, err
}

//Delete a ContactMessage
func (s *ContactMessageService) Delete(p model.ContactMessage) error {
	return s.Crud.Delete(&p)
}

//Get a ContactMessage by its ID.
func (s *ContactMessageService) Get(id int64) (model.ContactMessage, error) {
	var p model.ContactMessage
	err := s.Crud.FetchByID(&p, id)
	return p, err
}

//GetAll ContactMessage.
func (s *ContactMessageService) GetAll() ([]model.ContactMessage, error) {
	out := []model.ContactMessage{}
	stmt := `ORDER BY created_at DESC
   LIMIT 50
   `
	err := s.Crud.Fetch(&out, stmt)
	return out, err
}
