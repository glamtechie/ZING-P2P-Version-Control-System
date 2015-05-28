package zing

import (
	"sync"
	"fmt"
	"net"
	"net/rpc"
	"net/http"
)


type Server struct {
	// my index number
	id	 int

	// my ip address
	address  string

	// the prepare message queue
	preQueue []Version	

	// the lock for changing prepare message queue 
	lock	 *sync.Mutex

	// ready to serve or not
	ready	 bool
}

// global variable
var (
	GlobalBuffer []*Push
	IndexList    []int
)


func InitializeServer(fileName string) *Server {
	server := Server{}
	server.id = GetIndexNumber(fileName)
	
	addressList := GetIPList(fileName)
	server.address  = addressList[server.id]
	server.preQueue = make([]Version, 0)
	server.lock     = &sync.Mutex{}
	server.ready    = true
	return &server
}

func StartServer(instance *Server) error {
	server := rpc.NewServer()
	server.Register(instance)
	
	l, e := net.Listen("tcp", instance.address)
	if e != nil {
		return e
	}
	return http.Serve(l, server)
}


/*
 RPC function: ReceivePrepare
*/
func (self *Server) ReceivePrepare(prepare *Version, succ *bool) error {
	self.lock.Lock()
	defer self.lock.Unlock()

	fmt.Printf("Receive the prepare, Node: %d, Version: %d\n", prepare.NodeIndex, prepare.VersionIndex)

	// the server is not ready
	if !self.ready {
		return fmt.Errorf("Server not ready")
	}
	// it is a dummy version
	if prepare.NodeIndex == -1 && prepare.VersionIndex == -1 {
		return nil
	}

	if len(self.preQueue) == 0 {
		*succ = true
	} else {
		*succ = false
	}
	self.preQueue = append(self.preQueue, *prepare)
	return nil
}

func processChanges(push *Push, index int) []*Push {
	insertPoint := -1
	for key, value := range IndexList {
		if index < value {
			insertPoint = key
			break
		}
	}

	if insertPoint == -1 {
		insertPoint = len(IndexList)
	}
	IndexList    = append(IndexList[:insertPoint],    append([]int{index}, IndexList[insertPoint:]...)...)
	GlobalBuffer = append(GlobalBuffer[:insertPoint], append([]*Push{push}, GlobalBuffer[insertPoint:]...)...)

	if IndexList[0] == 0 {
		cutPoint := 1
		for i := 1; i < len(IndexList); i++ {
			if IndexList[i] - IndexList[i - 1] != 1 {
				cutPoint = i 
				break
			}
		}
		results := GlobalBuffer[:cutPoint]

		IndexList    = IndexList[cutPoint:]
		GlobalBuffer = GlobalBuffer[cutPoint:]
		return results
	} else {
		return make([]*Push, 0)
	}
}


func (self *Server) commitChanges(pushes []*Push) error {
	// commit the pushes to the file system
	for _, push := range pushes {
		if len(push.Patch) == 0 {
			continue
		} 

		e := zing_process_push("patch", push.Patch)
		if e != nil {
			panic("commit change error")
		}
	}
	return nil
}

/*
 RPC function: ReceivePush
*/
func (self *Server) ReceivePush(push *Push, succ *bool) error {
	var index int = -1
	var pushes []*Push 

	fmt.Printf("Receive the Push from Node: %d, Version: %d\n", push.Change.NodeIndex, push.Change.VersionIndex)
	fmt.Printf("Patch length: %d\n", len(push.Patch))

	self.lock.Lock()
	defer self.lock.Unlock()

	for i, prepare := range self.preQueue {
		if VersionEquality(prepare, push.Change) {
			index = i
			break
		}
	}
	if index == -1 {
		panic("No match prepare message")
	} else {
		pushes = processChanges(push, index)
	}
	
	// commit the changes
	self.commitChanges(pushes)
	if len(pushes) > 0 {
		self.preQueue = self.preQueue[len(pushes):]
	}
	*succ = true
	return nil
}

/*
 RPC function: ReceiveReady
*/
func (self *Server) ReceiveReady(address string, succ *bool) error {
	// this function should only called by its own client
	if address != self.address {
		panic("Not come from my self")
	}

	self.lock.Lock()
	defer self.lock.Unlock()

	self.ready = true
	*succ = true
	return nil
}


/*
 RPC function: RecevieIP
*/
func (self *Server) ReceiveIPChange(ipchange *IPChange, ipList *[]string) error {
	// TODO: write the ip changes to the metadata file
	return nil
}

func (self *Server) ReceiveIPListRequest(address string, ipList *[]string) error {
	list := GetIPList("info.txt")
	*ipList = list
	return nil
}


func (self *Server) PrepareQueueCheck(address string, result *bool) error {
	if address != self.address {
		panic("Not come from my self")
	}

	self.lock.Lock()
	defer self.lock.Unlock()

	if len(self.preQueue) != 0 {
		*result = false
	} else {
		*result = true
	}
	return nil
}
