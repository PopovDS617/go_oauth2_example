package server

import (
	"golangoauth2example/internal/auth"
	"net/http"
)

type Router struct {
	Handler http.Handler
	Pattern string
}

func NewRouterFromParams(routers []Router) []Router {
	mux := http.NewServeMux()

	result := make([]Router, 0, len(routers))

	for _, v := range routers {
		mux.Handle(v.Pattern, v.Handler)

		result = append(result, Router{Pattern: v.Pattern, Handler: v.Handler})
	}

	return result
}

func NewHTMLRouter() Router {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("templates/")))

	return Router{
		Handler: mux,
		Pattern: "",
	}

}

func NewCSSRouter() Router {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("assets/css/")))

	return Router{
		Handler: mux,
		Pattern: "/css",
	}

}

func NewImagesRouter() Router {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("assets/images/")))

	return Router{
		Handler: mux,
		Pattern: "/images",
	}

}

func NewGoogleOAUTHRouter() Router {
	mux := http.NewServeMux()

	patternLogin := "/login"
	patternCallback := "/callback"

	mux.HandleFunc(patternLogin+"/", auth.OAUTHGoogleLogin)
	mux.HandleFunc(patternCallback+"/", auth.OAUTHGoogleCallback)

	return Router{
		Handler: mux,
		Pattern: "/auth/google",
	}

}

func NewAccountRouter() Router {
	mux := http.NewServeMux()

	pattern := "/account"

	mux.HandleFunc("/", auth.HandleGetPersonalAccount)

	return Router{
		Pattern: pattern,
		Handler: mux,
	}
}
