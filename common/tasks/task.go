package tasks

import (
	"github.com/golang/glog"
	"time"
)

//Task
type Task struct {
	taskName string
	executable func()
}

//NewSimpleTask create a new task
func NewSimpleTask(taskName string, executable func()) *Task {
	return &Task{
		taskName:taskName,
		executable:executable,
	}
}

func (task *Task)run()  {
	defer func() {
		if err := recover() ; err != nil {
			glog.Errorf("task %s fail to execute, the reason is %v", task.taskName, err)
		}
	}()
	glog.Infof("start to run task %s", task.taskName)
	task.executable()
	glog.Infof("finish to run task %s", task.taskName)
}

type TimingTask struct {
	Task
	duration time.Duration
	runTime int64
	left int//the left time to run
}

//NewTimingTask create a new timing task, set the default left run time to infinity
func NewTimingTask(taskName string, executable func(), runTime int64, duration time.Duration) *TimingTask {
	return &TimingTask{
		Task: Task{
			taskName: taskName,
			executable:executable,
		},
		runTime: runTime,
		duration: duration,
		left: -1,
	}
}
//SetLeft set left run time to the specified number
func (task *TimingTask)SetLeft(left int)  {
	task.left = left
}