package zing

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

type Client struct {
	// my index number
	id int

	// the ip addres of my server
	server string

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

func (self *Client) Init() error {
	e := zing_init(0)
	if e != nil {
		return e
	}
	data := Data{Version{-1, -1, ""}, make([]string, 0)}
	e = writeFile(&data, METADATA_FILE)
	if e != nil {
		panic(e)
	}
	setAddressList([]string{self.server})

	setOwnIndex(0)
	log := []Push{Push{Version{-1, -1, ""}, make([]byte, 0)}}
	e = writeFile(&log, LOG_FILE)
	if e != nil {
		panic(e)
	}
	setVersion(0)
	return nil
}

func (self *Client) Clone(ip string) error {
	e := zing_init(0)
	if e != nil {
		return e
	}
	data := Data{Version{-1, -1, ""}, make([]string, 0)}
	e = writeFile(&data, METADATA_FILE)
	if e != nil {
		panic(e)
	}
	log := []Push{Push{Version{-1, -1, ""}, make([]byte, 0)}}
	e = writeFile(&log, LOG_FILE)
	if e != nil {
		panic(e)
	}
	setVersion(0)
	status := self.joinGroup(ip)
	if status == false {
		return fmt.Errorf("Cannot clone")
	}
	return nil

}

func (self *Client) Pull() error {
	if !IsServerRuning(self.server) {
		return fmt.Errorf("Server is not running")
	}
	if self.id == -1 {
		return fmt.Errorf("Not Initialized")
	}

	e := zing_pull("master")
	return e
}

func (self *Client) Commit(message string) error {
	if !IsServerRuning(self.server) {
		return fmt.Errorf("Server is not running")
	}
	if self.id == -1 {
		return fmt.Errorf("Not Initialized")
	}

	e := zing_commit(message)
	return e
}

func (self *Client) Add(filename string) error {
	if !IsServerRuning(self.server) {
		return fmt.Errorf("Server is not running")
	}
	if self.id == -1 {
		return fmt.Errorf("Not Initialized")
	}

	e := zing_add(filename)
	return e
}

func (self *Client) Push() error {
	if !IsServerRuning(self.server) {
		return fmt.Errorf("Server is not running")
	}
	
	er := self.Pull()  // pull before push
	if er != nil {
		return er
	}

	status := false		// check the prepare queue
	CheckPrepareQueue(self.server, self.server, &status)
	if status == false {
		return fmt.Errorf("Another push in progress, please pull and try again!")
	}

	e, data := zing_make_patch_for_push("master", "patch")
	if e != nil {
		return e
	} else if len(data) == 0 {		// don't need to issue a push if it doesn't need to
		return fmt.Errorf("Already up to data")
	}

	cversion := getVersion()
	prepare  := Version{NodeIndex: self.id, VersionIndex: cversion, NodeAddress: self.server}
	succ, bitMap := self.sendPrepare(&prepare)
	count := 0
	for i := 0; i < len(bitMap); i++ {
		if bitMap[i] {
			count++
		}
	}

	var bundle Asynchronous = Asynchronous{self.id, self.addressList, Push{}, bitMap}
	if (succ == false) || (count <= len(bitMap)/2) {
		fmt.Println("simultaneous push going on, please wait and push again")
		//self.sendPush(&Push{Change: prepare, Patch: make([]byte, 0)}, bitMap)
		bundle.Message = Push{Change: prepare, Patch: make([]byte, 0)}
	} else { 
		setVersion(cversion + 1)
		//self.sendPush(&Push{Change: prepare, Patch: data}, bitMap)
		bundle.Message = Push{Change: prepare, Patch: data}
	}
	succeed := false
	e = SendPushRequest(self.server, &bundle, &succeed)
	if e != nil {
		return fmt.Errorf("RPC function failed")
	}
	return nil
}

func (self *Client) sendPrepare(prepare *Version) (bool, []bool) {
	resultMap := make([]int, len(self.addressList))

	// send prepare message from first to last
	var group sync.WaitGroup
	for i := 0; i < len(self.addressList); i++ {
		address := self.addressList[i]
		group.Add(1)
		resultMap[i] = -1
		go SendPrepare(address, prepare, resultMap, i, &group)
	}
	group.Wait()

	liveBitMap := make([]bool, len(self.addressList))
	first := true
	index := -1
	for i := 0; i < len(self.addressList); i++ {
		if resultMap[i] != -1 {
			liveBitMap[i] = true
			if first {
				index = i
				first = false
			}
		} else {
			liveBitMap[i] = false
		}
	}

	fmt.Println(resultMap)
	var succeed bool = false
	if index != -1 {
		succeed = (resultMap[index] == 1)
		version := Version{NodeIndex: -1, VersionIndex: -1}
		group.Add(1)
		e := SendPrepare(self.addressList[index], &version, make([]int, 1), 0, &group)
		if e != nil {
			succeed = false
		}
	}
	return succeed, liveBitMap
}

func (self *Client) sendPush(push *Push, liveBitMap []bool) {
	if len(liveBitMap) != len(self.addressList) {
		panic("Node number not match")
	}

	var group sync.WaitGroup
	// send push changes from last to first
	for i := len(self.addressList) - 1; i >= 0; i-- {
		if liveBitMap[i] {
			address := self.addressList[i]
			succ := false
			group.Add(1)
			go SendPush(address, push, &succ, &group)
		}
	}

	group.Wait()
	return
}

func (self *Client) comeAlive() {
	succeed := false
	prepare := Version{NodeIndex: self.id, VersionIndex: -1, NodeAddress: self.server}
	pushes := Push{Change: prepare, Patch: []byte{}}
	bitMap := make([]bool, 0)

	time.Sleep(time.Second)
	for {
		succeed, bitMap = self.sendPrepare(&prepare)
		fmt.Println("Catching up")
		if succeed {
			break
		} else {
			self.sendPush(&pushes, bitMap)
		}

		live := false
		for _, value := range bitMap {
			if value {
				live = true
			}
		}
		if !live {
			succ := false
			SetReady(self.server, self.server, &succ)
			return
		}
	}

	ipList := getAddressList()
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
	pushes := Push{Change: prepare, Patch: []byte{}}
	ipList := make([]string, 0)
	bitMap := make([]bool, 0)

	for {
		ipList = make([]string, 0)
		err := RequestAddressList(address, make([]string, 0), &ipList)
		if err != nil {
			return false
		}
		self.addressList = ipList
		succeed, bitMap = self.sendPrepare(&prepare)
		fmt.Println("something wrong")
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
				commitChanges(pushList, self.id)
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

func (self *Client) Revert(commit_id string) error {
	if !IsServerRuning(self.server) {
		return fmt.Errorf("Server is not running")
	}
	if self.id == -1 {
		return fmt.Errorf("Not Initialized")
	}

	e := zing_revert(commit_id)
	return e
}

func (self *Client) Log() error {
	if !IsServerRuning(self.server) {
		return fmt.Errorf("Server is not running")
	}
	if self.id == -1 {
		return fmt.Errorf("Not Initialized")
	}
	e := zing_log()
	return e
}

func (self *Client) Status() error {
	if !IsServerRuning(self.server) {
		return fmt.Errorf("Server is not running")
	}
	if self.id == -1 {
		return fmt.Errorf("Not Initialized")
	}
	e := zing_status()
	return e
}
