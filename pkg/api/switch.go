package api

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"micze.io/mt/pkg/rabbitmq"
)

// switchEndpointRequest request data from endpoint
type switchEndpointRequest struct {
	Mac      string `json:"mac,omitempty" binding:"required"`
	IpAddr   string `json:"ipaddr,omitempty" binding:"required"`
	SysId    string `json:"sysid,omitempty" binding:"required"`
	Username string `json:"username,omitempty" binding:"required"`
}

// switchEndpoint handles endpoint requests to switch
func (server *Server) switchEndpoint(c *gin.Context) {
	var req switchEndpointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("error parsing request data:", err.Error())
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	msg := rabbitmq.SwitchEndpointParams{
		Mac:       req.Mac,
		IpAddr:    req.IpAddr,
		SysId:     req.SysId,
		Username:  req.Username,
		Timestamp: time.Now().Unix(),
		Id:        server.idcounter,
	}

	log.Printf("[M%06d] received switch message from \"%v\": %v\n", msg.Id, c.RemoteIP(), msg)
	log.Printf("[M%06d] publishing task to unprepared queue: %v\n", msg.Id, msg)
	err := server.rabbitmq.Publish(msg)
	if err != nil {
		log.Printf("[M%06d] error publishing task %v: %v\n", msg.Id, msg, err.Error())
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	server.idcounter++
	log.Printf("[M%06d] completed processing message: %v\n", msg.Id, msg)

	c.IndentedJSON(http.StatusOK, msg)
}
