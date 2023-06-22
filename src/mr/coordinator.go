package mr

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

type Coordinator struct {
	// Your definitions here.
	TaskMap map[string]bool
}

// Your code here -- RPC handlers for the worker to call.

func ReduceTask(args *ReduceTaskArgs, reply *ReduceTaskReply) {
	wordCountMap := args.wordCountMap
	reply.wordName = args.targetWordName

	for word, count := range wordCountMap {
		if word == args.targetWordName {
			reply.wordCount += count
		}
	}
}

// an example RPC handler.
//
// the RPC argument and reply types are defined in rpc.go.
func (c *Coordinator) Example(args *ExampleArgs, reply *ExampleReply) error {
	reply.Y = args.X + 1
	return nil
}

func (c *Coordinator) RequestTask(args *RequestTaskArgs, reply *RequestTaskReply) {
	c.dispatchTask(args, reply)
}

func (c *Coordinator) AllTaskDone() {

}

func (c *Coordinator) OnMapFinish(args *MapTaskArgs, reply *MapTaskReply) {
	if reply.success {
		c.TaskMap[reply.fileName] = true
	}
}

func (c *Coordinator) dispatchTask(args *RequestTaskArgs, reply *RequestTaskReply) {
	for fileName, status := range c.TaskMap {
		if !status {
			reply.taskType = Map
			reply.fileName = fileName
			break
		}
	}
}

// start a thread that listens for RPCs from worker.go
func (c *Coordinator) server() {
	rpc.Register(c)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	//sockname := coordinatorSock()
	//os.Remove(sockname)
	//l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

// main/mrcoordinator.go calls Done() periodically to find out
// if the entire job has finished.
func (c *Coordinator) Done() bool {
	ret := false

	// Your code here.
	for _, status := range c.TaskMap {
		if !status {
			return false
		}
	}
	ret = true
	return ret
}

// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
func MakeCoordinator(files []string, nReduce int) *Coordinator {
	c := Coordinator{
		TaskMap: make(map[string]bool),
	}

	// Your code here.
	// 初始化任务列表
	for _, fileName := range files {
		fmt.Println("fileName:", fileName)
		_, exists := c.TaskMap[fileName]
		if !exists {
			c.TaskMap[fileName] = false
		}
	}

	c.server()
	return &c
}
