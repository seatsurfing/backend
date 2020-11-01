package main

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

var _appInstance *App
var _appOnce sync.Once

func GetApp() *App {
	_appOnce.Do(func() {
		_appInstance = &App{}
	})
	return _appInstance
}

type App struct {
	Router        *mux.Router
	CleanupTicker *time.Ticker
}

func (a *App) InitializeDatabases() {
	RunDBSchemaUpdates()
	InitDefaultOrgSettings()
}

func (a *App) InitializeRouter() {
	a.Router = mux.NewRouter()
	routers := make(map[string]Route)
	routers["/location/{locationId}/space/"] = &SpaceRouter{}
	routers["/location/"] = &LocationRouter{}
	routers["/booking/"] = &BookingRouter{}
	routers["/organization/"] = &OrganizationRouter{}
	routers["/auth-provider/"] = &AuthProviderRouter{}
	routers["/auth/"] = &AuthRouter{}
	routers["/user/"] = &UserRouter{}
	routers["/stats/"] = &StatsRouter{}
	routers["/search/"] = &SearchRouter{}
	routers["/setting/"] = &SettingsRouter{}
	routers["/signup/"] = &SignupRouter{}
	routers["/fastspring/"] = &FastSpringRouter{}
	for route, router := range routers {
		subRouter := a.Router.PathPrefix(route).Subrouter()
		router.setupRoutes(subRouter)
	}
	a.setupStaticRoutes(a.Router)
	a.Router.PathPrefix("/").Methods("OPTIONS").HandlerFunc(CorsHandler)
	a.Router.Use(CorsMiddleware)
	a.Router.Use(VerifyAuthMiddleware)
}

func (a *App) InitializeTimers() {
	a.CleanupTicker = time.NewTicker(time.Minute * 5)
	go func() {
		for {
			select {
			case <-a.CleanupTicker.C:
				log.Println("Cleaning up expired database entries...")
				if err := GetAuthStateRepository().DeleteExpired(); err != nil {
					log.Println(err)
				}
				if err := GetSignupRepository().DeleteExpired(); err != nil {
					log.Println(err)
				}
			}
		}
	}()
}

func (a *App) setupStaticRoutes(router *mux.Router) {
	paths := []string{
		"/login",
		"/dashboard",
		"/locations",
		"/users",
		"/settings",
		"/bookings",
		"/search",
	}
	fs := http.FileServer(http.Dir(GetConfig().StaticFilesPath))
	stripPrefix := func(prefix string) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r2 := new(http.Request)
			*r2 = *r
			r2.URL = new(url.URL)
			*r2.URL = *r.URL
			r2.URL.Path = "/"
			fs.ServeHTTP(w, r2)
		})
	}
	for _, path := range paths {
		path = "/admin" + path
		router.PathPrefix(path).Handler(stripPrefix(path))
	}
	router.Path("/admin/").Handler(stripPrefix("/admin/"))
	router.PathPrefix("/admin/").Handler(http.StripPrefix("/admin/", fs))
}

func (a *App) Run(publicListenAddr string) {
	log.Println("Initializing REST services...")
	httpServer := &http.Server{
		Addr:         publicListenAddr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      a.Router,
	}
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatal(err)
			os.Exit(-1)
		}
	}()
	log.Println("HTTP Server listening on", publicListenAddr)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println("Shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	httpServer.Shutdown(ctx)
}
