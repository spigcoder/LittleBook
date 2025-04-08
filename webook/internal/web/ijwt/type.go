package ijwt

import (
	"github.com/gin-gonic/gin"
)

type Server interface {
	LoginJwt(c *gin.Context)
}
