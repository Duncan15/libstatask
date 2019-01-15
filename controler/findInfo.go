package controler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"libstatask/common/dates"
	"libstatask/common/dbs"
	"libstatask/common/retries"
	"libstatask/common/sorts"
	"net/http"
	"sort"
	"strconv"
	"time"
)

//input name:string
//		days:int64
func FindSeatsActionByName(ctx *gin.Context) {
	session := dbs.MySQL.NewSession()
	defer session.Close()
	content := []map[string]interface{}{}
	userName := ctx.Query("name")
	if userName == "" {
		panic("no name parameter")
	}
	daysStr := ctx.DefaultQuery("days", "1")
	days, _ := strconv.ParseInt(daysStr, 10, 64)
	user := new(dbs.User)
	if has, err := session.Where("user_name=?", userName).Get(user); err != nil {
		panic(err)
	} else if !has {
		ctx.JSON(http.StatusOK, gin.H{
			"errmsg": "user can't find",
		})
		return
	}
	userSeats := make([]dbs.UserSeat, 0)
	if err := session.Where("user_id=?", user.UserID).And("start_time>?", dates.GetStartTimeOfCurrentDay(time.Now()).Add(time.Duration((-days+1)*24)*time.Hour).String()).Asc("id").Find(&userSeats); err != nil {
		panic(err)
	}
	for k := range userSeats {
		unit := map[string]interface{}{
			"startTime":  userSeats[k].StartTime.String(),
			"seatName":   userSeats[k].SeatName,
			"updateTime": userSeats[k].UpdateTime,
		}
		content = append(content, unit)
	}
	data := map[string]interface{}{
		"content": content,
		"total":   len(content),
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": data,
	})

}

//input: seatName: string
//       days: int64
func FindUsersBySeatName(ctx *gin.Context) {
	session := dbs.MySQL.NewSession()
	defer session.Close()
	content := []map[string]interface{}{}
	seatName := ctx.Query("seatName")
	if seatName == "" {
		panic("no seatName parameter")
	}
	daysStr := ctx.DefaultQuery("days", "1")
	days, _ := strconv.ParseInt(daysStr, 10, 64)
	userSeats := []dbs.UserSeat{}
	if err := session.Where("seat_name=?", seatName).And("start_time>?", dates.GetStartTimeOfCurrentDay(time.Now()).Add(time.Duration((-days+1)*24)*time.Hour).String()).Asc("id").Find(&userSeats); err != nil {
		panic(err)
	}
	userIDs := []interface{}{}
	for k := range userSeats {
		userIDs = append(userIDs, userSeats[k].UserID)
	}
	users := []dbs.User{}
	if err := session.In("user_id", userIDs...).Find(&users); err != nil {
		panic(err)
	}
	IDNameMap := map[int64]string{}
	for k := range users {
		IDNameMap[users[k].UserID] = users[k].UserName
	}
	for k := range userSeats {
		content = append(content, map[string]interface{}{
			"startTime":  userSeats[k].StartTime,
			"updateTime": userSeats[k].UpdateTime,
			"user":       IDNameMap[userSeats[k].UserID],
		})
	}
	data := map[string]interface{}{}
	data["content"] = content
	data["total"] = len(content)
	ctx.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

//FindFrequentNeighboursByName
//input:	name:string
//			days:int64
func FindFrequentNeighboursByName(ctx *gin.Context) {
	content := []interface{}{}
	session := dbs.MySQL.NewSession()
	defer session.Close()
	name := ctx.Query("name")
	if name == "" {
		panic("no name parameter")
	}
	daysStr := ctx.DefaultQuery("days", "1")
	days, _ := strconv.ParseInt(daysStr, 10, 64)
	usr := new(dbs.User)

	var has bool
	if err := retries.Retry(5, func() error {
		var innerErr error
		has, innerErr = session.Where("user_name=?", name).Get(usr)
		return innerErr
	}); err != nil {
		panic(fmt.Errorf("can't get user by user_name, cause %v", err))
	} else if !has {
		ctx.JSON(http.StatusOK, gin.H{
			"detail": "no user find",
		})
		return
	}

	userSeats := make([]dbs.UserSeat, 0)

	if err := retries.Retry(5, func() error {
		return session.Where("user_id=?", usr.UserID).And("start_time>?", dates.GetStartTimeOfCurrentDay(time.Now()).Add(time.Duration((-days+1)*24)*time.Hour).String()).Asc("id").Find(&userSeats)
	}); err != nil {
		panic(fmt.Errorf("can't get actions by user_id and start_time, cause %v", err))
	}

	collector := make(chan []dbs.UserSeat, len(userSeats))
	for k := range userSeats {
		go findNeighbourAction(&userSeats[k], collector)
	}
	countMap := map[int64]int64{}
	for i := 0; i < len(userSeats); i++ {
		subActions := <-collector
		for k := range subActions {
			if v, has := countMap[subActions[k].UserID]; has {
				countMap[subActions[k].UserID] = v + 1
			} else {
				countMap[subActions[k].UserID] = 1
			}
		}
	}
	//glog.Infof("user{ID:%d,days:%d} 's neighbour action number is %d", usr.UserID, days, len(countMap))
	s := sorts.Int64KeyInterface{}
	for k, v := range countMap {
		o := sorts.Int64KeyStruct{
			Key:   v,
			Value: k,
		}
		s = append(s, o)
	}
	sort.Sort(sort.Reverse(s))
	l := s.Len()

	//return the value whose count is bigger than 2
	ids := []interface{}{}
	for i := range s {
		if s[i].Key < 2 {
			break
		}
		ids = append(ids, s[i].Value)
	}
	if l > len(ids) {
		l = len(ids)
	}

	IDUsrsMap := dbs.GetUsersByInIDs(ids...)
	for i := 0; i < l; i++ {
		content = append(content, map[string]interface{}{
			"name":  IDUsrsMap[s[i].Value.(int64)].UserName,
			"id":    IDUsrsMap[s[i].Value.(int64)].UserID,
			"count": s[i].Key,
		})
	}
	data := map[string]interface{}{
		"neighbours": map[string]interface{}{
			"content": content,
			"total":   len(content),
		},
		"actions": map[string]interface{}{
			"content": userSeats,
			"total":   len(userSeats),
		},
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": data,
	})

}

//findNeighbourAction neighbours action is action that happen close to current action in time and space
func findNeighbourAction(acton *dbs.UserSeat, ans chan<- []dbs.UserSeat) {
	defer func() {
		if err := recover(); err != nil {
			glog.Errorf("error when findNeighbourAction, cause %v", err)
			ans <- []dbs.UserSeat{}
		}
	}()
	if acton == nil {
		panic("action parameter can't be nil")
	}
	session := dbs.MySQL.NewSession()
	defer session.Close()
	seatEndIDStr := acton.SeatName[len(acton.SeatName)-3:]
	seatStartIDStr := acton.SeatName[0 : len(acton.SeatName)-3]
	seatEndID, err := strconv.ParseInt(seatEndIDStr, 10, 64)
	if err != nil {
		panic(fmt.Errorf("parse seatName error, casue %v", err))
	}
	seatEndIDs := []int64{}
	//get the front and the back 5 id
	for i := int64(1); i <= 5; i++ {
		up := seatEndID + i
		down := seatEndID - i
		seatEndIDs = append(seatEndIDs, up, down)
	}
	seatNames := []interface{}{}
	for _, v := range seatEndIDs {
		//filter the invalid end id
		if v <= 0 {
			continue
		}

		//build the valid seatName
		seatName := ""
		if v >= 100 {
			seatName = seatStartIDStr + strconv.FormatInt(v, 10)
		} else {
			seatName = seatStartIDStr + "0" + strconv.FormatInt(v, 10)
		}
		seatNames = append(seatNames, seatName)
	}

	//if the update time of a action is bigger than current action's start time
	//and the start time of the action is smaller than current action's update time
	//we think the action is close to current action
	neighourActions := []dbs.UserSeat{}

	if err := session.In("seat_name", seatNames...).Where("start_time<=?", acton.UpdateTime.String()).And("update_time>=?", acton.StartTime.String()).Find(&neighourActions); err != nil {
		panic(err)
	}
	//glog.Infof("acton{ID:%d}'s neighbour action number is %d", acton.ID, len(neighourActions))
	ans <- neighourActions
}
