package controler

import "github.com/gin-gonic/gin"

var (
	GetURLMap  = map[string]gin.HandlerFunc{}
	PostURLMap = map[string]gin.HandlerFunc{}
)

func init() {
	GetURLMap["/seatsAction"] = FindSeatsActionByName
	GetURLMap["/users"] = FindUsersBySeatName
	GetURLMap["/dbInfo"] = DataBaseInfo
	GetURLMap["/collectInfoNow"] = CollectSeatsActionNow
	GetURLMap["/frequentNeighbours"] = FindFrequentNeighboursByName
}
