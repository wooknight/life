package repository

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/wooknight/life/internal/domain/entity"
)

type thoughtRepository struct {
	db *sql.DB
}

func NewThoughtRepository(db *sql.DB) *thoughtRepository {
	return &thoughtRepository{
		db: db,
	}
}

// GET
func (r *thoughtRepository) Get(ctx context.Context, id int64) (*entity.Thought, error) {
	stmt, err := r.db.PrepareContext(ctx, "select * from thoughts where id = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var updatedAt sql.NullTime
	t := &entity.Thought{}
	row := stmt.QueryRowContext(ctx, id)
	err = row.Scan(&t.ID, &t.Thought, &t.ThoughtDescription, &t.ThoughtTags, &t.CreatedAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrThoughtNotFound
		}
		return nil, fmt.Errorf("%s: %w", ErrScanData, err)
	}
	if updatedAt.Valid {
		t.UpdatedAt = updatedAt.Time
	}
	return t, nil
}

func (r *thoughtRepository) List(ctx context.Context) ([]*entity.Thought, error) {
	stmt, err := r.db.PrepareContext(ctx, "select * from thoughts")
	if err != nil {
		return nil, fmt.Errorf("%s:%w", ErrPrepareStatement, err)
	}
	defer stmt.Close()
	var results []*entity.Thought
	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", ErrExecuteQuery, err)
	}
	for rows.Next() {
		var t entity.Thought
		var updatedAt sql.NullTime
		err := rows.Scan(&t.ID, &t.Thought, &t.CreatedAt, &updatedAt)
		if err != nil {
			return nil, fmt.Errorf("%s:%w", ErrScanData, err)
		}
		if updatedAt.Valid {
			t.UpdatedAt = updatedAt.Time
		}
		results = append(results, &t)

	}
	if len(results) == 0 {
		return nil, ErrThoughtNotFound
	}

	return results, nil
}

func (r *thoughtRepository) Create(ctx context.Context, idea *entity.Thought) (int64, error) {
	stmt, err := r.db.PrepareContext(ctx, "insert into thoughts (thought, descr, tags) values ($1, $2, $3) returning id")
	if err != nil {
		return 0, fmt.Errorf("%s:%w", ErrPrepareStatement, err)
	}
	defer stmt.Close()
	err = stmt.QueryRowContext(ctx, idea.Thought, idea.ThoughtDescription, idea.ThoughtTags).Scan(&idea.ID)
	if err != nil {
		return 0, fmt.Errorf("%s:%w", ErrExecuteQuery, err)
	}
	return idea.ID, nil
}

func (r *thoughtRepository) Search(ctx context.Context, searchStr string) ([]*entity.Thought, error) {
	stmt, err := r.db.PrepareContext(ctx, "select * from thoughts where lower(thought) like lower($1)")
	if err != nil {
		return nil, fmt.Errorf("%s:%w", ErrPrepareStatement, err)
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, "%"+searchStr+"%")
	if err != nil {
		return nil, fmt.Errorf("%s:%w", ErrExecuteQuery, err)
	}
	var results []*entity.Thought
	for rows.Next() {
		var t entity.Thought
		var updatedAt sql.NullTime
		err = rows.Scan(&t.ID, &t.Thought, &t.ThoughtDescription, &t.ThoughtTags, &t.CreatedAt, &updatedAt)
		if err != nil {
			return nil, fmt.Errorf("%s:%w", ErrScanData, err)
		}
		if updatedAt.Valid {
			t.UpdatedAt = updatedAt.Time
		}
		results = append(results, &t)
	}
	if len(results) == 0 {
		return nil, ErrThoughtNotFound
	}

	return results, nil
}

func (r *thoughtRepository) Update(ctx context.Context, t *entity.Thought) error {
	stmt, err := r.db.PrepareContext(ctx, "update thoughts set thought = $1,descr = $2 , tags = $3, updated_at = NOW() where id = $4")
	if err != nil {
		return fmt.Errorf("%s:%w", ErrPrepareStatement, err)
	}
	defer stmt.Close()
	result, err := stmt.ExecContext(ctx, t.Thought, t.ID)
	if err != nil {
		return fmt.Errorf("%s:%w", ErrExecuteQuery, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s:%w", ErrRetrieveRows, err)
	}
	if rowsAffected == 0 {
		return ErrThoughtNotFound
	}
	return nil
}

func (r *thoughtRepository) Delete(ctx context.Context, id int64) error {
	stmt, err := r.db.PrepareContext(ctx, "delete from thought where id = $1")
	if err != nil {
		return fmt.Errorf("%s:%w", ErrPrepareStatement, err)
	}
	defer stmt.Close()
	result, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return fmt.Errorf("%s:%w", ErrExecuteQuery, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s:%w", ErrRetrieveRows, err)
	}
	if rowsAffected == 0 {
		return ErrThoughtNotFound
	}
	return nil
}
