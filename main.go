package main

import (
	"net/http"
	"os"
	"log"
	"context"
	"os/signal"
	"time"
	"io"
	"github.com/kavu/go_reuseport"
	"fmt"
	"strconv"
)

func main() {

	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	listener, err := reuseport.Listen("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	mux := http.NewServeMux()

	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request!")
		time.Sleep(5 * time.Second)
		io.WriteString(w, "Finished!")
		fmt.Println("Request -done-")
	}))

	srv := &http.Server{Handler: mux}

	go func() {
		// service connections
		fmt.Println("Listening as: ",strconv.Itoa(os.Getpid()), " on Port: 8080")
		if err := srv.Serve(listener); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()

	<-stopChan // wait for SIGINT
	log.Println("Shutting down server...")

	// shut down gracefully, but wait no longer than 5 seconds before halting
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	srv.Shutdown(ctx)

	log.Println("Server gracefully stopped")
}
