package auth

import (
	"encoding/json"
	"net/http"

	"example/pkg/logger"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/pkg/errors"
)

const jwtUserKey = "user"

var tokenAuth *jwtauth.JWTAuth

func Init(secret string) {
	if tokenAuth == nil {
		tokenAuth = jwtauth.New(jwa.HS256.String(), []byte(secret), nil)
	}
}

func Verifier() func(http.Handler) http.Handler {
	return jwtauth.Verifier(tokenAuth)
}

func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, err := jwtauth.FromContext(r.Context())

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		
		u := &User{}
		jwtUser, ok := claims[jwtUserKey]
		if !ok {
			logger.Errorf(r.Context(), "user is not in jwt: %v", claims)
			errorJSON(w, r, http.StatusInternalServerError, errors.New("user is not in jwt"))
			return
		}

		jsonUser, _ := json.Marshal(jwtUser)
		err = json.Unmarshal(jsonUser, u)
		if err != nil {
			logger.Errorf(r.Context(), "cannot unmarshal user from jwt: %s", err)
			errorJSON(w, r, http.StatusInternalServerError, errors.New("cannot unmarshal user from jwt"))
			return
		}

		lang := r.Header.Get("Accept-Language")
		u.Lang = lang
		if u.Lang == "" {
			u.Lang = "en"
		}
		uCtx := InContext(r.Context(), u)
		next.ServeHTTP(w, r.WithContext(uCtx))
	})
}
