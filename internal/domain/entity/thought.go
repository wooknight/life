package entity

import (
	"time"
)

type Thought struct {
	ID int64 `json:"id"`
    Thought string  `json:"thought"`
    Created_on time.Time,
    Updated_on time.Time 
}

func NewThought (thought string) (*Thought, error){
	tht:= &Thought{
		Thought: thought,
		Created_on: time.Now(),
		Updated_on: time.Time{},
	}
	if err := tht.Validate(); err != nil {
		return nil , err 
	}
	return tht , nil 
}

func (t *Thought) Validate () error {
	if t.Thought == "" {
		return ErrInvalidThought
	}
	return nil 
}