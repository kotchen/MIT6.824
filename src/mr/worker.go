package mr

import (
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"
)

// Map functions return a slice of KeyValue.
type KeyValue struct {
	Key   string
	Value string
}

// use ihash(key) % NReduce to choose the reduce
// task number for each KeyValue emitted by Map.
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}

// main/mrworker.go calls this function.
func Worker(mapf func(string, string) []KeyValue,
	reducef func(string, []string) string) {

	// Your worker implementation here.
	//RequestTask(mapf, reducef)
	// uncomment to send the Example RPC to the coordinator.
	CallExample()

}

// example function to show how to make an RPC call to the coordinator.
//
// the RPC argument and reply types are defined in rpc.go.
func CallExample() {

	// declare an argument structure.
	args := ExampleArgs{}

	// fill in the argument(s).
	args.X = 99

	// declare a reply structure.
	reply := ExampleReply{}

	// send the RPC request, wait for the reply.
	// the "Coordinator.Example" tells the
	// receiving server that we'd like to call
	// the Example() method of struct Coordinator.
	ok := call("Coordinator.Example", &args, &reply)
	if ok {
		// reply.Y should be 100.
		fmt.Printf("reply.Y %v\n", reply.Y)
	} else {
		fmt.Printf("call failed!\n")
	}
}

func GetFileContent(fileName string) string {
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

	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	// 将文件内容转化为字符串
	return string(content)

}

type ByKey []KeyValue

func (o ByKey) Less(i, j int) bool { return o[i].Key < o[j].Key }
func (o ByKey) Len() int           { return len(o) }
func (o ByKey) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }

//func RequestTask(mapf func(string, string) []KeyValue, reducef func(string, []string) string) {
//args := RequestTaskArgs{}
//reply := RequestTaskReply{}
//ok := call("Coordinator.RequestTask", &args, &reply)
//if ok {
//switch reply.taskType {
//case Map:
//fileContent := GetFileContent(reply.fileName)
//mapRes := mapf(reply.fileName, fileContent)
//sort.Sort(ByKey(mapRes))

//oname := "mr-out-0"
//ofile, _ := os.Create(oname)

//case Reduce:
//reducef()
//}

//}
//}

func OnMapFinish(wordCount []KeyValue) {
	reply := MapTaskReply{}
	reply.success = true
	reply.wordCount = wordCount

}

func OnReduceFinish() {

}

// send an RPC request to the coordinator, wait for the response.
// usually returns true.
// returns false if something goes wrong.
func call(rpcname string, args interface{}, reply interface{}) bool {
	c, err := rpc.DialHTTP("tcp", "192.168.1.5"+":1234")
	//sockname := coordinatorSock()
	//c, err := rpc.DialHTTP("unix", sockname)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}
