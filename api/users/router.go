package users

import "github.com/gin-gonic/gin"

func AddUserRoutes(r *gin.Engine) {
	r.GET("/users", GetAllUsers)
	r.GET("/users/me", GetMe)
	r.POST("/users", CreateUser)
	r.GET("/users/:id", GetUserById)
	r.DELETE("/users/:id", DeleteUser)
	r.PUT("/users/:id", EditUser)
}
