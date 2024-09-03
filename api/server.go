package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"solo_simple-bank_tutorial/db/sqlc"
	"solo_simple-bank_tutorial/token"
	"solo_simple-bank_tutorial/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config util.Config
	store  sqlc.Store
	router *gin.Engine
	token  token.Maker
}

func NewServer(config util.Config, store sqlc.Store) (*Server, error) {
	// If using JWT, create the ecdsa file private key first and load it
	// using paseto = directly using TokenAPI

	// TODO: Readfile for private key (optional: if you're using the file)
	privateKeyBytes, err := ioutil.ReadFile("ecdsa_private_key.pem")
	if err != nil {
		return nil, err
	}

	// TODO: ReadEnv for privateKey (optional: generally used for secret key (github,etc))
	// privateKeyBytes, exists := os.LookupEnv("ECDSA_PRIVATE_KEY")
	// if !exists {
	// 	fmt.Println("Private Key environtment variable is not set")
	// }

	// if using paseto, change to token.NewPasetoMaker and used config.TokenAPI
	// and delete the read file above
	token, err := token.NewJWTMaker(string(privateKeyBytes))
	if err != nil {
		return nil, fmt.Errorf("cannot creat token maker: %v", err)
	}
	router := gin.Default()
	server := &Server{
		store:  store,
		router: router,
		token:  token,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.Routes(router)

	server.router = router
	return server, nil
}

func (s *Server) Routes(router *gin.Engine) {
	user := router.Group("/users")
	{
		user.POST("", s.CreateUser)
		user.POST("login", s.LoginUser)
	}

	hello := router.Group("/hello").Use(authMiddleware(s.token))
	{
		hello.GET("/", func(c *gin.Context) {
			payload, exists := c.Get(Authorization_Payload)
			if !exists {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			}
			c.JSON(200, gin.H{"message": "Welcome to the API!",
				"user": payload})
		}).Use(authMiddleware(s.token))
	}
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}
