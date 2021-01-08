package apiserver

//Physical realization of api server.

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/sessions"

	"github.com/SKQR01/goblog/internal/app/store/sqlstore"
)

//Start ...
func Start(config *Config) error {
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}

	defer db.Close()

	store := sqlstore.New(db)
	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))
	sessionStore.Options.HttpOnly = true
	sessionStore.Options.SameSite = http.SameSite(http.SameSiteNoneMode)
	sessionStore.Options.Secure = true

	srv := newServer(store, sessionStore)

	
	srv.logger.Println("Starting server...")
	//for https (requires ssl)
	return http.ListenAndServeTLS(config.BindAddr, "certs/localhost.pem", "certs/localhost-key.pem", srv)
}

func newDB(databasURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databasURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, err
}
