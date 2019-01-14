package work

import (
	"libstatask/common/dbs"
	"testing"
)

func init() {
	dbs.NewMySQL("127.0.0.1:3306", "root", "12345678", "libstatask")
}
func TestCollectLibDomainLink(t *testing.T) {
	CollectLibDomainLink()
}
