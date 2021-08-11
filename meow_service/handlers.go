package main

import (
	"github.com/guiaramos/go-meower/db"
	"github.com/guiaramos/go-meower/event"
	"github.com/guiaramos/go-meower/schema"
	"github.com/guiaramos/go-meower/util"
	"github.com/segmentio/ksuid"
	"html/template"
	"log"
	"net/http"
	"time"
)

var (
	ErrInvalidBody = "Invalid body"
	ErrFailCreate  = "Failed to create"
)

func createMeowHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		ID string `json:"id"`
	}

	ctx := r.Context()

	// Read params
	body := template.HTMLEscaper(r.FormValue("body"))
	if len(body) < 1 || len(body) > 140 {
		util.ResponseError(w, http.StatusBadRequest, ErrInvalidBody)
		return
	}

	// Create meow
	createdAt := time.Now().UTC()
	id, err := ksuid.NewRandomWithTime(createdAt)
	if err != nil {
		util.ResponseError(w, http.StatusInternalServerError, ErrFailCreate)
		return
	}

	meow := schema.Meow{
		ID:        id.String(),
		Body:      body,
		CreatedAt: createdAt,
	}
	if err := db.InsertMeow(ctx, meow); err != nil {
		log.Println(err)
		util.ResponseError(w, http.StatusInternalServerError, ErrFailCreate)
		return
	}

	// Publish event
	if err := event.PublishMeowCreated(meow); err != nil {
		log.Println(err)
	}

	// Return new meow
	util.ResponseOK(w, response{ID: meow.ID})
}
