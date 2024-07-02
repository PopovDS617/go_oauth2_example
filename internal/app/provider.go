package app

import (
	httpServer "golangoauth2example/internal/http"
	"golangoauth2example/internal/middleware"
	"net/http"
)

type App struct {
	httpServer httpServer.Server
	router     http.Handler
}

func (a *App) initRouter() http.Handler {
	if a.router == nil {

		basicMiddleware := middleware.CreateMiddlewareStack(middleware.Logging, middleware.CORS)

		mux := http.NewServeMux()

		routersList := make([]httpServer.Router, 0)

		// file routers
		htmlRouter := httpServer.NewHTMLRouter()
		htmlRouter.Handler = middleware.Caching(htmlRouter.Handler)

		cssRouter := httpServer.NewCSSRouter()
		cssRouter.Handler = middleware.Caching(cssRouter.Handler)

		imagesRouter := httpServer.NewImagesRouter()
		imagesRouter.Handler = middleware.Caching(imagesRouter.Handler)

		// oauth routers
		googleOAUTHRouter := httpServer.NewGoogleOAUTHRouter()

		// personal account routers
		accountRouter := httpServer.NewAccountRouter()
		accountRouter.Handler = middleware.Auth(accountRouter.Handler)

		routersList = append(routersList, htmlRouter, cssRouter, imagesRouter, googleOAUTHRouter, accountRouter)

		for _, v := range routersList {
			mux.Handle(v.Pattern+"/", basicMiddleware(http.StripPrefix(v.Pattern, v.Handler)))

		}

		a.router = mux
	}
	return a.router
}

func (a *App) InitServer() httpServer.Server {
	if a.httpServer == nil {
		a.httpServer = httpServer.NewHTTPServer(a.router)
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
