package apiserver

//Package of url handlers.

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/SKQR01/goblog/internal/app/model"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gorilla/websocket"
)

var postsMessageChannel chan []byte = make(chan []byte, 1024)
var postsPaginationMessageChannel chan []byte = make(chan []byte, 1024)

// ------------------------------------Data types------------------------------------
// ----------------------------------------------------------------------------------
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
	Data   *model.Post
}

func (m *postsMessage) Validate() error {
	return validation.ValidateStruct(
		m,
		validation.Field(&m.Method, validation.Required),
		validation.Field(&m.Data, validation.Required),
	)
}

var (
	createMethod string = "CREATE"
	//On frontend I`ll use ReactJS, because in this case this method is unnecessary, but it can be useful in some other cases.
	removeMethod string = "REMOVE"
	editMethod   string = "EDIT"
)

// ----------------------------------------------------------------------------------
// ----------------------------------------------------------------------------------

// -----------------------------------------Handlers and helpers-----------------------------------------
// ------------------------------------------------------------------------------------------------------

func writeMessage(post *model.Post, method string) {

	encodedPostsMessage, err := json.Marshal(post)

	if err != nil {
		log.Println(err)
		return
	}
	postsMessageChannel <- encodedPostsMessage
}

func (srv *server) handlePostsCreate() http.HandlerFunc {
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

		go writeMessage(post, createMethod)

		srv.respond(w, r, http.StatusCreated, post)
	}
}

func (srv *server) handlePostsRemove() http.HandlerFunc {
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

		err = srv.store.Post().Remove(req.Ids)
		if err != nil {
			srv.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		srv.respond(w, r, http.StatusOK, req.Ids)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (srv *server) websocketPostHandler() http.HandlerFunc {

	return func(rw http.ResponseWriter, r *http.Request) {

		upgrader.CheckOrigin = func(r *http.Request) bool { return true }

		ws, err := upgrader.Upgrade(rw, r, nil)

		if err != nil {
			return
		}

		go func() {
			for {
				rangeMessage := &postsPaginationMessage{
					PageNumber:     1,
					PaginationSize: 10,
				}

				//Guarantee of valid
				if err := ws.ReadJSON(rangeMessage); err != nil {
					log.Println(err.Error())
					ws.WriteJSON(err)
					continue
				}

				if err := rangeMessage.Validate(); err != nil {
					log.Println(err.Error())
					ws.WriteJSON(err)
					continue
				}

				encodedRangeMessage, err := json.Marshal(rangeMessage)
				if err != nil {
					ws.WriteJSON(err)
					continue
				}

				postsPaginationMessageChannel <- encodedRangeMessage
			}
		}()

		go func() {
			for {
				rangeMessage := <-postsPaginationMessageChannel
				decodedRangeMessage := &postsPaginationMessage{}

				if err := json.Unmarshal(rangeMessage, &decodedRangeMessage); err != nil {
					ws.WriteJSON(err)
					continue
				}
				records, err := srv.store.Post().GetRecords(decodedRangeMessage.PageNumber, decodedRangeMessage.PaginationSize)
				if err != nil {
					ws.WriteJSON(err)
					continue
				}

				ws.WriteJSON(&records)
			}
		}()
		go func() {
			for {
				chanMessage := &postsMessage{}
				if err := json.Unmarshal(<-postsMessageChannel, &chanMessage); err != nil {
					ws.WriteJSON(err)
					continue
				}
				ws.WriteJSON(chanMessage)
			}
		}()
	}
}

// ------------------------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------------------------
