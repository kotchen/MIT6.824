package mr

//
// RPC definitions.
//
// remember to capitalize all names.
//

import (
	"os"
	"strconv"
)

//
// example to show how to declare the arguments
// and reply for an RPC.
//

type TaskType int

const (
	Map    TaskType = iota // 0
	Reduce                 // 1
)

type ExampleArgs struct {
	X int
}

type ExampleReply struct {
	Y int
}

// Add your RPC definitions here.

type RequestTaskArgs struct {
}

type RequestTaskReply struct {
	taskType TaskType
	fileName string
}

type MapTaskArgs struct {
	fileName string
}

type MapTaskReply struct {
	success   bool
	fileName  string
	wordCount []KeyValue
}

type ReduceTaskArgs struct {
	wordCountMap   map[string]int
	targetWordName string
}

type ReduceTaskReply struct {
	success   bool
	wordName  string
	wordCount int
}

// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the coordinator.
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets.
func coordinatorSock() string {
	s := "/var/tmp/824-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}
