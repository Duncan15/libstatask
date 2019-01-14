package dbs

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"log"
	"os"
	"time"
)

var (
	//MySQL resource
	MySQL *xorm.Engine
	stop  bool
)

//NewMySQL initialize mysql
func NewMySQL(tcpAddr string, userName string, password string, dbName string) *xorm.Engine {
	url := userName + ":" + password + "@(" + tcpAddr + ")/" + dbName + "?charset=utf8"
	var err error //if declare err in the next line, MySQL would become a local variable in this method
	MySQL, err = xorm.NewEngine("mysql", url)
	if err != nil {
		panic(fmt.Sprintf("mysql initialize error: %v", err))
	}
	go keepAlive(tcpAddr, userName, password, dbName)
	return MySQL
}

//UseLocalMode invoke when run at local mode
func UseLocalMode() {
	if MySQL == nil {
		log.Fatalln("MySQL haven't been initialized")
	}
	MySQL.Logger().SetLevel(core.LOG_DEBUG)
	MySQL.ShowSQL(true)
}

//UseFileLogger invoke shen run at online mode
func UseFileLogger(fileName string) {
	f, err := os.Create(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	MySQL.SetLogger(xorm.NewSimpleLogger(f))
	MySQL.Logger().SetLevel(core.LOG_WARNING)
	MySQL.ShowSQL(true)

}

//CloseMySQL close mysql connection pool
func CloseMySQL() {
	stop = true
}
func keepAlive(tcpAddr string, userName string, password string, dbName string) {
KEEP_ALIVE:
	for !stop {
		time.Sleep(10 * time.Second)
		for i := 0; i < 5; i++ {
			if MySQL.Ping() == nil {
				continue KEEP_ALIVE
			}
		}
		NewMySQL(tcpAddr, userName, password, dbName)
	}
	MySQL.Close()
}
