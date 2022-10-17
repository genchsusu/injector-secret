package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang/glog"
)

func main() {
	parameters := DefaultParametersObject()

	// get command line parameters
	flag.IntVar(&parameters.port, "port", LookupIntEnv("CONFIG_PORT", parameters.port), "Webhook server port.")
	flag.StringVar(&parameters.excludeNamespaces, "excludeNamespaces", LookupStringEnv("CONFIG_EXCLUDE_NAMESPACES", parameters.excludeNamespaces), "Comma-separated namespace names to ignore.")
	flag.StringVar(&parameters.targetSecretName, "targetSecretName", LookupStringEnv("CONFIG_TARGET_SECRET_NAME", parameters.targetSecretName), "Name of the targetSecret secret we will create in the namespace")

	flag.Parse()

	glog.Infof("Running with config: %+v", parameters)

	whsvr, err := NewWebhookServer(
		&parameters,
		&http.Server{
			Addr: fmt.Sprintf(":%v", parameters.port),
		},
	)
	if err != nil {
		glog.Exitf("Could not create the Webhook server: %v", err)
	}

	// define http server and server handler
	mux := http.NewServeMux()
	mux.HandleFunc("/mutate", whsvr.serve)
	whsvr.server.Handler = mux

	// start webhook server in new rountine
	go func() {
		if err := whsvr.server.ListenAndServe(); err != nil {
			glog.Errorf("Failed to listen and serve webhook server: %v", err)
		}
	}()

	// listening OS shutdown singal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	glog.Infof("Got OS shutdown signal, shutting down webhook server gracefully...")
	if err := whsvr.server.Shutdown(context.Background()); err != nil {
		glog.Errorf("Error while shutting down: %v", err)
	}
}
