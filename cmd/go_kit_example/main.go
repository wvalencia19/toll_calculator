package main

import (
	"net"
	"net/http"
	"os"

	"github.com/go-kit/log"
	"github.com/wvalencia19/tolling/cmd/go_kit_example/aggsvc/aggendpoint"
	"github.com/wvalencia19/tolling/cmd/go_kit_example/aggsvc/aggservice"
	"github.com/wvalencia19/tolling/cmd/go_kit_example/aggsvc/aggtransport"
)

func main() {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	service := aggservice.New()
	endpoints := aggendpoint.New(service, logger)
	httpHandler := aggtransport.NewHTTPHandler(endpoints, logger)

	httpAddr := ":3001"

	{
		// The HTTP listener mounts the Go kit HTTP handler we created.
		httpListener, err := net.Listen("tcp", httpAddr)
		if err != nil {
			logger.Log("transport", "HTTP", "during", "Listen", "err", err)
			os.Exit(1)
		}

		logger.Log("transport", "HTTP", "addr", httpAddr)
		err = http.Serve(httpListener, httpHandler)

		if err != nil {
			panic(err)
		}

	}
}
