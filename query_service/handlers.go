package main

import (
	"context"
	"github.com/guiaramos/go-meower/db"
	"github.com/guiaramos/go-meower/event"
	"github.com/guiaramos/go-meower/schema"
	"github.com/guiaramos/go-meower/search"
	"github.com/guiaramos/go-meower/util"
	"log"
	"net/http"
	"strconv"
)

var (
	ErrQueryMissing = "Missing query parameter"
	ErrInvalidSkip  = "Invalid skip parameters"
	ErrInvalidTake  = "Invalid take parameters"
	ErrListMeows    = "Could not fetch meows"
)

func onMeowCreated(m event.MeowCreatedMessage) {
	meow := schema.Meow{
		ID:        m.ID,
		Body:      m.Body,
		CreatedAt: m.CreatedAt,
	}
	if err := search.InsertMeow(context.Background(), meow); err != nil {
		log.Println(err)
	}
}

func searchMeowsHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx := r.Context()

	// Read params
	query := r.FormValue("query")
	if len(query) == 0 {
		util.ResponseError(w, http.StatusBadRequest, ErrQueryMissing)
		return
	}

	// Search settings
	skip := getSkip(w, r.FormValue("skip"))
	take := getTake(w, r.FormValue("take"))

	// Search meows
	meows, err := search.SearchMeows(ctx, query, skip, take)
	if err != nil {
		log.Println(err)
		util.ResponseOK(w, []schema.Meow{})
		return
	}

	util.ResponseOK(w, meows)
}

func listMeowsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error

	// Search settings
	skip := getSkip(w, r.FormValue("skip"))
	take := getTake(w, r.FormValue("take"))

	// Fetch meows
	meows, err := db.ListMeows(ctx, skip, take)
	if err != nil {
		log.Println(err)
		util.ResponseError(w, http.StatusInternalServerError, ErrListMeows)
		return
	}

	util.ResponseOK(w, meows)

}

func getTake(w http.ResponseWriter, takeStr string) (take uint64) {
	take = uint64(100)
	var err error
	if len(takeStr) != 0 {
		take, err = strconv.ParseUint(takeStr, 10, 64)
		if err != nil {
			util.ResponseError(w, http.StatusBadRequest, ErrInvalidTake)
			return
		}
	}
	return take
}

func getSkip(w http.ResponseWriter, skipStr string) (skip uint64) {
	skip = uint64(0)
	var err error
	if len(skipStr) != 0 {
		skip, err = strconv.ParseUint(skipStr, 10, 64)
		if err != nil {
			util.ResponseError(w, http.StatusBadRequest, ErrInvalidSkip)
			return
		}
	}
	return skip
}
