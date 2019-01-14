package work

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"io/ioutil"
	"libstatask/common/dbs"
	"libstatask/common/nets"
	"libstatask/common/sys"
	"net/http"
	"strconv"
	"time"
)

func getInfoLink(roomID string, time time.Time) string {
	//roomID/2006-1-2/1234567890
	return fmt.Sprintf("http://10.12.162.31/ClientWeb/pro/ajax/device.aspx?room_id=%s&date=%s&act=get_rsv_sta&_=%d", roomID, time.Format("2006-01-02"), time.Unix()*1000)
}

//CollectInfoFromCurrentSeats collect info from current seats
func CollectInfoFromCurrentSeats(collectRule func(map[string]interface{})) func() {
	//return a function without parameters to fit the need of register a task to an executor
	return func() {
		if sys.IsRoot() && !nets.EasyPing("10.12.162.31") {
			glog.Warningf("can't find the target remote machine %v", "10.12.162.31")
			return
		}
		session := dbs.MySQL.NewSession()
		defer session.Close()
		var err error
		libUnitSlice := []*dbs.LibUnit{}
		if err = session.Asc("id").Where("deleted=?", dbs.Undeleted).Find(&libUnitSlice); err != nil {
			glog.Errorf("error happen when read from db, error: %v", err)
			return
		}
		for k := range libUnitSlice {
			//here no need to be concurrent, because it's possible to be blocked by the remote server
			infoLink := getInfoLink(libUnitSlice[k].RoomID, time.Now())
			collect(infoLink, collectRule)
		}
	}
}

//CollectSeatsAction collect seats action information
func collect(infoLink string, collectRule func(map[string]interface{})) {
	defer func() {
		if err := recover(); err != nil {
			glog.Warningf("uncaught error(%v)", err)
		}
	}()
	resp, err := http.Get(infoLink)
	if err != nil {
		glog.Warningf("error (%v) happens when invoking api %s", err, infoLink)
		return
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	ansMap := map[string]interface{}{}
	if err = json.Unmarshal(bodyBytes, &ansMap); err != nil {
		glog.Warningf("error (%v) happens when unmarshaling the response of api %s", err, infoLink)
		return
	}
	arr := ansMap["data"].([]interface{})
	for k := range arr {
		seat := arr[k].(map[string]interface{})
		collectRule(seat)
	}

}

//CollectSeatsInfoRule the rule to collect seat information
func CollectSeatsInfoRule(input map[string]interface{}) {
	session := dbs.MySQL.NewSession()
	defer session.Close()
	session.Begin()
	defer func() {
		if err := recover(); err != nil {
			glog.Warningf("uncaught error(%v)", err)
			session.Rollback()
			return
		}
		session.Commit()
	}()

	seatName := input["name"].(string)
	location := input["labName"].(string)
	roomID := strconv.FormatInt(int64(input["roomId"].(float64)), 10)

	seat := &dbs.Seat{
		SeatName: seatName,
		Location: location,
		RoomID:   roomID,
	}

	if has, err := session.Cols("seat_name", "location", "room_id").Exist(seat); err != nil {
		panic(err)
	} else if !has {
		session.InsertOne(seat)
	}
}

//CollectSeatsActionRule the rule to collect user_seat action
func CollectSeatsActionRule(input map[string]interface{}) {
	session := dbs.MySQL.NewSession()
	defer session.Close()
	session.Begin()
	defer func() {
		if err := recover(); err != nil {
			glog.Warningf("uncaught error(%v)", err)
			session.Rollback()
			return
		}
		session.Commit()
	}()

	//seatTsArr is the operations have been acted in this seat
	seatTsArr := input["ts"].([]interface{})
	if len(seatTsArr) == 0 { //now this is no operations in this seat
		return
	}
	for _, seatTsInterface := range seatTsArr {
		seatTs := seatTsInterface.(map[string]interface{})
		userIDStr := seatTs["accno"].(string)
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			panic(err)
		}
		userName := seatTs["owner"].(string)
		startTimeStr := seatTs["start"].(string)
		startTime, err := time.ParseInLocation("2006-01-02 15:04", startTimeStr, time.Local)
		if err != nil {
			panic(err)
		}
		seatName := input["name"].(string)

		if has, err := session.Where("user_id=?", userID).Exist(new(dbs.User)); err != nil {
			panic(err)
		} else if !has {
			session.InsertOne(&dbs.User{
				UserID:   userID,
				UserName: userName,
			})
		}

		action := new(dbs.UserSeat)
		action.UserID = userID
		action.StartTime = startTime
		action.SeatName = seatName
		if has, err := session.Get(action); err != nil {
			panic(err)
		} else if has { //if this record has existed
			action.UpdateTime = time.Now()
			if _, err := session.ID(action.ID).Update(action); err != nil {
				panic(err)
			}
		} else { //if no exist
			action.UpdateTime = time.Now()
			if _, err := session.InsertOne(action); err != nil {
				panic(err)
			}
		}
	}

}
