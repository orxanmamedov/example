package repository

import (
	"context"
	"example/pkg/storage/mysql"
)

type FData struct {
	UserID int64 `db:"user_id"`
	Count  int64 `db:"count"`
}

type Foo interface {
	F(ctx context.Context, userID int64) (*FData, error)
}

func (r *Repo) F(ctx context.Context, userID int64) (*FData, error) {
	var f FData

	err := r.db.GetContext(ctx, &f, "SELECT ? AS user_id, 2 AS count", userID)
	if err == mysql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return &f, nil
}
