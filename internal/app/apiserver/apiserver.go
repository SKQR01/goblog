package apiserver

import (
	"database/sql"
	"github.com/SKQR01/goblog/internal/app/store/sqlstore"
	"github.com/gorilla/sessions"
	"net/http"
)

func Start(config *Config) error {
	db, err := newDB(config.DatabaseURL)
	if err != nil{
		return nil
	}

	defer db.Close()

	store := sqlstore.New(db)
	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))
	srv := newServer(store, sessionStore)

	return http.ListenAndServe(config.BindAddr, srv)
}

func newDB(databasURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databasURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err !=nil{
		return nil, err
	}

	return db, err
}