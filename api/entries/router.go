package entries

import "github.com/gin-gonic/gin"

func AddEntryRoutes(r *gin.Engine) {
	r.GET("/entries", GetAllEntries)
	r.GET("/entries/me", GetMyEntries)
	r.POST("/entries", CreateEntry)
	r.DELETE("/entries/:id", DeleteEntry)
	r.PUT("/entries/:id", EditEntry)
}
