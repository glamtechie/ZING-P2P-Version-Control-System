package zing

import (
	"sync"
)


type Server struct {
	// my index number
	ID	 int

	// the prepare message queue
	preQueue []Version	

	// the lock for changing prepare message queue 
	lock	 *sync.Mutex
}


func InitializeServer(fileName string) *Server {
	// TODO: read the metadata from a file to initialize the Server
	return nil
}


/*
 RPC function: ReceivePrepare

*/
func (self *Server) ReceivePrepare(prepare *Version, succ *bool) error {
	self.lock.Lock()
	defer self.lock.Unlock()

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
	for i := 0; i < len(self.preQueue); i++ {
		if VersionEquality(self.preQueue[i], push.version) {
			index = i
			break
		}
	}	
	if index == -1 {
		panic("No match prepare message")
	} else {
		// TODO: main logic here	
	}

	// after handling the push, clean the prepare queue
	self.lock.Lock()
	index = -1
	for i := 0; i < len(self.preQueue); i++ {
		if VersionEquality(self.preQueue[i], push.version) {
			index = i
			break
		}
	}
	if index == -1 {
		panic("No match prepare message")
	} else {
		self.preQueue = append(self.preQueue[:index], self.preQueue[index + 1:]...)	
	}
	self.lock.Unlock()
	
	*succ = true
	return nil
}



