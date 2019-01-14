package sys

import (
	"github.com/golang/glog"
	"os/user"
	"strings"
)

//IsRoot judge whether the current user is root or not
func IsRoot() bool {
	if usr, err := user.Current(); err != nil {
		glog.Errorf("error happen when find the current user, %v", err)
		return false
	} else {
		if strings.ToUpper(usr.Name) == "ROOT" {
			return true
		}
		return false
	}
}
