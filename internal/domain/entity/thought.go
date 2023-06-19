package entity

import (
	"time"
)

type Thought struct {
	ID        int64  `json:"id"`
	Thought   string `json:"thought"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewThought(thought string) (*Thought, error) {
	tht := &Thought{
		Thought:   thought,
		CreatedAt: time.Now(),
		UpdatedAt: time.Time{},
	}
	if err := tht.Validate(); err != nil {
		return nil, err
	}
	return tht, nil
}

func (t *Thought) Validate() error {
	if t.Thought == "" {
		return ErrInvalidThought
	}
	return nil
}
