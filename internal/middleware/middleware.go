package middleware

import (
	"fmt"
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

		fmt.Println(r.URL)

		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		cookieValue := authCookie.Value

		if cookieValue != "true" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

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
