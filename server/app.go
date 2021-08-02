package main

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"time"

	gonnect "github.com/craftamap/atlas-gonnect"
	"github.com/craftamap/atlas-gonnect/middleware"
	"github.com/craftamap/atlas-gonnect/routes"
	"github.com/craftamap/atlas-gonnect/store"
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
	config := GetConfig()
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
	if config.OrgSignupEnabled {
		routers["/signup/"] = &SignupRouter{}
	}
	for route, router := range routers {
		subRouter := a.Router.PathPrefix(route).Subrouter()
		router.setupRoutes(subRouter)
	}
	a.setupStaticAdminRoutes(a.Router)
	a.setupStaticUserRoutes(a.Router)
	a.Router.PathPrefix("/").Methods("OPTIONS").HandlerFunc(CorsHandler)
	a.Router.Use(CorsMiddleware)
	a.Router.Use(VerifyAuthMiddleware)
}

func (a *App) InitializeDefaultOrg() {
	ids, err := GetOrganizationRepository().GetAllIDs()
	if err == nil && len(ids) == 0 {
		config := GetConfig()
		org := &Organization{
			Name:       config.InitOrgName,
			Language:   strings.ToLower(config.InitOrgLanguage),
			Country:    strings.ToUpper(config.InitOrgCountry),
			SignupDate: time.Now().UTC(),
		}
		GetOrganizationRepository().Create(org)
		GetSettingsRepository().Set(org.ID, SettingSubscriptionMaxUsers.Name, "10000")
		GetOrganizationRepository().AddDomain(org, config.InitOrgDomain, true)
		user := &User{
			OrganizationID: org.ID,
			Email:          config.InitOrgUser + "@" + config.InitOrgDomain,
			HashedPassword: NullString(GetUserRepository().GetHashedPassword(config.InitOrgPass)),
			OrgAdmin:       true,
			SuperAdmin:     false,
		}
		GetUserRepository().Create(user)
		GetOrganizationRepository().createSampleData(org)
	}
}

func (a *App) InitializeAtlassianConnect() {
	if !GetConfig().Development {
		os.Setenv("GONNECT_PROFILE", "prod")
	}
	config := a.getAtlassianConfig()
	descriptor := a.getAtlassianDescriptor()
	addon, err := gonnect.NewAddon(config, descriptor)
	if err != nil {
		panic(err)
	}
	store, err := store.New("postgres", GetConfig().PostgresURL)
	if err != nil {
		panic(err)
	}
	addon.Store = store
	subRouter := a.Router.PathPrefix("/confluence/").Subrouter()
	subRouter.Use(middleware.NewRequestMiddleware(addon, make(map[string]string)))
	routes.RegisterRoutes(addon, subRouter)
	router := &ConfluenceRouter{
		Addon: addon,
	}
	router.setupRoutes(subRouter)
}

func (a *App) getAtlassianConfig() *os.File {
	path, err := filepath.Abs("./res/atlassian-config.json")
	if err != nil {
		panic(err)
	}
	config, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	return config
}

func (a *App) getAtlassianDescriptor() *os.File {
	path, err := filepath.Abs("./res/atlassian-connect.json")
	if err != nil {
		panic(err)
	}
	config, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	return config
}

func (a *App) InitializeTimers() {
	a.CleanupTicker = time.NewTicker(time.Minute * 5)
	go func() {
		for {
			<-a.CleanupTicker.C
			log.Println("Cleaning up expired database entries...")
			if err := GetAuthStateRepository().DeleteExpired(); err != nil {
				log.Println(err)
			}
			if err := GetSignupRepository().DeleteExpired(); err != nil {
				log.Println(err)
			}
		}
	}()
}

func (a *App) stripStaticPrefix(fs http.Handler, prefix string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r2 := new(http.Request)
		*r2 = *r
		r2.URL = new(url.URL)
		*r2.URL = *r.URL
		r2.URL.Path = "/"
		fs.ServeHTTP(w, r2)
	})
}

func (a *App) setupStaticAdminRoutes(router *mux.Router) {
	const basePath = "/admin"
	paths := []string{
		"/login",
		"/dashboard",
		"/locations",
		"/users",
		"/settings",
		"/bookings",
		"/search",
	}
	fs := http.FileServer(http.Dir(GetConfig().StaticAdminUiPath))
	for _, path := range paths {
		path = basePath + path
		router.PathPrefix(path).Handler(a.stripStaticPrefix(fs, path))
	}
	router.Path(basePath + "/").Handler(a.stripStaticPrefix(fs, basePath+"/"))
	router.PathPrefix(basePath + "/").Handler(http.StripPrefix(basePath+"/", fs))
}

func (a *App) setupStaticUserRoutes(router *mux.Router) {
	const basePath = "/ui"
	paths := []string{
		"/login",
		"/search",
		"/bookings",
	}
	fs := http.FileServer(http.Dir(GetConfig().StaticBookingUiPath))
	for _, path := range paths {
		path = basePath + path
		router.PathPrefix(path).Handler(a.stripStaticPrefix(fs, path))
	}
	router.Path(basePath + "/").Handler(a.stripStaticPrefix(fs, basePath+"/"))
	router.PathPrefix(basePath + "/").Handler(http.StripPrefix(basePath+"/", fs))
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
