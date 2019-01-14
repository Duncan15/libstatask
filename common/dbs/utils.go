package dbs

import "github.com/golang/glog"

func GetUsersByInIDs(ids ...interface{}) map[int64]*User {
	session := MySQL.NewSession()
	defer session.Close()
	ans := map[int64]*User{}
	usrs := []User{}
	if err := session.In("user_id", ids...).Find(&usrs); err != nil {
		glog.Warningf("GetUsersByInIDs fail, cause %v", err)
		return ans
	}
	for k := range usrs {
		ans[usrs[k].UserID] = &usrs[k]
	}
	return ans
}
