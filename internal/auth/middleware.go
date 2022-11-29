package auth

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Y-bro/go-graphql/internal/users"
	"github.com/Y-bro/go-graphql/pkg/jwt"
)

type contextKey struct {
	Name string
}

var userCtxKey = &contextKey{"user"}

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			header := r.Header.Get("Authorization")

			if header == "" {
				fmt.Println("Here")
				next.ServeHTTP(w, r)
				return
			}

			tokenStr := header

			fmt.Println("Here")
			username, err := jwt.ParseToken(tokenStr)
			fmt.Println("Here")

			if err != nil {
				http.Error(w, "Token invalid", http.StatusUnauthorized)
				return
			}

			user := users.User{Username: username}

			id, err := users.GetUserIdByUsername(username)

			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			user.ID = strconv.Itoa(id)

			ctx := context.WithValue(r.Context(), userCtxKey, &user)

			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)

		})
	}
}

func ForContext(ctx context.Context) *users.User {
	raw, _ := ctx.Value(userCtxKey).(*users.User)

	return raw
}
