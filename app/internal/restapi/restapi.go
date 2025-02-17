package restapi

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/PetStores/go-simple/internal/petstore/pet"

	"github.com/PetStores/go-simple/internal/petstore/category"

	"go.uber.org/zap"

	"github.com/gorilla/mux"
)

// RESTAPI represents a REST API business logic server.
type RESTAPI struct {
	server http.Server
	errors chan error
	logger *zap.SugaredLogger
}

// New returns a new instance of the REST API server.
func New(logger *zap.SugaredLogger, port int, categoryController *category.Controller, petController *pet.Controller) *RESTAPI {
	router := mux.NewRouter()

	peth := petHandlers{
		categoryController: categoryController,
		petController:      petController,
	}
	router.HandleFunc("/pet", peth.addPet()).Methods(http.MethodPost)
	//router.HandleFunc("/pet").Methods(http.MethodPut)

	return &RESTAPI{
		server: http.Server{
			Addr:    net.JoinHostPort("", strconv.Itoa(port)),
			Handler: router,
		},
		errors: make(chan error, 1),
		logger: logger,
	}
}

// Start diagnostics server.
func (rapi *RESTAPI) Start() {
	go func() {
		rapi.errors <- rapi.server.ListenAndServe()
		close(rapi.errors)
	}()
}

// Stop diagnostics server.
func (rapi *RESTAPI) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return rapi.server.Shutdown(ctx)
}

// Notify returns a channel to notify the caller about errors.
// If you receive an error from the channel you should stop the application.
func (rapi *RESTAPI) Notify() <-chan error {
	return rapi.errors
}
