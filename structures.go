package zing

import (
	"net/rpc"
)

type Version struct {
	// the node that make this change
	nodeIndex int

	// the version of the change in this node
	version  int
}

func VersionEquality(a, b Version) bool {
	if a.nodeIndex == b.nodeIndex && a.version == b.version {
		return true
	} else {
		return false
	}
}


type Push struct {
	// the verstion corresponded to this push
	version Version

	// a list of diff files, map from filename to diff.
	diffList map[string]string
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


