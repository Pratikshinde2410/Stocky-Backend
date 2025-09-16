package handlers

import "net/http"

func StatsHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
}


