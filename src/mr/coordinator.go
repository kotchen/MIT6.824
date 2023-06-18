package mr

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strings"
)

type TaskContext struct {
	fileName string
	status   bool
}

type Coordinator struct {
	// Your definitions here.
	TaskList []TaskContext
}

// Your code here -- RPC handlers for the worker to call.

// 统计一个文件中的单词个数，返回该文件中每个出现的单词的个数
func MapTask(args *MapTaskArgs, reply *MapTaskReply) {
	reply.success = false
	fileName := args.fileName
	// open a file
	file, err := os.Open(fileName)
	defer func() {
		err := file.Close()
		if err != nil {
			// Handle the error
			fmt.Println("Failed to close the file:", err)
		}
	}()

	if err != nil {
		log.Fatalf("cannot open %v", fileName)
	}

	reader := bufio.NewReader(file)
	//循环读取文件的内容
	for {
		str, err := reader.ReadString('\n') //读到一个换行就结束
		if err == io.EOF {                  //io.EOF表示文件的末尾
			break
		}
		words := strings.Split(str, " ")
		for _, word := range words {
			if word != "" {
				reply.wordCountMap[word]++
			}
		}
	}
	reply.success = true
}

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

func (c *Coordinator) dispatchTask(args *RequestTaskArgs, reply *RequestTaskReply) {
	for _, task := range c.TaskList {
		if !task.status {
			reply.taskType = Map
			reply.fileName = task.fileName
			break
		}
	}
}

// start a thread that listens for RPCs from worker.go
func (c *Coordinator) server() {
	rpc.Register(c)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := coordinatorSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
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
	for _, task := range c.TaskList {
		if !task.status {
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
	c := Coordinator{}

	// Your code here.
	// 初始化任务列表
	for _, fileName := range files {
		fmt.Println("fileName:", fileName)
		c.TaskList = append(c.TaskList, TaskContext{fileName, false})
	}

	c.server()
	return &c
}
