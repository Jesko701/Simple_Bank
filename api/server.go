package api

import (
	"solo_simple-bank_tutorial/db/sqlc"
	"solo_simple-bank_tutorial/util"

	"github.com/gin-gonic/gin"
)

type Server struct {
	config util.Config
	store  sqlc.Store
	router *gin.Engine
}

func NewServer(store sqlc.Store) (*Server, error) {
	router := gin.Default()
	server := &Server{
		store:  store,
		router: router,
	}

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome to the API!"})
	})

	server.router = router
	return server, nil
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}
