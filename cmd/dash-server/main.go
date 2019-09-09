// dash-server is the dash command line server.
//
// Usage:
//
//    dash-server [-datadir <datadir>]
//
// The server will listen for incoming DASH experiment requests and
// will keep serving them until it is interrupted.
//
// The `-datadir <datadir>` flag specifies the directory where to write
// measurement results. By default is the current working directory.
//
// The server will listen on `:80`.
//
// The server will emit access logs on the standard output using the
// usual format. The server will emit error logging on the standard
// error using github.com/apex/log's JSON format.
package main

import (
	"context"
	"flag"
	"net/http"
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/json"
	"github.com/gorilla/handlers"
	"github.com/m-lab/go/prometheusx"
	"github.com/m-lab/go/rtx"
	"github.com/neubot/dash/server"
)

var (
	flagDatadir  = flag.String("datadir", ".", "directory where to save results")
)

func main() {
	log.Log = &log.Logger{
		Handler: json.New(os.Stderr),
		Level:   log.DebugLevel,
	}
	flag.Parse()
	promServer := prometheusx.MustServeMetrics()
	defer promServer.Close()
	mux := http.NewServeMux()
	handler := server.NewHandler(*flagDatadir)
	handler.StartReaper(context.Background())
	handler.RegisterHandlers(mux)
	handler.Logger = log.Log
	loggingHandler := handlers.LoggingHandler(os.Stdout, mux)
	rtx.Must(
		http.ListenAndServe(":80", loggingHandler),
		"ListenAndServe failed",
	)
}
