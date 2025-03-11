package main

import (
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/runtime"
)

var once sync.Once
var instance runtime.Instance

func RunFlex(w http.ResponseWriter, r *http.Request) {
	runtime.CommonPreInit()
	once.Do(func() {
		instance = runtime.GetFlexRuntime()
	})

	err := runtime.RunFlex(instance)
	if err != nil {
		load.Logrus.WithError(err).Fatal("flex: failed to run runtime")
	}

	runtime.CommonPostInit()
}

func main() {
	load.Logrus.Info("main: starting")
	listenAddr := ":8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		listenAddr = ":" + val
	}

	http.HandleFunc("/RunFlex", RunFlex)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
