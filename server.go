package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gitlab.com/JOsorio01/josorio-gin-app/controller"
	"gitlab.com/JOsorio01/josorio-gin-app/middlewares"
	"gitlab.com/JOsorio01/josorio-gin-app/service"
)

var (
	videoService    service.VideoService       = service.New()
	videoController controller.VideoController = controller.New(videoService)
)

func main() {
	// Initialize server
	server := gin.New()
	// Statics
	server.Static("/css", "./templates/css")
	server.LoadHTMLGlob("./templates/*.html")
	// Middlewares
	server.Use(gin.Recovery(), gin.Logger())
	// Api Group
	apiRoutes := server.Group("/api", middlewares.BasicAuth())
	{
		apiRoutes.GET("/videos", func(c *gin.Context) {
			c.JSON(http.StatusOK, videoController.FindAll())
		})
		apiRoutes.POST("/videos", func(c *gin.Context) {
			err := videoController.Save(c)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusCreated, gin.H{"message": "Video Created!"})
			}
		})
	}
	// View Group
	viewRoutes := server.Group("/view")
	{
		viewRoutes.GET("/videos", videoController.ShowAll)
	}
	port := os.Getenv("PORT")
	server.Run(":" + port)
}
