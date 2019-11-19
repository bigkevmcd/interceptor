package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/bigkevmcd/interceptor/pkg/interception/pullrequest"
)

var (
	port = flag.Int("port", 8080, "port to listen on")
)

func main() {
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		match, err := pullrequest.MatchPullRequestAction(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !match {
			fmt.Printf("pull-request match: %v\n", match)
		}
	})

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
