package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/arabian9ts/geekfes/code"
)

var (
	statusCode = map[code.ErrCode]int{
		code.OK:              http.StatusOK,
		code.InvalidArgument: http.StatusBadRequest,
		code.NotFound:        http.StatusNotFound,
		code.Internal:        http.StatusInternalServerError,
	}
)

func apiError(c *gin.Context, err error) {
	status := statusCode[code.From(err)]
	c.AbortWithError(status, err)
}

func badRequest(c *gin.Context, err error) {
	c.AbortWithError(http.StatusBadRequest, err)
}

func write(c *gin.Context, data interface{}) {
	c.Header("Cache-Control", "s-maxage=1,max-age=1")
	c.JSON(http.StatusOK, data)
}
