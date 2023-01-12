package api

import (
	"github.com/gin-gonic/gin"
	"micze.io/mt/pkg/ad"
	"micze.io/mt/pkg/rabbitmq"
)

type Server struct {
	router    *gin.Engine
	rabbitmq  *rabbitmq.PublisherSubscriber
	ad        *ad.AD
	idcounter int64
}

// Start runs rest api server
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// NewServer creates new api server instance
func NewServer(publisher *rabbitmq.PublisherSubscriber, ad *ad.AD) *Server {
	server := &Server{}
	router := gin.Default()

	router.POST("/switchEndpoint", server.switchEndpoint)
	router.POST("/configureEndpoint", server.configureEndpoint)

	server.router = router
	server.rabbitmq = publisher
	server.ad = ad
	server.idcounter = 0
	return server
}

// errorResponse formats error messages
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
