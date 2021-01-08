package apiserver

//Package of url handlers.

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/SKQR01/goblog/internal/app/model"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// ------------------------------------Data types------------------------------------
// ----------------------------------------------------------------------------------

var (
	createMethod     string = "CREATE"
	removeMethod     string = "REMOVE"
	editMethod       string = "EDIT"
	errorMethod      string = "ERROR"
	paginationMethod string = "PAGINATION"
)

type postsPaginationMessage struct {
	PageNumber     int `json:"pageNumber"`
	PaginationSize int `json:"paginationSize"`
}

func (m *postsPaginationMessage) Validate() error {
	return validation.ValidateStruct(
		m,
		validation.Field(&m.PageNumber, validation.Required),
		validation.Field(&m.PaginationSize, validation.Required),
	)
}

type postsMessage struct {
	Method string
	Data   interface{}
}

//Posts websocket view conn, home posts view, improve get records and find commands, secure connection, auth

// ----------------------------------------------------------------------------------
// ----------------------------------------------------------------------------------

// -----------------------------------------Handlers and helpers-----------------------------------------
// ------------------------------------------------------------------------------------------------------

func (srv *server) handleGetUserPosts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rangeMessage := &postsPaginationMessage{
			PageNumber:     1,
			PaginationSize: 10,
		}
		keys := r.URL.Query()

		if keys["pageNumber"] != nil {
			pageNumber, err := strconv.Atoi(keys["pageNumber"][0])
			if err == nil {
				rangeMessage.PageNumber = pageNumber
			}
		}

		if keys["paginationSize"] != nil {
			paginationSize, err := strconv.Atoi(keys["paginationSize"][0])
			if err == nil {
				rangeMessage.PaginationSize = paginationSize
			}
		}

		records, err := srv.store.Post().GetRecords(rangeMessage.PageNumber,
			rangeMessage.PaginationSize,
			r.Context().Value(ctxKeyUser).(*model.User).ID,
		)
		if err != nil {
			srv.error(w, r, http.StatusInternalServerError, err)
			return
		}
		srv.respond(w, r, http.StatusOK, records)
	}
}

func (srv *server) handlePostsCreate(hub *Hub) http.HandlerFunc {
	type request struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}

		var err error = nil

		err = json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			srv.error(w, r, http.StatusBadRequest, err)
			return
		}

		post := model.NewPost()

		post.Title = req.Title
		post.Content = req.Content
		post.SetOwnerID(r.Context().Value(ctxKeyUser).(*model.User).ID)

		err = srv.store.Post().Create(post)
		if err != nil {
			srv.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		//send message to hub (to make users possible to see changes)
		chanMessage := &postsMessage{
			Method: createMethod,
			Data:   post,
		}

		encodedChanMessage, _ := json.Marshal(chanMessage)

		hub.broadcast <- encodedChanMessage

		srv.respond(w, r, http.StatusCreated, post)
	}
}

func (srv *server) handlePostDetailView() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			srv.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		post, err := srv.store.Post().Find(id)
		if err != nil {
			srv.error(w, r, http.StatusInternalServerError, err)
			return
		}
		srv.respond(w, r, http.StatusOK, post)
	}
}

func (srv *server) handlePostsRemove(hub *Hub) http.HandlerFunc {
	type request struct {
		Ids []int `json:"postsToRemoveIds"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}

		var err error = nil

		err = json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			srv.error(w, r, http.StatusBadRequest, err)
			return
		}

		err = srv.store.Post().Remove(req.Ids, r.Context().Value(ctxKeyUser).(*model.User).ID)
		if err != nil {
			srv.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		//send message to hub (to make users possible to see changes)
		chanMessage := &postsMessage{
			Method: removeMethod,
			Data:   req.Ids,
		}

		encodedChanMessage, _ := json.Marshal(chanMessage)

		hub.broadcast <- encodedChanMessage

		srv.respond(w, r, http.StatusOK, req.Ids)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (srv *server) websocketPostHandler(hub *Hub) http.HandlerFunc {

	return func(rw http.ResponseWriter, r *http.Request) {

		go hub.run()

		upgrader.CheckOrigin = func(r *http.Request) bool { return true }

		conn, err := upgrader.Upgrade(rw, r, nil)

		if err != nil {
			log.Println(err)
			return
		}

		client := &Client{
			hub:  hub,
			conn: conn,
			send: make(chan []byte, 256),
		}
		client.hub.register <- client

		go client.writePump()
		go client.readPump(srv)
	}
}

// ------------------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------------------
