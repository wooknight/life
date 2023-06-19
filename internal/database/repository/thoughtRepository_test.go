package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/wooknight/life/internal/domain/entity"
)

func TestGetThought(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewThoughtRepository(db)

	thought := &entity.Thought{
		ID:        1,
		Thought:   "Let's Go Further!",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Time{}.UTC(),
	}

	t.Run("OK", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "thought", "created_at", "updated_at"}).
			AddRow(thought.ID, thought.Thought, thought.CreatedAt, thought.UpdatedAt)

		mock.ExpectPrepare("select \\* from thoughts WHERE id = $1").
			ExpectQuery().
			WithArgs(thought.ID).
			WillReturnRows(rows)

		gotThought, err := repo.Get(context.Background(), thought.ID)
		assert.NoError(t, err)
		assert.Equal(t, thought, gotThought)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Prepare Failed", func(t *testing.T) {
		mock.ExpectPrepare("select \\* from thoughts WHERE id = ").
			WillReturnError(sql.ErrConnDone)

		gotThought, err := repo.Get(context.Background(), thought.ID)
		assert.Error(t, err)
		assert.Empty(t, gotThought)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Query Failed", func(t *testing.T) {
		mock.ExpectPrepare("select \\* from thoughts WHERE id = ").
			ExpectQuery().
			WithArgs(thought.ID).
			WillReturnError(sql.ErrConnDone)

		gotThought, err := repo.Get(context.Background(), thought.ID)
		assert.Error(t, err)
		assert.Empty(t, gotThought)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Not Found", func(t *testing.T) {
		mock.ExpectPrepare("select \\* from thoughts WHERE id = ").
			ExpectQuery().
			WithArgs(thought.ID).
			WillReturnError(sql.ErrNoRows)

		gotThought, err := repo.Get(context.Background(), thought.ID)
		assert.Error(t, err)
		assert.Empty(t, gotThought)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestListThought(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewThoughtRepository(db)

	thoughts := []*entity.Thought{
		{
			ID:        1,
			Thought:   "Let's Go Further!",
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Time{},
		}, {
			ID:        2,
			Thought:   "Let's Go!",
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Time{},
		},
	}

	t.Run("OK", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "thought", "created_at", "updated_at"})
		for _, thought := range thoughts {
			rows = rows.AddRow(thought.ID, thought.Thought, thought.CreatedAt, thought.UpdatedAt)
		}

		mock.ExpectPrepare("select \\* from thoughts").
			ExpectQuery().
			WillReturnRows(rows)

		gotThoughts, err := repo.List(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, len(thoughts), len(gotThoughts))

		for i := range gotThoughts {
			assert.Equal(t, thoughts[i], gotThoughts[i])
		}

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Prepare Failed", func(t *testing.T) {
		mock.ExpectPrepare("select \\* from thoughts").
			WillReturnError(sql.ErrConnDone)

		gotThought, err := repo.List(context.Background())
		assert.Error(t, err)
		assert.Empty(t, gotThought)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Query Failed", func(t *testing.T) {
		mock.ExpectPrepare("select \\* from thoughts").
			ExpectQuery().
			WillReturnError(sql.ErrConnDone)

		gotThought, err := repo.List(context.Background())
		assert.Error(t, err)
		assert.Empty(t, gotThought)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Empty thought", func(t *testing.T) {
		mock.ExpectPrepare("select \\* from thoughts").
			ExpectQuery().
			WillReturnRows(&sqlmock.Rows{})

		gotThoughts, err := repo.List(context.Background())
		assert.Error(t, err)
		assert.Nil(t, gotThoughts)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSearchThoughts(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewThoughtRepository(db)

	thoughts := []*entity.Thought{
		{
			ID:      1,
			Thought: "Let's Go Further!",

			UpdatedAt: time.Time{},
			CreatedAt: time.Now().UTC(),
		}, {
			ID:      2,
			Thought: "Let's Go!",

			UpdatedAt: time.Time{}.UTC(),
			CreatedAt: time.Now().UTC(),
		},
	}

	t.Run("OK", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "thought", "updated_at", "created_at"})
		for _, thought := range thoughts {
			rows = rows.AddRow(thought.ID, thought.Thought, thought.UpdatedAt, thought.CreatedAt)
		}

		mock.ExpectPrepare("select \\* from thoughts WHERE LOWER\\(thought\\) LIKE LOWER\\(\\$1\\)").
			ExpectQuery().
			WithArgs("%Let's Go%").
			WillReturnRows(rows)

		gotThoughts, err := repo.Search(context.Background(), "Let's Go")

		assert.NoError(t, err)
		assert.Equal(t, len(thoughts), len(gotThoughts))

		for i := range gotThoughts {
			assert.Equal(t, thoughts[i], gotThoughts[i])
		}

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Prepare Failed", func(t *testing.T) {
		mock.ExpectPrepare("select \\* from thoughts WHERE LOWER\\(thought\\) LIKE LOWER\\(\\$1\\)").
			WillReturnError(sql.ErrConnDone)

		gotThoughts, err := repo.Search(context.Background(), "Let's Go")

		assert.Error(t, err)
		assert.Empty(t, gotThoughts)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Query Failed", func(t *testing.T) {
		mock.ExpectPrepare("select \\* from thoughts WHERE LOWER\\(thought\\) LIKE LOWER\\(\\$1\\)").
			ExpectQuery().
			WillReturnError(sql.ErrConnDone)

		gotThoughts, err := repo.Search(context.Background(), "Let's Go")

		assert.Error(t, err)
		assert.Empty(t, gotThoughts)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Empty Thought", func(t *testing.T) {
		mock.ExpectPrepare("select \\* from thoughts WHERE LOWER\\(thought\\) LIKE LOWER\\(\\$1\\)").
			ExpectQuery().
			WillReturnRows(&sqlmock.Rows{})

		gotThoughts, err := repo.Search(context.Background(), "Let's Go")

		assert.Error(t, err)
		assert.Empty(t, gotThoughts)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Not Found", func(t *testing.T) {
		mock.ExpectPrepare("select \\* from thoughts WHERE LOWER\\(thought\\) LIKE LOWER\\(\\$1\\)").
			ExpectQuery().
			WithArgs("%Let's Go%").
			WillReturnRows(&sqlmock.Rows{})

		gotThoughts, err := repo.Search(context.Background(), "Let's Go")

		assert.Error(t, err)
		assert.Empty(t, gotThoughts)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCreateThought(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewThoughtRepository(db)

	thought := &entity.Thought{
		ID:      1,
		Thought: "Let's Go Further!",

		UpdatedAt: time.Time{},
		CreatedAt: time.Now().UTC(),
	}

	t.Run("OK", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id"}).
			AddRow(1)

		mock.ExpectPrepare("INSERT INTO thoughts").
			ExpectQuery().
			WithArgs(thought.Thought).
			WillReturnRows(rows)

		id, err := repo.Create(context.Background(), thought)
		assert.NoError(t, err)
		assert.NotEmpty(t, id)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Prepare Failed", func(t *testing.T) {
		mock.ExpectPrepare("INSERT INTO thoughts").
			WillReturnError(sql.ErrConnDone)

		id, err := repo.Create(context.Background(), thought)
		assert.Error(t, err)
		assert.Empty(t, id)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Query Failed", func(t *testing.T) {
		mock.ExpectPrepare("INSERT INTO thoughts").
			ExpectQuery().
			WillReturnError(sql.ErrConnDone)

		id, err := repo.Create(context.Background(), thought)
		assert.Error(t, err)
		assert.Empty(t, id)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUpdateThought(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewThoughtRepository(db)

	thought := &entity.Thought{
		ID:      1,
		Thought: "Let's Go Further!",

		UpdatedAt: time.Time{},
		CreatedAt: time.Now().UTC(),
	}

	t.Run("OK", func(t *testing.T) {
		mock.ExpectPrepare("UPDATE thoughts").
			ExpectExec().
			WithArgs(thought.Thought, thought.ID).
			WillReturnResult(sqlmock.NewResult(int64(thought.ID), 1))

		err := repo.Update(context.Background(), thought)
		assert.NoError(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Prepare Failed", func(t *testing.T) {
		mock.ExpectPrepare("UPDATE thoughts").
			WillReturnError(sql.ErrConnDone)

		err := repo.Update(context.Background(), thought)
		assert.Error(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Prepare Failed", func(t *testing.T) {
		mock.ExpectPrepare("UPDATE thoughts").
			ExpectExec().
			WillReturnError(sql.ErrConnDone)

		err := repo.Update(context.Background(), thought)
		assert.Error(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Not Found", func(t *testing.T) {
		mock.ExpectPrepare("UPDATE thoughts").
			ExpectExec().
			WithArgs(thought.Thought, thought.ID).
			WillReturnResult(sqlmock.NewResult(int64(thought.ID), 0))

		err := repo.Update(context.Background(), thought)
		assert.Error(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDeleteThought(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewThoughtRepository(db)

	thought := &entity.Thought{
		ID:      1,
		Thought: "Let's Go Further!",

		UpdatedAt: time.Time{},
		CreatedAt: time.Now().UTC(),
	}

	t.Run("OK", func(t *testing.T) {
		mock.ExpectPrepare("DELETE from thoughts WHERE id =").
			ExpectExec().
			WithArgs(thought.ID).
			WillReturnResult(sqlmock.NewResult(int64(thought.ID), 1))

		err := repo.Delete(context.Background(), thought.ID)
		assert.NoError(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Prepare Failed", func(t *testing.T) {
		mock.ExpectPrepare("DELETE from thoughts WHERE id =").
			WillReturnError(sql.ErrConnDone)

		err := repo.Delete(context.Background(), thought.ID)
		assert.Error(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Exec Failed", func(t *testing.T) {
		mock.ExpectPrepare("DELETE from thoughts WHERE id =").
			ExpectExec().
			WillReturnError(sql.ErrConnDone)

		err := repo.Delete(context.Background(), thought.ID)
		assert.Error(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Not Found", func(t *testing.T) {
		mock.ExpectPrepare("DELETE from thoughts WHERE id =").
			ExpectExec().
			WithArgs(thought.ID).
			WillReturnResult(sqlmock.NewResult(int64(thought.ID), 0))

		err := repo.Delete(context.Background(), thought.ID)
		assert.Error(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
