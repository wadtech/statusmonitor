package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/wadtech/statusmonitor/service"
	"log"
	"net/http"
)

//@todo move this out
var workers = 10

var port = flag.String("port", ":8080", "http service port")

var myService = service.NewService("Port 8080 on localhost", "localhost", "8080")

func main() {
	//spin up workers based on settings given (hardcoded at first)

	go myService.Check()

	http.HandleFunc("/", check)
	log.Println("Now waiting on 127.0.0.1", *port)
	log.Fatal(http.ListenAndServe(*port, nil))
}

func check(w http.ResponseWriter, r *http.Request) {
	json := false

	accept := r.Header.Get("Accept")
	if accept == "application/json" {
		json = true
	}

	//@todo render json template if json = true, else html template from status struct
	var buf bytes.Buffer

	buf.WriteString(myService.Description)
	buf.WriteString("\n")
	buf.WriteString(fmt.Sprintf("%t", myService.Ok))
	buf.WriteString("\n")
	buf.WriteString("\n")
	buf.WriteString("Statusmonitor by Wadtech 2014\n")
	if json {
		buf.WriteString("JSON response requested\n")
	}
	buf.WriteTo(w)
}
