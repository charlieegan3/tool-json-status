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
		w.Header().Set("Access-Control-Allow-Origin", "https://charlieegan3.com")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")
		if r.Method == "OPTIONS" {
			return
		}

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

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(statusJSON))
	}
}
