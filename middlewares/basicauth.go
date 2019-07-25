package middlewares

import (
	"crypto/subtle"
	"github.com/go-chi/render"
	"goGFG/controllers"
	"net/http"
)

type auth struct {
	USERNAME string
	PASSWORD string
}

var authCredentials = auth{
	"shaumux",
	"secretPassword",
}

func BasicAuthMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()

		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(authCredentials.USERNAME)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(authCredentials.PASSWORD)) != 1 {
			render.Render(w, r, controllers.UnAuthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
