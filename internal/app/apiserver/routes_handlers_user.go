package apiserver

//Package of url handlers.

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/SKQR01/goblog/internal/app/model"
)

type ctxKey int8

const (
	sessionName        = "userSession"
	ctxKeyUser  ctxKey = iota
	ctxKeyRequestID
)

var (
	errIncorrectCredentials = errors.New("incorrect email or password")
	errNotAuth              = errors.New("not authenticated")
)

func (srv *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	//TODO:глянуть тут
	srv.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (srv *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func (srv *server) handleUsersCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			srv.error(w, r, http.StatusBadRequest, err)
			return
		}

		user := &model.User{
			Email:    req.Email,
			Password: req.Password,
		}
		if err := srv.store.User().Create(user); err != nil {
			srv.error(w, r, http.StatusUnprocessableEntity, err)
		}

		user.Sanitaze()
		srv.respond(w, r, http.StatusCreated, user)
	}
}

func (srv *server) handleUsersSessionsCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			srv.error(w, r, http.StatusBadRequest, err)
			return
		}

		user, err := srv.store.User().FindByEmail(req.Email)
		if err != nil || !user.ComparePassword(req.Password) {
			srv.error(w, r, http.StatusUnauthorized, errIncorrectCredentials)
			return
		}

		session, err := srv.sessionStore.Get(r, sessionName)
		if err != nil {
			srv.error(w, r, http.StatusInternalServerError, err)
			return
		}

		session.Values["user_id"] = user.ID
		if err := srv.sessionStore.Save(r, w, session); err != nil {
			srv.error(w, r, http.StatusInternalServerError, err)
			return
		}

		srv.respond(w, r, http.StatusOK, nil)
	}
}

func (srv *server) handleHome() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO:непонятка тут
		srv.respond(w, r, http.StatusOK, r.Context().Value(ctxKeyUser).(*model.User))
	})
}

// func (srv *server) createPostHandler() http.HandlerFunc {
// 	type request struct {
// 		Title   string `json:"title"`
// 		Content string `json:"content"`
// 		OwnerID int    `json:"owner"`
// 	}
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		req := &request{}
// 		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
// 			srv.error(w, r, http.StatusBadRequest, err)
// 			return
// 		}

// 		post := &model.Post{
// 			Title:   req.Title,
// 			Content: req.Content,
// 			OwnerID: req.OwnerID,
// 		}

// 		if err := srv.store.Post().Create(post); err != nil {
// 			srv.respond(w, r, http.StatusInternalServerError, err)
// 		}
// 		srv.respond(w, r, http.StatusOK, nil)
// 	}
// }

// func (srv *server) webSocketHandler() http.HandlerFunc {

// }

