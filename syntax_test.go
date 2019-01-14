package main

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"strconv"
	"testing"
)

func init() {
	flag.Set("logtostderr", fmt.Sprintf("%t", true))
}
func TestSyntax(t *testing.T) {
	flagChan := make(chan bool)
	go func() {
		flagChan <- true
	}()
	go func() {
		flagChan <- false
	}()
	select {
	case flag := <-flagChan:
		fmt.Println(flag)

	}
}

func TestGlog(t *testing.T) {
	glog.Info("hhh")
}
func TestString(t *testing.T) {
	fmt.Println("5C251"[len("5C251")-3:])
}

func TestParseInt(t *testing.T) {
	v, _ := strconv.ParseInt("056", 10, 64)
	fmt.Println(v)
}
