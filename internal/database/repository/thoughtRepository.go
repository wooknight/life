package repository

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/wookinight/life/internal/domain/entity"
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
func (r *thoughtRepository) Get(ctx context.Context, id int) (*thoughtRepository, error) {
	stmt, err := r.db.PrepareContext(ctx, "SELECT * FROM thought WHERE id = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	t := &entity.Thought()
	stmt.QueryRowContext(ctx, id)

}
