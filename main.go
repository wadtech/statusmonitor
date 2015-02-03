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
	"os"
	"time"
)

var configurator *config.Config
var work = make(chan *service.Service, 10)
var verbose bool

var htmlResponse = template.Must(template.New("htmlResp").Parse(templateHtml))

func main() {
	var configFile string

	flag.StringVar(&configFile, "c", "./config.json", "Path to config.json file")
	flag.BoolVar(&verbose, "v", false, "Verbose logging information")
	flag.Parse()

	configurator, err := config.NewConfig(configFile)
	if err != nil {
		log.Printf("Failed to read config from '%s': %s\n", configFile, err.Error())
		os.Exit(1)
	}

	logMessage("Configuration read from", configurator.Filepath)

	portFlag := configurator.Config.Port
	workerFlag := configurator.Config.Workers

	for i, _ := range configurator.Config.Services {
		work <- &configurator.Config.Services[i]
	}

	for i := 0; i < workerFlag; i++ {
		go readFrom(work)
	}

	http.HandleFunc("/", check)
	http.HandleFunc("/favicon.ico", handleFavicon)

	listen := net.JoinHostPort("127.0.0.1", portFlag)
	logMessage("Now waiting on ", listen)

	log.Fatal(http.ListenAndServe(listen, nil))
}

func readFrom(work chan *service.Service) {
	for {
		select {
		case service := <-work:
			service.Check()
			logMessage(service.Description, service.Ok)
			work <- service
		default:
		}
		time.Sleep(5 * time.Second)
	}
}

func check(w http.ResponseWriter, r *http.Request) {
	logMessage("Status page requested")

	if r.Header.Get("Accept") == "application/json" {
		formatted, err := json.Marshal(configurator.Config.Services)
		if err != nil {
			logMessage(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		tpl, err := template.New("json").Parse(string(formatted))
		if err != nil {
			logMessage(err)
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
	logMessage("Favicon requested")
	http.Error(w, "Not Found", http.StatusNotFound)
}

func logMessage(message ...interface{}) {
	if verbose {
		log.Println(message)
	}
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
