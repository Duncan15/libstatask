package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"net/http"
)

func CusRecovery() gin.HandlerFunc {
	return func(context *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				glog.Errorf("panic error, cause %v", err)
				context.JSON(http.StatusOK, gin.H{
					"errmsg": "panic",
					"detail": err,
				})
			}
		}()
		context.Next()
	}
}
