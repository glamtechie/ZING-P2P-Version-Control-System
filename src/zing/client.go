package zing

import (
	"errors"
)

type Client struct {
	// my index number
	id	int

	// the ip addres of my server
	server	string

	// the ip address of all the server
	addressList []string
}


func InitializeClient(fileName string) *Client {
	client := Client{}
	client.id = GetIndexNumber(fileName)
	client.addressList = GetIPList(fileName)

	client.server = client.addressList[client.id]
	//SetVersionNumber("VersionNumber", 0)
	return &client
}

func (self *Client) Pull() error {
	e := zing_pull("master")
	return e
}

func (self *Client) Commit() error {
	e := zing_commit()
	return e
}

func (self *Client) Add(filename string) error {
	e := zing_add(filename)
	return e
}


func (self *Client) Push() error {
	status := false
	CheckPrepareQueue(self.server, self.server, &status)
	if status == false {
		return errors.New("Another push in progress, please pull and try again!")
	}
	
	cversion     := GetVersionNumber(".zing/VersionNumber");
	prepare      := Version{NodeIndex: self.id, VersionIndex: cversion, NodeAddress: self.server}
	succ, bitMap := self.sendPrepare(&prepare)
	count := 0
	for i := 0; i < len(bitMap); i++ {
		if bitMap[i] {
			count++
		}
	}
	if (succ == false) || (count <= len(bitMap) / 2) {
		 self.sendAbort(bitMap, cversion)
		 return nil
	}

	e, data := zing_make_patch_for_push("master", "patch")
	if e != nil{
		return e
	}

	SetVersionNumber(".zing/VersionNumber", cversion + 1)
	self.sendPush(&Push{Change: prepare, Patch: data}, bitMap)
	return nil
}


func (self *Client) sendPrepare(prepare *Version) (bool, []bool) {
	firstNode  := false
	//firstIndex := -1
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
			if !firstNode {
				succeed    = succ
				firstNode  = false
				//firstIndex = i
			}
			liveBitMap[i] = true
		}
	}

	/*
	if firstIndex != -1 {
		version := Version{NodeIndex: -1, VersionIndex: -1}
		succ    := false
		e       := SendPrepare(self.addressList[firstIndex], &version, &succ)
		if e != nil {
			succeed = false
		}
	}*/

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

func (self *Client) sendAbort(liveBitMap []bool, cversion int) {
	ar := make([]byte, 0)
	vr := Version{NodeIndex: self.id, VersionIndex: cversion, NodeAddress: self.server}
	self.sendPush(&Push{Change: vr, Patch: ar}, liveBitMap)
}


func (self *Client) comeAlive() {
	succeed := false
	bitMap  := make([]bool, 0)
	prepare := Version{NodeIndex: self.id, VersionIndex: -1, NodeAddress: self.server}
	pushes  := Push{Change: prepare, Patch: []byte{}}
	
	for {		// try to go through this competing push
		succeed, bitMap = self.sendPrepare(&prepare)
		if succeed {
			break
		} else {
			self.sendPush(&pushes, bitMap)
		}
	}

	// here we read missing data from some node
	for key, address := range self.addressList {
		if bitMap[key] {
			e := ReadMissingData(address)
			if e == nil {
				break
			}
		}
	}

	succ := false
	SetReady(self.server, self.server, &succ)
	self.sendPush(&pushes, bitMap)
}

// new node join the group
// don't use this function yet
func (self *Client) joinGroup(address string) bool {	
	ipList := make([]string, 0)
	err    := RequestAddressList(address, self.server, &ipList)
	if err != nil {
		return false
	}
 
 	succeed := false
	bitMap  := make([]bool, 0)
	prepare := Version{NodeIndex: -1, VersionIndex: -1, NodeAddress: self.server}
	pushes  := Push{Change: prepare, Patch: []byte{}}

	// try to go through this competing push, only try once
	succeed, bitMap = self.sendPrepare(&prepare)
	if !succeed {
		self.sendPush(&pushes, bitMap)
		return false
	}

	// here we read missing data from some node
	for key, ip := range ipList {
		if bitMap[key] {
			e := ReadMissingData(ip)
			if e == nil {
				break
			}
		}
	}

	self.id = len(ipList)		// the last one
	self.addressList = append(ipList, self.server)

	succ := false
	SetReady(self.server, self.server, &succ)
	self.sendPush(&pushes, bitMap)
	return true
}
