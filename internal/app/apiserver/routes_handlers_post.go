package apiserver

//Package of url handlers.

import (
	"encoding/json"
	"net/http"

	"github.com/SKQR01/goblog/internal/app/model"
)

func (srv *server) createWebsocketPostHandler() http.HandlerFunc {
	//it`s not socket for a while
	type request struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		OwnerID int    `json:"owner"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			srv.error(w, r, http.StatusBadRequest, err)
			return
		}

		post := &model.Post{
			Title:   req.Title,
			Content: req.Content,
			OwnerID: req.OwnerID,
		}

		if err := srv.store.Post().Create(post); err != nil {
			srv.respond(w, r, http.StatusInternalServerError, err)
		}
		srv.respond(w, r, http.StatusOK, nil)
	}
}
