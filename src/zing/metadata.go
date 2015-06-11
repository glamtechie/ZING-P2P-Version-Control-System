package zing

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
)

const (
	METADATA_FILE = ".zing/metadata.txt"
	LOG_FILE=".zing/log.txt"
)

//file stuff
func writeLog(push Push){
	var data []Push
	e := readFile(LOG_FILE, &data)
	if e != nil {
		panic(e)
	}

	data=append(data,push)

	e = writeFile(&data, LOG_FILE)
	if e != nil {
		panic(e)
	}
}

func getPushDiff(ver Version)([]Push){
	var data []Push
	e := readFile(LOG_FILE, &data)
	if e != nil {
		panic(e)
	}

	var i int
	for i=0;i<len(data)-1;i++{
		if data[i].Change.NodeIndex==ver.NodeIndex && data[i].Change.VersionIndex==ver.VersionIndex{
			return data[i+1:]
		}
	}

	return make([]Push,0)
}

func getLastVer()(Version){
	var data []Push
	e := readFile(LOG_FILE, &data)
	if e != nil {
		panic(e)
	}

	return data[len(data)-1].Change
}

func getBytes(data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func getInterface(bts []byte, data interface{}) error {
	buf := bytes.NewBuffer(bts)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(data)
	if err != nil {
		return err
	}
	return nil
}

func writeFile(data interface{}, filename string) error {
	bts, e := getBytes(data)
	if e != nil {
		return e
	}
	err := ioutil.WriteFile(filename, bts, 0644)
	if err != nil {
		return err
	}

	return nil
}

func readFile(filename string, data interface{}) error {
	bts, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = getInterface(bts, data)
	if err != nil {
		return err
	}
	return nil
}

type Data struct {
	Myown Version
	All   []string
}

func getAddressList() []string {
	var data Data
	e := readFile(METADATA_FILE, &data)
	if e != nil {
		panic(e)
	}
	return data.All
}

func getOwnIndex() int {
	var data Data
	e := readFile(METADATA_FILE, &data)
	if e != nil {
		panic(e)
	}
	return data.Myown.NodeIndex

}

func getVersion() int {
	var data Data
	e := readFile(METADATA_FILE, &data)
	if e != nil {
		panic(e)
	}
	return data.Myown.VersionIndex
}

func setAddressList(iplist []string) {
	var data Data
	e := readFile(METADATA_FILE, &data)
	if e != nil {
		panic(e)
	}
	data.All = iplist
	e = writeFile(&data, METADATA_FILE)
	if e != nil {
		panic(e)
	}
}

func setVersion(ver int) {
	var data Data
	e := readFile(METADATA_FILE, &data)
	if e != nil {
		panic(e)
	}
	data.Myown.VersionIndex = ver
	e = writeFile(&data, METADATA_FILE)
	if e != nil {
		panic(e)
	}
}

func setOwnIndex(me int) {
	var data Data
	e := readFile(METADATA_FILE, &data)
	if e != nil {
		panic(e)
	}
	data.Myown.NodeIndex = me
	e = writeFile(&data, METADATA_FILE)
	if e != nil {
		panic(e)
	}
}
