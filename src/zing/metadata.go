package zing

import (
	"os"
	"strconv"
	"fmt"
	"bytes"
)

func GetIPList(filename string) []string {
	//return []string{}
	return []string{"137.110.90.199:27321", "137.110.90.91:27321" }
}

func GetIndexNumber(filename string) int {
	return 0
}

func GetVersionNumber(filename string) int {
	file, e := os.Open(filename)
	if e != nil {
		return 0
	}
	data := make([]byte, 16)
	file.Read(data)

	tail 	  := bytes.Index(data, []byte{0})
	result, _ := strconv.Atoi(string(data[:tail]))
	file.Close()
	return result
}

func SetVersionNumber(filename string, version int) {
	file, e := os.Create(filename)
	if e != nil {
		panic("Can't open the metadata file")
	}

	data := []byte(strconv.Itoa(version))
	file.Write(data)
	file.Close()
	return 
}
