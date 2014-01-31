package main

import (
	"encoding/json"
	"flag"
	"github.com/wadtech/statusmonitor/service"
	"html/template"
	"log"
	"net"
	"net/http"
)

var portFlag = flag.String("port", "8080", "http service port")
var workerFlag = flag.Int("workers", 10, "Number of concurrent workers testing ports")

var myService = service.NewService("Port 8080 on localhost", "localhost", "8080")

func main() {
	flag.Parse()

	//@todo spin up workers based on settings given (hardcoded at first)

	go myService.Check()

	http.HandleFunc("/", check)
	log.Println("Now waiting on 127.0.0.1", *portFlag)
	log.Fatal(http.ListenAndServe(net.JoinHostPort("127.0.0.1", *portFlag), nil))
}

func check(w http.ResponseWriter, r *http.Request) {
	accept := r.Header.Get("Accept")

	if accept == "application/json" {
		formatted, err := json.Marshal(myService)
		if err != nil {
			log.Println("Couldn't Marshal Service", myService.Description)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

		t, err := template.New("json").Parse(string(formatted))
		if err != nil {
			log.Println("Couldn't Render JSON Template")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		t.Execute(w, myService)
	} else {
		//@todo render html template
		t, err := template.New("html").Parse(html)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

		t.Execute(w, myService)
	}
}

var html = `<html><head><title>Status Monitor</title></head>
<body>
  <h3>Status</h3>
  <ul>
    <li><strong>{{.Description}}</strong> {{.Ok}}</li>
  </ul>
  <a href="http://github.com/wadtech/statusmonitor">Status monitor by Wadtech</a>
</body>
</html>`
