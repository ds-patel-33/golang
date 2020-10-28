package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {

	// l := log.New(os.Stdout, "Product-api", log.LstdFlags)
	// hh := handlers.NewHello(l)
	// sm := http.NewServeMux()
	// sm.Handle("/", hh)

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		log.Println("Hello World")
		d, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(rw, "Ooopss", http.StatusBadRequest)
			// rw.WriteHeader(http.StatusBadRequest)
			// rw.Write([]byte("Oppppppsss"))
			return
		}

		// log.Printf("Data %s\n", d)
		fmt.Fprintf(rw, "Date %s\n", d)
	})

	http.HandleFunc("/goodbye", func(http.ResponseWriter, *http.Request) {
		log.Println("GoodBye World")
	})

	http.ListenAndServe(":9090", nil)
}
