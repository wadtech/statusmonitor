package main

import (
	"encoding/json"
	"flag"
	"github.com/wadtech/statusmonitor/config"
	"github.com/wadtech/statusmonitor/service"
	"html/template"
	"log"
	"net"
	"net/http"
	"time"
)

var configurator *config.Config
var work chan *service.Service

var htmlResponse = template.Must(template.New("htmlResp").Parse(templateHtml))

func main() {
	var configFile string

	flag.StringVar(&configFile, "c", "./config.json", "Path to config.json file")
	flag.Parse()

	configurator = config.NewConfig(configFile)

	log.Println("Configuration read from", configurator.Filepath)

	portFlag := configurator.Config.Port
	workerFlag := configurator.Config.Workers

	work = make(chan *service.Service, 10)
	for i := 0; i < workerFlag; i++ {
		go readFrom(work)
	}

	log.Println("ranging the services")
	for _, v := range configurator.Config.Services {
		log.Println("shoving ", v, " into channel")
		work <- &v
	}

	log.Println("adding the handler func")
	http.HandleFunc("/", check)
	http.HandleFunc("/favicon.ico", handleFavicon)

	listen := net.JoinHostPort("127.0.0.1", portFlag)
	log.Println("Now waiting on ", listen)

	log.Fatal(http.ListenAndServe(listen, nil))
}

func readFrom(work chan *service.Service) {
	select {
	case <-work:
		log.Println("hi9?")
		service := <-work
		service.Check()
		work <- service

		time.Sleep(10 * time.Second)
	}
}

func check(w http.ResponseWriter, r *http.Request) {
	log.Println("Status page requested")

	if r.Header.Get("Accept") == "application/json" {
		formatted, err := json.Marshal(configurator.Config.Services)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		tpl, err := template.New("json").Parse(string(formatted))
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		tpl.Execute(w, formatted)
	} else {
		if err := htmlResponse.Execute(w, configurator.Config.Services); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// 'handle'... aka dismiss without a second thought
func handleFavicon(w http.ResponseWriter, r *http.Request) {
	log.Println("Favicon requested")
	http.Error(w, "Not Found", http.StatusNotFound)
}

const templateHtml = `<html><head><title>Status Monitor</title></head>
<body>
  <h3>Status</h3>
  <ul>
    {{range .}}<li><strong>{{.Description}}</strong> {{.Ok}}</li>{{end}}
  </ul>
  <a href="http://github.com/wadtech/statusmonitor">Status monitor by Wadtech</a>
</body>
</html>`
