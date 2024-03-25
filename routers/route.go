package routers

import (
	"final-project-go/controllers/commentcontroller"
	"final-project-go/controllers/photocontroller"
	"final-project-go/controllers/socialmediacontroller"
	"final-project-go/controllers/usercontroller"
	"final-project-go/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// User group
	users := router.Group("/users")
	{
		users.POST("/register", usercontroller.Register)
		users.POST("/login", usercontroller.Login)
		users.PUT("/", usercontroller.Update)
		users.DELETE("/", usercontroller.Delete)
	}

	// Photos group
	photos := router.Group("/photos")
	photos.Use(middlewares.AuthMiddleware())
	{
		photos.POST("/", photocontroller.CreatePhoto)
		photos.GET("/", photocontroller.GetAll)
		photos.GET("/:id", photocontroller.GetOne)
		photos.PUT("/:id", photocontroller.UpdatePhoto)
		photos.DELETE("/:id", photocontroller.DeletePhoto)
	}

	// Comments group
	//comments := router.Group("/comments")
	//{
	//	comments.POST("/", commentcontroller.CreateComment)
	//	comments.GET("/", commentcontroller.GetAll)
	//	comments.GET("/:id", commentcontroller.GetOne)
	//	comments.PUT("/:id", commentcontroller.UpdateComment)
	//	comments.DELETE("/:id", commentcontroller.DeleteComment)
	//}

	comments := router.Group("/comments")
	comments.Use(middlewares.AuthMiddleware()) // Terapkan middleware di sini
	{
		comments.POST("/", commentcontroller.CreateComment)
		comments.GET("/", commentcontroller.GetAll)
		comments.GET("/:id", commentcontroller.GetOne)
		comments.PUT("/:id", commentcontroller.UpdateComment)
		comments.DELETE("/:id", commentcontroller.DeleteComment)
	}

	socialmedias := router.Group("/socialmedias")
	socialmedias.Use(middlewares.AuthMiddleware())
	{
		socialmedias.POST("/", socialmediacontroller.CreateSocialMedia)
		socialmedias.GET("/", socialmediacontroller.GetAll)
		socialmedias.GET("/:id", socialmediacontroller.GetOne)
		socialmedias.PUT("/:id", socialmediacontroller.UpdateSocialMedia)
		socialmedias.DELETE("/:id", socialmediacontroller.DeleteSocialMedia)
	}

	return router
}
