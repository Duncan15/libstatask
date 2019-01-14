package nets

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/tatsushid/go-fastping"
	"net"
	"time"
)

//EasyPing ping the target ip addr
func EasyPing(ip string) bool {
	p := fastping.NewPinger()
	p.AddIP(ip)
	flagChan := make(chan bool, 2)
	fmt.Println("start to ping")
	p.OnRecv = func(addr *net.IPAddr, duration time.Duration) {
		glog.Infof("IP Address: %s receive, RTT %v", addr.String(), duration)
		flagChan <- true
	}
	p.OnIdle = func() {
		glog.Warningf("timeout for ping IP %s", ip)
		flagChan <- false
		close(flagChan)
	}
	p.Run()
	flag := <-flagChan
	for range flagChan {
	}
	return flag
}
