package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iho/bitly/shortener"
	"go.uber.org/fx"
)

func CreateURLHandler(c *gin.Context, s shortener.Shortener) {
	var request UrlCreationRequest
	if c.BindJSON(&request) == nil {
		key, err := s.Save(c.Request.Context(), request.Url)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, UrlCreationResponse{Code: key})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
	}
}

func GetURLHandler(c *gin.Context, s shortener.Shortener) {
	var request UrlGetRequest
	if err := c.ShouldBindUri(&request); err == nil {
		url, err := s.Load(c.Request.Context(), request.Code)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Redirect(http.StatusMovedPermanently, url)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "code is required"})
	}

}

func registerRoutes(r *gin.Engine, s shortener.Shortener) {
	r.POST("/urls", func(c *gin.Context) {
		CreateURLHandler(c, s)
	})
	r.GET("/urls/:code/", func(c *gin.Context) {
		GetURLHandler(c, s)
	})
}
func NewGinHTTPServer(lc fx.Lifecycle) *gin.Engine {
	return gin.Default()
}

// Module for go fx
var Module = fx.Options(fx.Invoke(registerRoutes), fx.Provide(NewGinHTTPServer))
