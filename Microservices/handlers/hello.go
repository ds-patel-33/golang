package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Hello struct {
	l *log.Logger
}

func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

func (h *Hello) ServerHttp(rw http.ResponseWriter, r *http.Request) {
	h.l.Println("Hello World")
	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "Ooopss", http.StatusBadRequest)
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Oppppppsss"))
		return
	}

	log.Printf("Data %s\n", d)
	fmt.Fprintf(rw, "Date %s\n", d)
}
