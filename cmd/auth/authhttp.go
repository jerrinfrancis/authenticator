package main

import (
	"io"
	"log"
	"net/http"

	"github.com/jerrinfrancis/authenticator/user"
)

var mux map[string]func(http.ResponseWriter, *http.Request)

type myHandler struct{}

func (h *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h, ok := mux[r.URL.String()]; ok {
		h(w, r)
		return
	}
	io.WriteString(w, "My Server: "+r.URL.String())
}

func main() {
	server := http.Server{
		Addr:    ":8080",
		Handler: &myHandler{},
	}
	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	mux["/login"] = user.LoginHandler
	mux["/register"] = user.RegisterHandler
	log.Println("Server listening at : ", server.Addr)
	server.ListenAndServe()

}
