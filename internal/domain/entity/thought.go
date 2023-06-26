package entity

import (
	"time"

	"github.com/lib/pq"
)

// https://stackoverflow.com/questions/55899890/convert-a-postgres-row-into-golang-struct-with-array-field
type Thought struct {
	ID                 int64          `json:"id"`
	Thought            string         `json:"thought"`
	ThoughtDescription string         `json:"description" pg:"description"`
	ThoughtTags        pq.StringArray `db:"tags" json:"tags"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
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
