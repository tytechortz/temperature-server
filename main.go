package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

var (
	dataPath = flag.String("output", "/tmp/temperature.csv", "path to temperature file")
	port     = flag.String("port", "8080", "port upon which to listen")
	lock     sync.Mutex
)

func Save(w http.ResponseWriter, r *http.Request) {
	lock.Lock()
	f, err := os.Create(*dataPath)
	if err != nil {
		log.Println("couldn't open file for writing", err)
		return
	}

	defer func() {
		f.Close()
		lock.Unlock()
	}()

	var m map[string]float64
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		log.Println("couldn't write file", err)
		return
	}
	fmt.Fprintf(f, "%f\n", m["temperature"])
}

func Get(w http.ResponseWriter, r *http.Request) {
	lock.Lock()
	f, err := os.Open(*dataPath)
	if err != nil {
		log.Println("couldn't open file for reading", err)
		return
	}

	defer func() {
		f.Close()
		lock.Unlock()
	}()
	io.Copy(w, f)
}

func main() {
	fmt.Printf("Started server at http://localhost%v.\n", port)
	http.HandleFunc("/", Get)
	http.HandleFunc("/save", Save)
	fmt.Println(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", *port), nil))
}
