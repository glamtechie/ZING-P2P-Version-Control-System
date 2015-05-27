package zing

import (
	"os"
	"strconv"
	"io/ioutil"
	"bytes"
	"encoding/gob"
)

const (
	file="data.txt"
)
func GetIPList(filename string) []string {
	return []string {"137.110.91.41:27321", "137.110.90.199:27321", "137.110.92.134:27321"}
}

func GetIndexNumber(filename string) int {
	return 2
}

func GetVersionNumber(filename string) int {
	file, e := os.Open(filename)
	if e != nil {
		//panic("Can't open the metadata file")
		return 0
	}

	data := make([]byte, 8)
	file.Read(data)
	result, _ := strconv.ParseInt(string(data), 10, 32)
	return int(result)
}

func SetVersionNumber(filename string, version int) {
	file, e := os.Create(filename)
	if e != nil {
		panic("Can't open the metadata file")
	}

	data := []byte(strconv.Itoa(version))
	file.Write(data)
	return
}


//file stuff

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

func writeFile(data interface{}, filename string)error{
	bts,e:=getBytes(data)
	if e!=nil{
		return e
	}
	err := ioutil.WriteFile(filename, bts, 0644)
	if err!=nil{
		return err
	}

	return nil
}

func readFile(filename string, data interface{})(error){
	bts, err := ioutil.ReadFile(filename)
    if err != nil {
    	return err
    }
    err = getInterface(bts, &data)
	if err != nil {
		return err
	}
	return nil
}

type Data struct {
	Myown Version
	All []string
}

func getIPs()([]string) {
	var data Data
	e:=readFile(file,&data)
	if e!=nil{
		panic(e)
	}
	return data.All
}

func getOwnIndex()(int){
	var data Data
	e:=readFile(file,&data)
	if e!=nil{
		panic(e)
	}
	return data.Myown.NodeIndex
}

func getVersion()(int){
	var data Data
	e:=readFile(file,&data)
	if e!=nil{
		panic(e)
	}
	return data.Myown.VersionIndex
}

func setIPs(iplist []string){
	var data Data
	e:=readFile(file,&data)
	if e!=nil{
		panic(e)
	}
	data.All=iplist
	e=writeFile(&data,file)
	if e!=nil{
		panic(e)
	}
}

func setVersion(ver int){
	var data Data
	e:=readFile(file,&data)
	if e!=nil{
		panic(e)
	}
	data.Myown.VersionIndex=ver
	e=writeFile(&data,file)
	if e!=nil{
		panic(e)
	}
}

func setOwnIndex(me int){
	var data Data
	e:=readFile(file,&data)
	if e!=nil{
		panic(e)
	}
	data.Myown.NodeIndex=me
	e=writeFile(&data,file)
	if e!=nil{
		panic(e)
	}
}

