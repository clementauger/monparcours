package model

import (
	"time"
)

//Protest model
type Protest struct {
	ID          int64      `json:"id"`
	AuthorID    string     `json:"author_id"`
	Title       string     `json:"title" validate:"required,max=60" conform:"text"`
	Protest     string     `json:"protest" validate:"required,max=200" conform:"text"`
	Organizer   string     `json:"organizer" validate:"max=200" conform:"text"`
	Description string     `json:"description" validate:"max=200" conform:"text"`
	Public      bool       `json:"public"`
	Password    string     `json:"password" validate:"iffalse=Public,max=20" conform:"text"`
	GatherAt    time.Time  `json:"gather_at" validate:"required"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
	Steps       []Step     `json:"steps" validate:"max=10,min=1,required,dive,required"`
}

//Step model
type Step struct {
	ID        int64     `json:"id"`
	ProtestID int64     `json:"protest_id"`
	Place     string    `json:"place" validate:"required,max=60" conform:"text"`
	Details   string    `json:"details" validate:"max=200" conform:"text"`
	Lat       float64   `json:"lat" validate:"required,latitude"`
	Lng       float64   `json:"lng" validate:"required,longitude"`
	GatherAt  time.Time `json:"gather_at" validate:""`
}

type ProtestService interface {
	Insert(p Protest) (Protest, error)
	Delete(p Protest) error
	Get(id int64) (Protest, error)
	GetWithPassword(id int64, pwd string) (Protest, error)
	GetByAuthorID(authorID string) ([]Protest, error)
	GetByIDs(ids ...int64) ([]Protest, error)
	SearchProtests(title, protest, organizer *string, startDate, endDate *time.Time, atLat, atLng *float64, ray float64) ([]Protest, error)
}
type StepService interface {
	Insert(p Step) (Step, error)
	InsertSteps(p Protest) (Protest, error)
	InsertAll(steps []Step) ([]Step, error)
	Delete(p Step) error
	Get(id int64) (Step, error)
	GetAll(protestID int64) ([]Step, error)
	GetSteps(p Protest) (Protest, error)
	GetProtectedSteps(p Protest) (Protest, error)
	FindStepsAround(atDate time.Time, withinTime time.Duration, atLat, atLng float64) ([]Step, error)
}
