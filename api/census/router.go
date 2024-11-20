package census

import "github.com/gin-gonic/gin"

func AddCensusRoutes(r *gin.Engine) {
	r.GET("/census", GetCensus)
}
