package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cmullendore/ephem/src/core"
	"github.com/cmullendore/ephem/src/ephem"
	"github.com/cmullendore/ephem/src/webui"
)

func main() {

	log.Println("Starting...")
	log.Println(os.Getwd())

	var conf = LoadMainConfig()

	var se core.ISecretsEngine

	switch engine := conf.SecretsEngine; engine {
	case "ephem":
		se = ephem.CreateEngine()
	case "vault":
		//se =
	}

	var wi core.IWebInterface

	switch mode := conf.Mode; mode {
	case "webui":
		wi = webui.NewUIServer(se)
	case "grpc":
		//wi =
	case "rest":
		//wi =
	}

	wi.Listen()

	var interrupt = make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(interrupt)

	<-interrupt

}
