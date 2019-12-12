package main

import (
	"github.com/becosuke/guestbook-api/application/controller"
	"github.com/becosuke/guestbook-api/config"
	"log"
	"net/http"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime)
}

func main() {
	serveMux := http.NewServeMux()
	controller.Register(serveMux)

	log.Fatal(http.ListenAndServe(config.GetConfig().Rest.Addr, serveMux))
}
