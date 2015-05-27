package zing

type Client struct {
	// my index number
	id	int

	// the ip addres of my server
	server	string

	// the ip address of all the server
	addressList []string

	//current version
	cversion int
}


func InitializeClient(fileName string) *Client {
	// TODO: read metadata from a file to initialize the client
	return nil
}

func (self *Client) Pull()error{
	e:=zing_pull("master")
	return e
}

func (self *Client) Commit()error{
	e:=zing_commit()
	return e
}

func (self *Client) Add(filename string)error{
	e:=zing_add(filename)
	return e
}

func (self *Client) Push()error{
	status,bitMap:=sendPrepare(&Version{self.id,self.cversion+1})
}

func (self *Client) sendPrepare(prepare *Version) (bool, []bool) {
	firstNode  := true
	firstIndex := -1
	succeed    := false
	liveBitMap := make([]bool, len(self.addressList))

	// send prepare message from first to last
	for i := 0; i < len(self.addressList); i++ {
		address := self.addressList[i]
		succ    := false
		e       := SendPrepare(address, prepare, &succ)

		if e != nil {
			liveBitMap[i] = false
		} else {
			if firstNode {
				succeed    = succ
				firstNode  = false
				firstIndex = i
			}
			liveBitMap[i] = true
		}
	}

	if firstIndex != -1 {
		version := Version{NodeIndex: -1, VersionIndex: -1}
		succ    := false
		e       := SendPrepare(self.addressList[firstIndex], &version, &succ)
		if e != nil {
			succeed = false
		}
	}

	// TODO: write self.addressList back to metadata file
	return succeed, liveBitMap
}


func (self *Client) sendPush(push *Push, liveBitMap []bool) {
	if len(liveBitMap) != len(self.addressList) {
		panic("Node number not match")
	}

	// send push changes from last to first
	for i := len(self.addressList) - 1; i >= 0; i-- {
		if liveBitMap[i] {
			address := self.addressList[i]
			succ    := false
			SendPush(address, push, &succ)
		}
	}
}



func (self *Client) comeAlive() {
	ipchange := IPChange{Index: self.ID, IP: self.server}

	// tell everybody I am alive
	for i := 0; i < len(self.addressList); i++ {
		address := self.addressList[i]
		succ    := false
		SendIPChange(address, &ipchange, &succ)
	}

	succeed := false
	bitMap  := make([]bool, 0)
	prepare := Version{NodeIndex: self.ID, VersionIndex: -1}
	pushes  := Push{Change: prepare, DiffList: make(map[string]string)}
	for {
		succeed, bitMap = self.sendPrepare(&prepare)
		if succeed {
			break
		} else {
			self.sendPush(&pushes, bitMap)
		}
	}

	succ := false
	SetReady(self.server, self.server, &succ)
	self.sendPush(&pushes, bitMap)
}


