package tasks

import (
	"github.com/golang/glog"
	"time"
)

type TaskScheduler interface {
	Run()
}

//QueueScheduler scheduler to run task sequentially
type QueueScheduler struct {
	taskQueue []*Task
}

//NewQueueScheduler initialize a task scheduler
func NewQueueScheduler() *QueueScheduler {
	return &QueueScheduler{}
}

//RegisterTask register a new task
func (scheduler *QueueScheduler) RegisterTask(task *Task) {
	scheduler.taskQueue = append(scheduler.taskQueue, task)
}

//Run sync method
func (scheduler *QueueScheduler) Run() {
	for _, task := range scheduler.taskQueue {
		task.run()
	}
}

//TimingScheduler scheduler to run task at fixed time
type TimingScheduler struct {
	stopChan   chan bool
	taskChan   chan *TimingTask
	ticker     <-chan time.Time
	curPos     int64
	taskMatrix [][]*TimingTask
}

//NewTimingScheduler initialize a timing scheduler
func NewTimingScheduler() *TimingScheduler {
	taskMaxtrix := make([][]*TimingTask, 3600)

	return &TimingScheduler{
		stopChan:   make(chan bool),
		taskChan:   make(chan *TimingTask, 10*10),
		ticker:     time.Tick(time.Second),
		curPos:     time.Now().Unix() % 3600,
		taskMatrix: taskMaxtrix,
	}
}

func (scheduler *TimingScheduler) RegisterTask(task *TimingTask) {
	pos := task.runTime % 3600
	scheduler.taskMatrix[pos] = append(scheduler.taskMatrix[pos], task)
}

//Run async method
func (scheduler *TimingScheduler) Run() {
	glog.Info("start to run the timing scheduler")
	go func() {
		//start 10 goroutine to run all the task until this scheduler is stoped
		for i := 0; i < 10; i++ {
			go func() {
				for task := range scheduler.taskChan {
					task.run()
				}
				glog.Info("exit the worker goroutine")
				<-scheduler.stopChan
			}()
		}

		scheduler.curPos = time.Now().Unix() % 3600

		//loop here until the executor is stopped
	PRODUCE_LOOP:
		for {
			select {
			case <-scheduler.ticker:
				curTimeUnix := time.Now().Unix()
				taskSlice := scheduler.taskMatrix[scheduler.curPos]      //temperately store the current position's taskSlice
				scheduler.taskMatrix[scheduler.curPos] = []*TimingTask{} //clear the current postion's taskSlice
				if len(taskSlice) != 0 {                                 //if the size of task slice is zero, don't care about it

					for i := range taskSlice {

						//if left time is bigger than zero, this task have limited run time
						//if left time is smaller than zero, this task have infinity run time
						//if left time is equal to zero, this task run finish
						if taskSlice[i].left != 0 {

							//because when go to the current position,
							//the currentTimeUnix must be bigger than or equal to task.runTime
							//if the curTimeUnix is smaller than runTime, the situation is that the task shouldn't run at this round
							if taskSlice[i].runTime <= curTimeUnix {
								scheduler.taskChan <- taskSlice[i]
								taskSlice[i].runTime = curTimeUnix + int64(taskSlice[i].duration.Seconds())
								pos := taskSlice[i].runTime % 3600
								scheduler.taskMatrix[pos] = append(scheduler.taskMatrix[pos], taskSlice[i])

								//if the left time bigger than zero, decrease it
								if taskSlice[i].left > 0 {
									taskSlice[i].left--
								}
							} else { //if the current task no run, just store it back to the current position's taskSlice
								scheduler.taskMatrix[scheduler.curPos] = append(scheduler.taskMatrix[scheduler.curPos], taskSlice[i])
							}

						}
					}

				}
				//increase the position to next
				scheduler.curPos++
				scheduler.curPos %= 3600
			case <-scheduler.stopChan:
				close(scheduler.taskChan)
				break PRODUCE_LOOP

			}
		}

		<-scheduler.stopChan

	}()
}

func (scheduler *TimingScheduler) Close() {
	glog.Info("exit the timing scheduler")
	//inform the produce_loop to exit
	scheduler.stopChan <- true

	//wait the ten worker to exit
	for i := 0; i < 10; i++ {
		scheduler.stopChan <- true
	}

	//block here to wait the produce_loop exit
	scheduler.stopChan <- true

}
