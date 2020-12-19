package apiserver

//Abstract server in a common implementation.

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"

	"github.com/SKQR01/goblog/internal/app/store"
)

type server struct {
	router *mux.Router
	logger *logrus.Logger
	//TODO:понять почему так, а не *store.Store
	store        store.Store
	sessionStore sessions.Store
}

func (srv *server) configureRouter() {
	//TODO:если что смотреть здесь
	srv.router.Use(srv.setRequestID)
	srv.router.Use(srv.logRequest)
	srv.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
	srv.router.HandleFunc("/users", srv.handleUsersCreate()).Methods("POST")
	srv.router.HandleFunc("/sessions", srv.handleUsersSessionsCreate()).Methods("POST")

	private := srv.router.PathPrefix("/private").Subrouter()
	private.Use(srv.authenticateUser)
	private.HandleFunc("/home", srv.handleHome()).Methods("GET")
}

func (srv *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.router.ServeHTTP(w, r)
}

func newServer(store store.Store, sessionsStore sessions.Store) *server {
	srv := &server{
		router:       mux.NewRouter(),
		logger:       logrus.New(),
		store:        store,
		sessionStore: sessionsStore,
	}
	srv.configureRouter()
	return srv
}
