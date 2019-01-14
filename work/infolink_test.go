package work

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"libstatask/common/dbs"
	"testing"
	"time"
)

func init() {
	flag.Set("logtostderr", fmt.Sprintf("%t", true))
	dbs.NewMySQL("127.0.0.1:3306", "root", "12345678", "libstatask")
}
func TestBuildInfoLink(t *testing.T) {
	fmt.Println(getInfoLink("100496304", time.Now()))
}

func TestCollectSeatsInfo(t *testing.T) {
	CollectInfoFromCurrentSeats(CollectSeatsInfoRule)()
}
func TestCollectSeatsAction(t *testing.T) {
	CollectInfoFromCurrentSeats(CollectSeatsActionRule)()
}

func TestCollectInfoFromSpecifiedDomain(t *testing.T) {
	glog.Info("start test")
	collect(getInfoLink("100496278", time.Now()), CollectSeatsActionRule)
	glog.Info("end test")
}
