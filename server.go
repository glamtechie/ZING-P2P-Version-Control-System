package zing

type Server struct {
	// my index number
	ID	    int

	// the ip address of all the server
	// whether server or client should maintain this, not both
	// addressList []string

	// the prepare message queue
	prepare	    *Version
}



func InitializeServer(fileName string) *Server {
	// TODO: read the metadata from a file to initialize the Server

	return nil
}


/*
 RPC function: ReceivePrepare

*/
func (self *Server) ReceivePrepare(prepare *Version, succ *bool) error {
	if self.prepare == nil {
		self.prepare = prepare
		*succ = true
	} else {
		*succ = false
	}

	return nil
}

/*
 RPC function: ReceivePush

*/
func (self *Server) ReceivePush(push *Push, succ *bool) error {
	if self.prepare == nil {
		panic("Prepare is nil")
	} else {
		index   := self.prepare.nodeIndex
		version := self.prepare.version	
		if index != push.version.nodeIndex || version != push.version.version {
			panic("Don't match the prepare message")
		} else {
			// TODO: main logic here
		}
	}

	*succ = true
	self.prepare = nil
	return nil
}



