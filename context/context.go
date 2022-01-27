package context

import (
	"context"
	"goweb_v1/models"
)

const (
	userKey privatekey = "user"
)

// to prevent controller override the string type key is to use a non-export
// type
type privatekey string

func WithUser(ctx context.Context, user *models.User) context.Context {
	// returns a context with key "user" and value the current user
	return context.WithValue(ctx, userKey, user)
}

func User(ctx context.Context) *models.User {
	// ctx values only return value
	if temp := ctx.Value(userKey); temp != nil {
		// type conversion, convert temp to user type
		if user, ok := temp.(*models.User); ok {
			return user
		}
	}
	return nil
}
