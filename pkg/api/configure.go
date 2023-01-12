package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type userConfigRequest struct {
	Username string `json:"username,omitempty" binding:"required"`
}

type userConfig struct {
	Username string   `json:"username,omitempty" binding:"required"`
	Groups   []string `json:"groups,omitempty" binding:"required"`
}

func (server *Server) configureEndpoint(c *gin.Context) {
	var req userConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	groups := server.ad.GetUserGroups("ble")
	userConfig := &userConfig{
		Username: req.Username,
		Groups:   groups,
	}

	c.IndentedJSON(http.StatusOK, userConfig)
}
