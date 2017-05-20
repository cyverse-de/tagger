package restapi

import (
	"crypto/tls"
	"net/http"
	"os"

	version "github.com/cyverse-de/version"
	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	swag "github.com/go-openapi/swag"
	graceful "github.com/tylerb/graceful"

	"github.com/cyverse-de/tagger/restapi/operations"
	"github.com/cyverse-de/tagger/restapi/operations/status"

	status_impl "github.com/cyverse-de/tagger/restapi/impl/status"
)

// This file is safe to edit. Once it exists it will not be overwritten

//go:generate swagger generate server --target .. --name tagger --spec ../swagger.yml

// Command line options that aren't managed by go-swagger.
var options struct {
	CfgPath     string `long:"config" default:"/etc/iplant/de/tagger.yaml" description:"The path to the config file"`
	ShowVersion bool   `short:"v" long:"version" description:"Print the app version and exit"`
}

func configureFlags(api *operations.TaggerAPI) {
	api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{
		swag.CommandLineOptionsGroup{
			ShortDescription: "Service Options",
			LongDescription:  "",
			Options:          &options,
		},
	}
}

// Initialize the service.
func initService() error {
	if options.ShowVersion {
		version.AppVersion()
		os.Exit(0)
	}

	return nil
}

// Clean up when the service exits.
func cleanup() {
}

func configureAPI(api *operations.TaggerAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.StatusGetHandler = status.GetHandlerFunc(status_impl.BuildStatusHandler(SwaggerJSON))

	api.ServerShutdown = cleanup

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *graceful.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json
// document.  So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return middleware.Redoc(middleware.RedocOpts{}, handler)
}
