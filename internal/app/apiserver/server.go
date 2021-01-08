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
	hub := newHub()

	srv.router.Use(srv.setRequestID)
	//TODO:logger conflicts with with websocket connetctions
	// srv.router.Use(srv.logRequest)

	headersOK := []string{"Accept", "Access-Control-Allow-Credentials",
		"Access-Control-Allow-Origin",
		"Access-Control-Allow-Headers",
		"X-Requested-With",
		"Content-Type",
		"Content-Length",
		"Accept-Encoding",
		"Authorization",
		"Set-Cookie",
		"Content-Length",
	}
	originsOK := []string{"https://localhost:3000"}

	srv.router.Use(
		handlers.CORS(
			handlers.AllowCredentials(),
			handlers.AllowedHeaders(headersOK),
			handlers.AllowedOrigins(originsOK),
		))
		
	srv.router.Handle("/users", srv.handleUsersCreate()).Methods("POST", "OPTIONS")
	srv.router.HandleFunc("/sessions", srv.handleUsersSessionsCreate()).Methods("POST", "OPTIONS")

	srv.router.HandleFunc("/posts-feed", srv.websocketPostHandler(hub))
	srv.router.HandleFunc("/posts/{id:[0-9]+}", srv.handlePostDetailView()).Methods("GET", "OPTIONS")

	private := srv.router.PathPrefix("/private").Subrouter()
	private.Use(srv.authenticateUser)
	private.HandleFunc("/home", srv.handleHome()).Methods("GET", "OPTIONS")
	private.HandleFunc("/home/posts", srv.handleGetUserPosts()).Methods("GET", "OPTIONS")
	private.HandleFunc("/sessions/logout", srv.handleUsersSessionsRemove()).Methods("POST", "OPTIONS")
	private.HandleFunc("/posts/create", srv.handlePostsCreate(hub)).Methods("POST", "OPTIONS")
	private.HandleFunc("/posts/remove", srv.handlePostsRemove(hub)).Methods("POST", "OPTIONS")
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
