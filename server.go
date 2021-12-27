package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gitlab.com/JOsorio01/josorio-gin-app/controller"
	"gitlab.com/JOsorio01/josorio-gin-app/middlewares"
	"gitlab.com/JOsorio01/josorio-gin-app/repository"
	"gitlab.com/JOsorio01/josorio-gin-app/service"
)

var (
	videoRepository repository.VideoRepository = repository.NewVideoRepository()
	videoService    service.VideoService       = service.New(videoRepository)
	loginService    service.LoginService       = service.NewLoginService()
	jwtService      service.JWTService         = service.NewJWTService()

	videoController controller.VideoController = controller.New(videoService)
	loginController controller.LoginController = controller.NewLoginController(loginService, jwtService)
)

func main() {
	// Initialize server
	server := gin.New()
	// Statics
	server.Static("/css", "./templates/css")
	server.LoadHTMLGlob("./templates/*.html")
	// Middlewares
	server.Use(gin.Recovery(), gin.Logger())

	// Login: Authentication + Token creation
	server.POST("/login", func(c *gin.Context) {
		token := loginController.Login(c)
		if token != "" {
			c.JSON(http.StatusOK, gin.H{
				"token": token,
			})
		} else {
			c.JSON(http.StatusUnauthorized, nil)
		}
	})
	// JWT Authorization Middleware applies to "/api" only.
	apiRoutes := server.Group("/api", middlewares.AuthotizeJWT())
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
		apiRoutes.PUT("/videos/:id", func(c *gin.Context) {
			err := videoController.Update(c)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusCreated, gin.H{"message": "Video Updated!"})
			}
		})
		apiRoutes.DELETE("/videos/:id", func(c *gin.Context) {
			err := videoController.Delete(c)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusCreated, gin.H{"message": "Video Deleted!"})
			}
		})
	}
	// View Group
	viewRoutes := server.Group("/view")
	{
		viewRoutes.GET("/videos", videoController.ShowAll)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	server.Run(":" + port)

}
