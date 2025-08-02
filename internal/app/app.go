package app

import (
    "net/http"
    "github.com/gorilla/mux"
)

type App struct {
    Router *mux.Router
}

func NewApp() *App {
    app := &App{
        Router: mux.NewRouter(),
    }
    app.initializeRoutes()
    return app
}

func (a *App) initializeRoutes() {
    a.Router.HandleFunc("/", a.homeHandler).Methods("GET")
}

func (a *App) homeHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Welcome to Arise-test!"))
}

func (a *App) Run(addr string) error {
    return http.ListenAndServe(addr, a.Router)
}