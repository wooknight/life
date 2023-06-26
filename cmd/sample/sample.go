package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"
)

// https://stackoverflow.com/questions/55899890/convert-a-postgres-row-into-golang-struct-with-array-field
type Thought struct {
	ID                 int64          `json:"id"`
	Thought            string         `json:"thought" db:"thought"`
	ThoughtDescription string         `json:"description" db:"descr"`
	ThoughtTags        pq.StringArray `db:"tags" json:"tags"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func NewThought(thought, desc string, tags []string) (*Thought, error) {
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
	if t.Thought == "" {
		return errors.New("invalid thought")
	}
	return nil
}

func main() {
	// open database
	db, err := sql.Open("postgres", "host=0.0.0.0 port=5400 user=postgres password=password dbname=life sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var entry Thought
	entry.Thought = "luvda"
	entry.ThoughtDescription = "lassoon"
	entry.ThoughtTags = []string{"luvda", "lasoon", "kanda", "batata"}
	insertDynStmt := `insert into "thoughts"("thought", "descr","tags") values($1, $2,$3)`
	_, err = db.Exec(insertDynStmt, entry.Thought, entry.ThoughtDescription, pq.Array(entry.ThoughtTags))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	var thgtarr []Thought
	query := `SELECT id,thought, descr,tags,created_at, updated_at FROM thoughts`
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var cur Thought
		var updatedAt sql.NullTime
		var tags []string
		err = rows.Scan(&cur.ID, &cur.Thought, &cur.ThoughtDescription, pq.Array(&tags), &cur.CreatedAt, &updatedAt)
		if err != nil {
			log.Fatal(err)
		}
		if updatedAt.Valid {
			cur.UpdatedAt = updatedAt.Time
		}
		if len(tags) > 0 {
			cur.ThoughtTags = tags
		}
		thgtarr = append(thgtarr, cur)
	}

	if len(thgtarr) == 0 {
		fmt.Println("Could not find any data")
	}

	fmt.Println(thgtarr)

	// var foo Thought
	// err = db.QueryRowx(`SELECT thought, description, tags FROM thoughts where id = ? LIMIT 1`).StructScan(&foo)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(foo)
}
