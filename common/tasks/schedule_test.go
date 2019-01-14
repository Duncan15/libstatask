package tasks

import (
	"flag"
	"fmt"
	"testing"
	"time"
)
func init() {
	flag.Set("logtostderr", fmt.Sprintf("%t", true))
}

func TestTimingScheduler(t *testing.T) {
	scheduler := NewTimingScheduler()
	task := NewTimingTask("testTask", func() {
		fmt.Printf("print in testTask, now is %s\n", time.Now().String())
	}, time.Now().Unix() + 5, 10 * time.Second)
	scheduler.RegisterTask(task)
	scheduler.Run()
	time.Sleep(30 * time.Second)
	scheduler.Close()
	time.Sleep(30 * time.Second)
}
