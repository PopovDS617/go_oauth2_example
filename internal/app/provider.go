package app

import (
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

		mw := middleware.CreateMiddlewareStack(middleware.Logging)

		mux := http.NewServeMux()

		routersList := make([]server.Router, 0)

		// file routers
		htmlRouter := server.NewHTMLRouter()
		cssRouter := server.NewCSSRouter()
		imagesRouter := server.NewImagesRouter()

		// oauth routers
		googleOAUTHRouter := server.NewGoogleOAUTHRouter()

		// personal account routers
		accountRouter := server.NewAccountRouter()
		accountRouter.Handler = middleware.Auth(accountRouter.Handler)

		routersList = append(routersList, htmlRouter, cssRouter, imagesRouter, googleOAUTHRouter, accountRouter)

		for _, v := range routersList {
			mux.Handle(v.Pattern+"/", mw(http.StripPrefix(v.Pattern, v.Handler)))

		}

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
