package initialize

import "short-code/task"

var (
	flushBindDataTask = task.FlushBindDataTask{}
)

func Tasks() {
	flushBindDataTask.Run()
}
