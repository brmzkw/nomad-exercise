package webservice

import (
	"net/http"
	"regexp"

	"github.com/brmzkw/nomad-exercise/internal/nomad"
	"github.com/gin-gonic/gin"
)

type WebService struct {
	nomad *nomad.Nomad
}

func NewWebService(nomad *nomad.Nomad) *WebService {
	return &WebService{nomad: nomad}
}

type PageRequest struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	Script bool   `json:"script"`
}

type PageResponse struct {
	URL string `json:"url"`
}

func (ws *WebService) createPage(c *gin.Context) {
	var service PageRequest
	if err := c.BindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// name should only be alphanumeric
	if !regexp.MustCompile(`^[a-z-]+$`).MatchString(service.Name) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name should only be alphanumeric"})
		return
	}

	url, err := ws.nomad.CreatePage(service.Name, service.URL, service.Script)
	if err != nil {
		return
	}

	c.IndentedJSON(http.StatusCreated, PageResponse{
		URL: url,
	})
}

func (ws *WebService) Run() {
	router := gin.Default()
	router.POST("/services", ws.createPage)

	router.Run("localhost:3000")
}
