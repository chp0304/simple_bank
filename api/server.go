package api

import (
	db "github.com/chp0304/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server is the HTTP server for the simple bank API
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server with the given store
func NewServer(store db.Store) *Server {
	router := gin.Default()
	server := &Server{
		store:  store,
		router: router,
	}
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)
	router.POST("/transfers", server.createTransfer)
	return server
}

// Start runs the HTTP server on a speacific address
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

// errorResponse returns a standardized error response
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
