[![Go](https://github.com/dronezzzko/easyshutdown/actions/workflows/linters.yml/badge.svg)](https://github.com/dronezzzko/easyshutdown/actions/workflows/linters.yml)

# easyshutdown
Gracefully shutdown your Go services in just one line. 

This package supports: 
- HTTP and HTTPS servers
- GRPC server
- OpenTelemetry tracers
- And more

## Usage
### Simple example
```go
package main

import (
	"log"
	"net/http"

	es "github.com/dronezzzko/easyshutdown"
)

func main() {
	srv := &http.Server{
		Addr: ":8080",
	}

	go func() {
		log.Println("starting HTTP server", srv.Addr)

		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("HTTP server stopped: %s\n", err.Error())
		}
	}()

	sd, _ := es.NewShutdown(
		es.WithHTTPServer(srv),
	)
	sd.Graceful()
}
```

```terminal
go run .
```

Press ``ctrl+c``:
```
2023/02/04 19:25:04 starting HTTP server :8080
easyshutdown 2023/02/04 19:25:05 Shutting down HTTP/HTTPS server
2023/02/04 19:25:08 HTTP server stopped: http: Server closed
exit status 1
```

Also see [options.go](options.go) for all available options and supported services.
