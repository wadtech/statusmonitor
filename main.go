package main

import (
	"encoding/json"
	"flag"
	"github.com/wadtech/statusmonitor/checker"
	"github.com/wadtech/statusmonitor/service"
	"html/template"
	"log"
	"net"
	"net/http"
)

var portFlag string
var workerFlag int
var configFile string

var myService = service.NewService("Port 8080 on localhost", "localhost", "8080")

func main() {
	flag.StringVar(&portFlag, "p", "8080", "http service port")
	flag.IntVar(&workerFlag, "w", 10, "Number of concurrent workers")
	flag.StringVar(&configFile, "c", "config.json", "Path to config.json file")

	flag.Parse()

	ch := make(chan *service.Service)
	for i := 0; i < workerFlag; i++ {
		go checker.Listen(ch)
	}

	//@todo this will become a looped thing based on how many jobbers are in the whatsit
	ch <- myService

	http.HandleFunc("/", check)
	http.HandleFunc("/favicon.ico", handleFavicon)

	log.Println("Now waiting on 127.0.0.1", portFlag)
	log.Fatal(http.ListenAndServe(net.JoinHostPort("127.0.0.1", portFlag), nil))
}

func check(w http.ResponseWriter, r *http.Request) {
	log.Println("Status page requested")

	if r.Header.Get("Accept") == "application/json" {
		formatted, err := json.Marshal(myService)
		if err != nil {
			log.Println("Couldn't Marshal Service", myService.Description)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

		tpl, err := template.New("json").Parse(string(formatted))
		if err != nil {
			log.Println("Couldn't Render JSON Template")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

		log.Println("Rendering json response")
		w.Header().Set("Content-Type", "application/json")
		tpl.Execute(w, myService)
	} else {
		tpl, err := template.New("html").Parse(html)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

		log.Println("Rendering html template")
		tpl.Execute(w, myService)
	}
}

func handleFavicon(w http.ResponseWriter, r *http.Request) {
	log.Println("Favicon requested")
	http.Error(w, "Not Found", http.StatusNotFound)
}

// @todo move the template somewhere else.
var html = `<html><head><title>Status Monitor</title></head>
<body>
  <h3>Status</h3>
  <ul>
    <li><strong>{{.Description}}</strong> {{.Ok}}</li>
  </ul>
  <a href="http://github.com/wadtech/statusmonitor">Status monitor by Wadtech</a>
</body>
</html>`
