package model

import (
	"time"
)

//ContactMessage model
type ContactMessage struct {
	ID         int64      `json:"id" sql:"oid"`
	ReturnAddr string     `json:"returnaddr" validate:"required,max=60" conform:"text"`
	Subject    string     `json:"subject" validate:"required,max=60" conform:"text"`
	Body       string     `json:"body" validate:"required,max=200" conform:"text"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
}

type ContactMessageService interface {
	Insert(p ContactMessage) (ContactMessage, error)
	Delete(p ContactMessage) error
	Get(id int64) (ContactMessage, error)
	GetAll() ([]ContactMessage, error)
}
