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

func NewHTMLRouter() []Router {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("templates/")))

	routers := []Router{{
		Handler: mux,
		Pattern: "",
	}}

	return routers
}

func NewCSSRouter() []Router {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("assets/css/")))

	routers := []Router{{
		Handler: mux,
		Pattern: "/css",
	},
	}

	return routers
}

func NewImagesRouter() []Router {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("assets/images/")))

	routers := []Router{{
		Handler: mux,
		Pattern: "/images",
	}}

	return routers
}

func NewGoogleOAUTHRouter() []Router {
	mux := http.NewServeMux()

	patternLogin := "/auth/google/login"
	patternCallback := "/auth/google/callback"

	mux.HandleFunc(patternLogin+"/", auth.OAUTHGoogleLogin)
	mux.HandleFunc(patternCallback+"/", auth.OAUTHGoogleCallback)

	routers := []Router{{
		Handler: mux,
		Pattern: patternLogin,
	},
		{Handler: mux, Pattern: patternCallback},
	}

	return routers

}

func NewAccountRouter() Router {
	mux := http.NewServeMux()

	pattern := "/account"

	mux.HandleFunc(pattern+"/", auth.HandleGetPersonalAccount)

	return Router{
		Pattern: pattern,
		Handler: mux,
	}
}
