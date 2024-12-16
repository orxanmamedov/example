package auth

import "context"

type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "chi context value " + k.name
}

var (
	userCtxKey = &contextKey{"user"}
)

func FromContext(ctx context.Context) User {
	value := ctx.Value(userCtxKey)
	if value == nil {
		return User{
			ID:     -1,
			CityID: -1,
			Lang:   "en",
		}
	}
	return value.(User)
}

func InContext(ctx context.Context, u *User) context.Context {
	return context.WithValue(ctx, userCtxKey, *u)
}
