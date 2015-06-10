package zing

import (
	"net/rpc"
	"sync"
)

const (
	INVALIDIP = "0.0.0.0:0"
)

type Version struct {
	// the node that make this change
	NodeIndex int

	// the version of the change in this node
	VersionIndex int

	// the adress of sending node
	NodeAddress string
}

func VersionEquality(a, b Version) bool {
	if a.NodeIndex == b.NodeIndex && a.VersionIndex == b.VersionIndex && a.NodeAddress == b.NodeAddress {
		return true
	} else {
		return false
	}
}

func IsServerRuning(address string) bool {
	var group sync.WaitGroup
	version := Version{NodeIndex: -1, VersionIndex: -1, NodeAddress: INVALIDIP}
	value   := -1
	group.Add(1)
	e       := SendPrepare(address, &version, &value, &group)
	if e != nil {
		return false
	} else {
		return true
	}
}

type Push struct {
	// the verstion corresponded to this push
	Change Version

	// a list of diff files, map from filename to diff.
	Patch []byte
}

func SendPrepare(address string, prepare *Version, indicate *int, w *sync.WaitGroup) error {
	defer w.Done()
	conn, e := rpc.DialHTTP("tcp", address)
	if e != nil {
		return e
	}

	var succ bool
	e = conn.Call("Server.ReceivePrepare", prepare, &succ)
	if e != nil {
		conn.Close()
		return e
	}

	if succ {
		*indicate = 1
	} else {
		*indicate = 0
	}
	return conn.Close()
}

func SendPush(address string, push *Push, succ *bool, w *sync.WaitGroup) error {
	defer w.Done()
	conn, e := rpc.DialHTTP("tcp", address)
	if e != nil {
		return e
	}

	e = conn.Call("Server.ReceivePush", push, succ)
	if e != nil {
		conn.Close()
		return e
	}
	return conn.Close()
}

func SetReady(address, ip string, succ *bool) error {
	conn, e := rpc.DialHTTP("tcp", address)
	if e != nil {
		return e
	}

	e = conn.Call("Server.ReceiveReady", ip, succ)
	if e != nil {
		conn.Close()
		return e
	}
	return conn.Close()
}

// for asynchronous Push
type Asynchronous struct {
	// the client object
	Object 	*Client
	// the push message
	Message *Push
	// the live bit map
	LiveMap []bool
}

func SendPushRequest(address string, bundle *Asynchronous, succ *bool) error {
	conn, e := rpc.DialHTTP("tcp", address)
	if e != nil {
		return e
	}

	e = conn.Call("Server.AsynchronousPush", bundle, succ)
	if e != nil {
		conn.Close()
		return e
	}
	return conn.Close()
}

func RequestAddressList(address string, argList []string, ipList *[]string) error {
	conn, e := rpc.DialHTTP("tcp", address)
	if e != nil {
		return e
	}

	e = conn.Call("Server.ReturnAddressList", argList, ipList)
	if e != nil {
		conn.Close()
		return e
	}
	return conn.Close()
}

func CheckPrepareQueue(address, ip string, result *bool) error {
	conn, e := rpc.DialHTTP("tcp", address)
	if e != nil {
		return e
	}

	e = conn.Call("Server.PrepareQueueCheck", ip, result)
	if e != nil {
		conn.Close()
		return e
	}
	return conn.Close()
}

// New node read missing data from another node.
func ReadMissingData(address string, ver Version, pushes *[]Push) error {
	conn, e := rpc.DialHTTP("tcp", address)
	if e != nil {
		return e
	}

	e = conn.Call("Server.ReturnMissingData", ver, pushes)
	if e != nil {
		conn.Close()
		return e
	}
	return conn.Close()
}
