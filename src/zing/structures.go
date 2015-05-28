package zing

import (
	"net/rpc"
)


const (
	INVALIDIP = "0.0.0.0:0"
)

type Version struct {
	// the node that make this change
	NodeIndex    int

	// the version of the change in this node
	VersionIndex int

	// the adress of sending node
	NodeAddress  string
}

func VersionEquality(a, b Version) bool {
	if a.NodeIndex == b.NodeIndex && a.VersionIndex == b.VersionIndex && a.NodeAddress == b.NodeAddress {
		return true
	} else {
		return false
	}
}



type Push struct {
	// the verstion corresponded to this push
	Change Version

	// a list of diff files, map from filename to diff.
	Patch []byte
}


func SendPrepare(address string, prepare *Version, succ *bool) error {
	conn, e := rpc.DialHTTP("tcp", address)
	if e != nil {
        	return e
    	}

    	e = conn.Call("Server.ReceivePrepare", prepare, succ)
    	if e != nil {
        	conn.Close()
        	return e
    	}
    	return conn.Close()
}

func SendPush(address string, push *Push, succ *bool) error {
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

func RequestAddressList(address, ip string, ipList *[]string) error {
	conn, e := rpc.DialHTTP("tcp", address)
	if e != nil {
		return e
	}

	e = conn.Call("Server.ReturnAddressList", ip, ipList)
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
func ReadMissingData(address string) error {
	return nil
}



