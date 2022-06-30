package main

import (
	"context"
	"fmt"
	"imageuploader/config"
	apihandlers "imageuploader/internal/apihandlers"
	imageservice "imageuploader/internal/image"
	imagerepository "imageuploader/internal/image/repository"
	"imageuploader/pkg/utils/files"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

const (

	// Timeouts to prevent slow-loris attacks and for avoiding build-ups
	serverReadTimeout  = 10 * time.Second
	serverIdleTimeout  = 10 * time.Second
	serverWriteTimeout = 10 * time.Second
)

func initImageHandler(ctx context.Context) (*imageservice.Handler, error) {
	imageRepo, err := imagerepository.NewBoltRepository(config.LocalMemPath)
	if err != nil {
		return nil, err
	}

	imageService := imageservice.NewService(imageRepo)

	return imageservice.NewHandlers(imageService), nil
}

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	files.CreateDir(config.LocalMemPath)
	files.CreateDir(config.LocalImagesDirectory)

	imageHandlers, err := initImageHandler(ctx)
	if err != nil {
		log.Fatalf("Error initialising image service %v", err)
	}
	api := apihandlers.Api{
		ImageHandlers: imageHandlers,
	}

	r := gin.New()
	api.CreateRoutes(r)
	r.Static("/home", "./public")
	http.Handle("/", r)
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	// Set up the http server
	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%s", config.HttpServerPort),
		ReadTimeout:  serverReadTimeout,
		IdleTimeout:  serverIdleTimeout,
		WriteTimeout: serverWriteTimeout,
	}

	go srv.ListenAndServe()

	// Block for exit
	<-termChan

	// Cleanup
	srv.Shutdown(ctx)
	cancelFunc()
}
