package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"transaction/internal/application"
)

/*
import (
	"net/http"

	"transaction/internal/adapters/http/rest"
)


func main() {

	logger := log.New(os.Stderr, "", log.Lshortfile)
	port := 8080
	logger.Printf("Starting server on port %d\n", port)
	router := rest.InitHandlers()
	logger.Fatalln(http.ListenAndServe(":8080", router))

}
*/

func main() {

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGTERM,
		os.Interrupt,
	)

	defer cancel()

	go application.Start(ctx)
	<-ctx.Done()

	application.Stop()
}
