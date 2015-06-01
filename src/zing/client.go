package zing

import (
	"fmt"
	"net"
	"os"
)

type Client struct {
	// my index number
	id	int

	// the ip addres of my server
	server	string

	// the ip address of all the server
	addressList []string
}


func InitializeClient() *Client {
	client := Client{}
	if _, err := os.Stat(METADATA_FILE); os.IsNotExist(err) {
		client.id = -1
		client.addressList = make([]string, 0)
	} else {
		client.id = getOwnIndex()
		client.addressList = getAddressList()
	}
	addrs, _ := net.InterfaceAddrs()
    for _, address := range addrs {
        if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
            	client.server = ipnet.IP.String() + ":27321"
            	break
            }
        }
    }
	return &client
}

func (self *Client) Init() error{
	e:=zing_init(0)
	if e!=nil{
		return e
	}
    data := Data{Version{-1,-1,""},make([]string,0)}
    e 	  = writeFile(&data, METADATA_FILE)
    if e != nil{
    	panic(e)
    }

    setAddressList([]string{self.server})

	setOwnIndex(0)
	log  := []Push{ Push{Version{-1,-1, self.server}, make([]byte,0)} }
	e     = writeFile(&log, LOG_FILE)
	if e != nil {
		panic(e)
	}
	setVersion(0)
	return nil
}

func (self *Client) Clone(ip string) error{
	e:=zing_init(0)
	if e != nil{
		return e
	}
	data := Data{Version{-1,-1,""},make([]string,0)}
    e 	  = writeFile(&data, METADATA_FILE)
    if e != nil {
    	panic(e)
    }
	writeLog(Push{Version{-1,-1, self.server},make([]byte,0)})
	setVersion(0)
	status:=self.joinGroup(ip)
	if status==false{
		return fmt.Errorf("Cannot clone")
	}
	return nil

}

func (self *Client) Pull() error {
	if self.id == -1 {
		return fmt.Errorf("Not Initialized")
	}

	e := zing_pull("master")
	return e
}

func (self *Client) Commit(message string) error {
	if self.id == -1 {
		return fmt.Errorf("Not Initialized")
	}

	e := zing_commit(message)
	return e
}

func (self *Client) Add(filename string) error {
	if self.id == -1 {
		return fmt.Errorf("Not Initialized")
	}

	e := zing_add(filename)
	return e
}


func (self *Client) Push() error {
	er:=self.Pull()
	if er!=nil{
		return er
	}
	status := false
	CheckPrepareQueue(self.server, self.server, &status)
	if status == false {
		return fmt.Errorf("Another push in progress, please pull and try again!")
	}

	cversion     := getVersion()
	prepare      := Version{NodeIndex: self.id, VersionIndex: cversion, NodeAddress: self.server}
	succ, bitMap := self.sendPrepare(&prepare)
	count := 0
	for i := 0; i < len(bitMap); i++ {
		if bitMap[i] {
			count++
		}
	}
	if (succ == false) || (count <= len(bitMap) / 2) {
		fmt.Println("Push failed, abort")
		self.sendAbort(bitMap, cversion)
		return nil
	}
	e, data := zing_make_patch_for_push("master", "patch")
	if e != nil{
		return e
	}

	setVersion(cversion + 1)
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
	prepare := Version{NodeIndex: self.id, VersionIndex: -1, NodeAddress: self.server}
	pushes  := Push{Change: prepare, Patch: []byte{}}
	bitMap  := make([]bool, 0)

	for {
		succeed, bitMap = self.sendPrepare(&prepare)
		if succeed {
			break
		} else {
			self.sendPush(&pushes, bitMap)
		}
	}

	ipList  := getAddressList()
	resList := make([]string, 0)
	for key, ip := range ipList {
		if bitMap[key] {
			resList = make([]string, 0)
			RequestAddressList(ip, ipList, &resList)

			localVer := getLastVer()
			pushList := make([]Push, 0)
			ReadMissingData(ip, localVer, &pushList)
			if len(pushList) == 1 && len(pushList[0].Patch) == 0 {
				pushList = getPushDiff(pushList[0].Change)
				ReadMissingData(ip, localVer, &pushList)
			} else {
				commitChanges(pushList, self.id)
			}
		}
	}

	if len(resList) > len(ipList) {
		setAddressList(resList)
	}

	succ := false
	SetReady(self.server, self.server, &succ)
	self.sendPush(&pushes, bitMap)
	return
}


func (self *Client) joinGroup(address string) bool {
 	succeed := false
	prepare := Version{NodeIndex: -1, VersionIndex: -1, NodeAddress: self.server}
	pushes  := Push{Change: prepare, Patch: []byte{}}
	ipList  := make([]string, 0)
	bitMap  := make([]bool, 0)

	for {
		ipList  = make([]string, 0)
		err    := RequestAddressList(self.server, make([]string, 0), &ipList)
		if err != nil {
			return false
		}
		succeed, bitMap = self.sendPrepare(&prepare)
		if succeed {
			break
		} else {
			self.sendPush(&pushes, bitMap)
		}
	}
	for key, ip := range ipList {
		if bitMap[key] {
			localVer := getLastVer()
			pushList := make([]Push, 0)
			e := ReadMissingData(ip, localVer, &pushList) // read from any of the node
			if e == nil {
				break
			}
		}
	}

	ipList = append(ipList, self.server)
	setAddressList(ipList)
	setVersion(0)
	setOwnIndex(len(ipList) - 1)

	succ := false
	SetReady(self.server, self.server, &succ)
	self.sendPush(&pushes, bitMap)
	return true
}

func (self *Client) Revert(commit_id string)(error){
	if self.id == -1 {
		return fmt.Errorf("Not Initialized")
	}

	e:=zing_revert(commit_id)
	return e
}

func (self *Client) Log()(error){
	if self.id == -1 {
		return fmt.Errorf("Not Initialized")
	}
	e:=zing_log()
	return e
}

func (self *Client) Status()(error){
	if self.id == -1 {
		return fmt.Errorf("Not Initialized")
	}
	e:=zing_status()
	return e
}


