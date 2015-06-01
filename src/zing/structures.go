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

func IsServerRuning(address string) bool {
	version := Version{NodeIndex: -1, VersionIndex: -1, NodeAddress: INVALIDIP}
	succ    := false
	e       := SendPrepare(address, &version, &succ)
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



