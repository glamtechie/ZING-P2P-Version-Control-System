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
	ID	 int

	// my ip address
	address  string

	// the prepare message queue
	preQueue []Version	

	// the lock for changing prepare message queue 
	lock	 *sync.Mutex

	// ready to serve or not
	ready	 bool
}


func InitializeServer(fileName string) *Server {
	// TODO: read the metadata from a file to initialize the Server
	return nil
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


/*
 RPC function: ReceivePush
*/
func (self *Server) ReceivePush(push *Push, succ *bool) error {	
	index := -1
	self.lock.Lock()
	for i, prepare := range self.preQueue {
		if VersionEquality(prepare, push.Change) {
			index = i
			break
		}
	}
	self.lock.Unlock()
	if index == -1 {
		panic("No match prepare message")
	} else {
		// TODO: main logic here	
	}

	// after handling the push, clean the prepare queue
	self.lock.Lock()
	index = -1
	for i, prepare := range self.preQueue {
		if VersionEquality(prepare, push.Change) {
			index = i
			break
		}
	}
	if index == -1 {
		panic("No match prepare message")
	} else {
		self.preQueue = append(self.preQueue[:index], self.preQueue[index+1:]...)
	}
	self.lock.Unlock()
	
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
func (self *Server) ReceiveIPChange(ipchange *IPChange, succ *bool) error {
	// TODO: write the ip changes to the metadata file
	*succ = true
	return nil
}



