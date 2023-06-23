package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server"
	"github.com/hashicorp/terraform-provider-aws/internal/provider"
)

func main() {
	debugFlag := flag.Bool("debug", false, "Start provider in debug mode.")
	flag.Parse()

	serverFactory, _, err := provider.ProtoV5ProviderServerFactory(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		http.ListenAndServe("localhost:8080", nil)
	}()

	var serveOpts []tf5server.ServeOpt

	if *debugFlag {
		serveOpts = append(serveOpts, tf5server.WithManagedDebug())
	}

	logFlags := log.Flags()
	logFlags = logFlags &^ (log.Ldate | log.Ltime)
	log.SetFlags(logFlags)

	err = tf5server.Serve(
		"registry.terraform.io/hashicorp/aws",
		serverFactory,
		serveOpts...,
	)

	if err != nil {
		log.Fatal(err)
	}
}
