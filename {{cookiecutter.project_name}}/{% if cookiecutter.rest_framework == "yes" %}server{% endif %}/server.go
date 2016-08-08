// Echo Server

package server

import (
	"fmt"
	"net/http"

	"{{ cookiecutter.project_name|lower }}/config"
	"{{ cookiecutter.project_name|lower }}/hal"
	"{{ cookiecutter.project_name|lower }}/logger"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/fasthttp"
)

// Configuration
var c config.ServerConfig

// Always initialise an Echo instance
func New() *echo.Echo {
	server := echo.New()
	// Add healthcheck handler
	server.GET("/__healthcheck__", healthcheckHandler)
	// Set default error handler
	server.SetHTTPErrorHandler(errorHandler)

	return server
}

func RunServer(server *echo.Echo) {
	// Runs the web service
	logger.Fields(logger.F{"config": c}).Info("Starting Web Server")
	server.Run(fasthttp.New(fmt.Sprintf("%s:%d", c.Bind, c.Port)))
}

// Updates the c variable with configuration
func Configure() {
	c = config.NewServerConfig()
}

// Common error handler func that returns proper formatted HAL documents
func errorHandler(err error, ctx echo.Context) {
	code := http.StatusInternalServerError
	msg := http.StatusText(code)
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message
	}
	switch code {
	case http.StatusNotFound:
		hal.NotFound(ctx)
	case http.StatusInternalServerError:
		hal.ServerError(ctx, err)
	default:
		hal.Response(code, ctx, hal.Resource(ctx, &hal.Error{Message: msg}))
	}
}
