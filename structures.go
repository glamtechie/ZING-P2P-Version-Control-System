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
}

func VersionEquality(a, b Version) bool {
	if a.NodeIndex == b.NodeIndex && a.VersionIndex == b.VersionIndex {
		return true
	} else {
		return false
	}
}


type Push struct {
	// the verstion corresponded to this push
	Change Version

	// a list of diff files, map from filename to diff.
	DiffList map[string]string
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

type IPChange struct {
	// the index of the machine
	Index	int

	// the ip address of this machine
	IP	string
}


func SendIPChange(address string, ipchange *IPChange, succ *bool) error {
	conn, e := rpc.DialHTTP("tcp", address)
	if e != nil {
		return e
	}

	e = conn.Call("Server.ReceiveIPChange", ipchange, succ)
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






