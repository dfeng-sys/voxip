package main

import (
	"fmt"
	"net/http"

	mux "github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"voxip.com/voxip/voxip"
)

var (
	port string = "9001"
)

func init() {
	log.SetLevel(log.TraceLevel)
	log.SetFormatter(
		&log.TextFormatter{
			ForceColors:   true,
			ForceQuote:    true,
			FullTimestamp: true,
			PadLevelText:  true,
		},
	)
}

func main() {

	log.Info("main func")

	voxip.Init("db/GeoLite2-Country.mmdb")

	handleRequests(port)
}

func handleRequests(port string) {

	log.Info("listening for requests...")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/api/v1/ip", voxip.RequestHandler).Methods("GET")
	router.HandleFunc("/api/v1/update", voxip.UpdateHandler).Methods("PUT")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
