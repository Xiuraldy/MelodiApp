package userConfig

import (
	"github.com/gin-gonic/gin"
)

func AddUserConfigRoutes(r *gin.Engine) {
	r.PUT("/user-config", updateUserConfig)
	r.GET("/user-config", getUserConfig)
}
