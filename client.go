package zing

type Client struct {
	// my index number
	ID	int

	// the ip addres of my server 
	server	string

	// the ip address of all the server
	addressList []string
}

func InitializeClient(fileName String) *Client {
	// TODO: read metadata from a file to initialize the client

	return nil
}


func (self *Client) sendPrepare(prepare *Version) {

	// from the first index to the last
	for i := 0; i < self.addressList.count; i++ {
		address := self.addressList[i]
	}
}

func (self *Client) sendPush(push *Push) {
	
	// from the last index to the first
	for i := self.addressList.count - 1; i >= 0; i-- {
	
	}

}


