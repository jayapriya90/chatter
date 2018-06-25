package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/jayapriya90/chatter/backend"
	log "github.com/sirupsen/logrus"
)

var port = flag.Int("port", 8888, "port to run the chatter service")

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true, 
		TimestampFormat: time.RFC3339Nano,
	})
	log.SetLevel(log.DebugLevel)
	flag.Parse()
	server := backend.NewServer()
	go server.Run()
	http.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
		backend.ServeWebSocket(server, w, r)
	})
	log.Infof("Listening and serving on %d", *port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}