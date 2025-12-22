package model

import (
	"time"
)

type Timestamps struct {
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func NewTimestamps() Timestamps {
	return Timestamps{
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

func (t *Timestamps) Update() *Timestamps {
	t.UpdatedAt = time.Now()
	return t
}

func (t *Timestamps) Delete() *Timestamps {
	now := time.Now()
	t.DeletedAt = &now
	return t
}
