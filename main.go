package main

import (
	"context"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginAdapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginAdapter.GinLambda

// album represents data about a record album.
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var albumMap = map[string]album{
	"1": {ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	"2": {ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	"3": {ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func main() {
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.POST("/albums/add/:id", postAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.DELETE("albums/delete/:id", deleteAlbumByID)

	environment := os.Getenv("GIN_MODE")
	if environment == "release" {
		ginLambda = ginAdapter.New(router)
		lambda.Start(Handler)
	} else {
		router.Run("localhost:8080")
	}
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, request)
}

// getAlbums responds with the list of all albums
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albumMap)
}

func postAlbums(c *gin.Context) {
	var newAlbum album
	id := c.Param("id")
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	albumMap[id] = newAlbum
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	elem, ok := albumMap[id]
	if ok {
		c.IndentedJSON(http.StatusOK, elem)
		return
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}
func deleteAlbumByID(c *gin.Context) {
	id := c.Param("id")
	elem, ok := albumMap[id]
	if ok {
		delete(albumMap, id)
		c.IndentedJSON(http.StatusOK, elem)
		return
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}
