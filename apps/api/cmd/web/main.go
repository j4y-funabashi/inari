package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
	"github.com/julienschmidt/httprouter"
)

func monthsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	out := []app.Collection{}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(out)
}

func main() {
	router := httprouter.New()

	router.GET("/api/months", monthsHandler)

	log.Fatal(http.ListenAndServe(":8090", router))
}
