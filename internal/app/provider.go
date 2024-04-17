package app

import (
	"golangoauth2example/internal/auth"
	"golangoauth2example/internal/middleware"
	"golangoauth2example/internal/server"
	"net/http"
)

type App struct {
	httpServer server.Server
	router     http.Handler
}

func (a *App) initRouter() http.Handler {
	if a.router == nil {
		middleware := middleware.CreateMiddlewareStack(middleware.Logging, middleware.Auth)

		mux := http.NewServeMux()

		routersList := make([]server.Router, 0)

		htmlRouters := server.NewHTMLRouter()
		cssRouters := server.NewCSSRouter()
		imagesRouters := server.NewImagesRouter()

		routersList = append(routersList, htmlRouters...)
		routersList = append(routersList, cssRouters...)
		routersList = append(routersList, imagesRouters...)

		for _, v := range routersList {
			mux.Handle(v.Pattern+"/", http.StripPrefix(v.Pattern, v.Handler))
		}
		mux.HandleFunc("/auth/google/login/", auth.OAUTHGoogleLogin)
		mux.HandleFunc("/auth/google/callback/", auth.OAUTHGoogleCallback)
		mux.Handle("/account/", middleware(server.NewAccountRouter().Handler))

		a.router = mux
	}
	return a.router
}

func (a *App) InitServer() server.Server {
	if a.httpServer == nil {
		a.httpServer = server.NewHTTPServer(a.router)
	}

	return a.httpServer
}

func (a *App) Start() error {
	return a.httpServer.Run()
}

func NewApp() *App {
	app := &App{}

	app.initRouter()
	app.InitServer()

	return app

}
