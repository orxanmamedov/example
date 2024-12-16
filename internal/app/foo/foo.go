package foo

import (
	"context"
	"example/internal/repository"
	"example/pkg/storage/mysql"
	"github.com/google/uuid"

	"github.com/pkg/errors"
)

type Foo struct {
	r repository.Foo
}

var ErrNotFound = errors.New("not found")

func New(r repository.Foo) *Foo {
	return &Foo{
		r: r,
	}
}

func (f *Foo) F(ctx context.Context, userID int64) (*F, error) {
	resp, err := f.r.F(ctx, userID)
	if err == mysql.ErrNoRows {
		return nil, errors.Wrap(ErrNotFound, "foo")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "problem foo for user")
	}
	return &F{
		ID:     uuid.NewString(),
		UserID: resp.UserID,
		Count:  resp.Count,
	}, nil
}
