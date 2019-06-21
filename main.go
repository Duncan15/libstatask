package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"gopkg.in/gcfg.v1"
	"libstatask/cfgs"
	"libstatask/common/dbs"
	"libstatask/common/middlewares"
	"libstatask/common/tasks"
	"libstatask/controler"
	"libstatask/work"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	preRun()
	defer glog.Flush()
	defer dbs.CloseMySQL()
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
	glog.Info("start the library statistic service")
	scheduler := tasks.NewTimingScheduler()

	libLinkTask := tasks.NewTimingTask("collectLibDomainLink", work.CollectLibDomainLink, time.Now().Unix(), 24*time.Hour)
	seatsInfoTask := tasks.NewTimingTask("collectSeatsInfo", work.CollectInfoFromCurrentSeats(work.CollectSeatsInfoRule), time.Now().Unix(), 24*time.Hour)
	seatsActionTask := tasks.NewTimingTask("collectSeatsAction", work.CollectInfoFromCurrentSeats(work.CollectSeatsActionRule), time.Now().Unix(), 10*time.Minute)
	scheduler.RegisterTask(libLinkTask)
	scheduler.RegisterTask(seatsInfoTask)
	scheduler.RegisterTask(seatsActionTask)

	r := gin.New()
	r.Use(middlewares.CusRecovery())
	for k, v := range controler.GetURLMap {
		r.GET(k, v)
	}
	for k, v := range controler.PostURLMap {
		r.POST(k, v)
	}
	go r.Run(":3333")

	scheduler.Run()
	select {
	case <-sigc:
		scheduler.Close()
		glog.Info("stop the library statistic service")
	}
}

func preRun() {
	flag.Parse()
	if *cfgs.MODE == "online" {
		if err := gcfg.ReadFileInto(cfgs.Conf, "./config_online.ini"); err != nil {
			log.Fatalln(err)
		}

		//the following configuration must be set when before any glog's output method is invoked
		//because the log_dir would be used when first output and only be used once time
		flag.Set("log_dir", cfgs.Conf.Log.LogAddress) //Log files will be written to this directory instead of the default temporary directory
		//flag.Set("logtostderr", "false")               //Logs are written to standard error instead of to files
		//flag.Set("alsologtostderr", "false")           //Logs are written to standard error as well as to files
		//flag.Set("stderrthreshold", "ERROR")           //Log events at or above this severity are logged to standard error as well as to files
	} else if *cfgs.MODE == "local" {
		if err := gcfg.ReadFileInto(cfgs.Conf, "./config_local.ini"); err != nil {
			log.Fatalln(err)
		}
		flag.Set("logtostderr", "true") //Logs are written to standard error instead of to files
	} else {
		log.Fatalln("run in an unknown mode, exit")
	}
	dbs.NewMySQL(cfgs.Conf.MySQL.TcpAddress, cfgs.Conf.MySQL.UserName, cfgs.Conf.MySQL.Password, cfgs.Conf.MySQL.DbName)
	dbs.UseFileLogger(cfgs.Conf.MySQL.LogAddress)
}
