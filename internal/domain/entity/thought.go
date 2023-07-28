package entity

import (
	"time"

	"github.com/lib/pq"
)

// https://stackoverflow.com/questions/55899890/convert-a-postgres-row-into-golang-struct-with-array-field
type Thought struct {
	ID                 int64          `json:"id"`
	Thought            string         `json:"thought" pg:"thought"`
	ThoughtDescription string         `json:"description" pg:"descr"`
	ThoughtTags        pq.StringArray `pq:"tags" json:"tags"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func NewThought(thought string, desc string, tags []string) (*Thought, error) {
	tht := &Thought{
		Thought:            thought,
		ThoughtDescription: desc,
		ThoughtTags:        tags,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Time{},
	}
	if err := tht.Validate(); err != nil {
		return nil, err
	}
	return tht, nil
}

func (t *Thought) Validate() error {
	if t.Thought == "" || t.ThoughtDescription == "" || len(t.ThoughtTags) == 0 {
		return ErrInvalidThought
	}
	return nil
}
