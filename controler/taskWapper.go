package controler

import (
	"github.com/gin-gonic/gin"
	"libstatask/work"
	"net/http"
)

func CollectSeatsActionNow(ctx *gin.Context) {
	run := work.CollectInfoFromCurrentSeats(work.CollectSeatsActionRule)
	go run()
	ctx.JSON(http.StatusOK, gin.H{
		"detail": "success to run",
	})
}
