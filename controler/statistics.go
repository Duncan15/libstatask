package controler

import (
	"github.com/gin-gonic/gin"
	"libstatask/common/dbs"
	"net/http"
)

func DataBaseInfo(ctx *gin.Context) {
	session := dbs.MySQL.NewSession()
	defer session.Close()
	data := map[string]interface{}{}
	if num, err := session.Where("deleted=?", dbs.Undeleted).Count(new(dbs.Seat)); err != nil {
		panic(err)
	} else {
		data["seatsNum"] = num
	}
	if num, err := session.Where("deleted=?", dbs.Undeleted).Count(new(dbs.User)); err != nil {
		panic(err)
	} else {
		data["userNum"] = num
	}
	if num, err := session.Where("deleted=?", dbs.Undeleted).Count(new(dbs.UserSeat)); err != nil {
		panic(err)
	} else {
		data["bookingActionsNum"] = num
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": data,
	})

}
