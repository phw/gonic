package main

import (
	"log"
	"net/http"
	"time"

	"github.com/sentriz/gonic/db"
	"github.com/sentriz/gonic/handler"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type middleware func(next http.HandlerFunc) http.HandlerFunc

func newChain(wares ...middleware) middleware {
	return func(final http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			last := final
			for i := len(wares) - 1; i >= 0; i-- {
				last = wares[i](last)
			}
			last(w, r)
		}
	}
}

func main() {
	address := ":5000"
	cont := handler.Controller{
		DB: db.New(),
	}
	withWare := newChain(
		cont.LogConnection,
		cont.EnableCORS,
		cont.CheckParameters,
	)
	mux := http.NewServeMux()
	mux.HandleFunc("/rest/ping", withWare(cont.Ping))
	mux.HandleFunc("/rest/ping.view", withWare(cont.Ping))
	mux.HandleFunc("/rest/stream", withWare(cont.Stream))
	mux.HandleFunc("/rest/stream.view", withWare(cont.Stream))
	mux.HandleFunc("/rest/getMusicDirectory", withWare(cont.GetMusicDirectory))
	mux.HandleFunc("/rest/getMusicDirectory.view", withWare(cont.GetMusicDirectory))
	mux.HandleFunc("/rest/getCoverArt", withWare(cont.GetCoverArt))
	mux.HandleFunc("/rest/getCoverArt.view", withWare(cont.GetCoverArt))
	mux.HandleFunc("/rest/getIndexes", withWare(cont.GetIndexes))
	mux.HandleFunc("/rest/getIndexes.view", withWare(cont.GetIndexes))
	mux.HandleFunc("/rest/getLicense", withWare(cont.GetLicence))
	mux.HandleFunc("/rest/getLicense.view", withWare(cont.GetLicence))
	mux.HandleFunc("/", withWare(cont.NotFound))
	server := &http.Server{
		Addr:         address,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	log.Printf("starting server at `%s`\n", address)
	err := server.ListenAndServe()
	if err != nil {
		log.Printf("when starting server: %v\n", err)
	}
}
