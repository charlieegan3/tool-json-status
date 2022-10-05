package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/doug-martin/goqu/v9"
)

func BuildLatestHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	goquDB := goqu.New("postgres", db)

	return func(w http.ResponseWriter, r *http.Request) {
		sel := goquDB.From("jsonstatus.data").Select("value").Where(goqu.C("key").Eq("status")).Limit(1)
		var statusJSON string
		found, err := sel.Executor().ScanVal(&statusJSON)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("failed to get current status: %s", err)))
			return
		}

		if !found {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// set the content type to JSON
		w.Header().Set("Content-Type", "application/json")

		w.Write([]byte(statusJSON))
	}
}
