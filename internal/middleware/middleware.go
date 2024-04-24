package middleware

import (
	"context"
	"golangoauth2example/internal/auth"
	"golangoauth2example/internal/utils"
	"log"
	"net/http"
)

type Middleware func(http.Handler) http.Handler

func CreateMiddlewareStack(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			x := xs[i]
			next = x(next)
		}
		return next
	}

}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authCookie, err := r.Cookie("gotestsession")

		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		cookieValue := authCookie.Value

		user, err := auth.ParseJWTToken(cookieValue)

		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		ctx := context.Background()

		ctx = context.WithValue(ctx, "user", *user)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func Logging(wrappedHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		lrw := utils.NewLoggingResponseWriter(w)
		wrappedHandler.ServeHTTP(lrw, r)

		statusCode := lrw.StatusCode

		log.Printf("%s %s %d\n", r.Method, r.URL.Path, statusCode)
	})
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "htt://localhost:8000")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		next.ServeHTTP(w, r)
	})
}

func Caching(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Cache-Control", "public, max-age=86400")

		next.ServeHTTP(w, r)
	})
}
