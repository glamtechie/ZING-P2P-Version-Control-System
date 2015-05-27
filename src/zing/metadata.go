package zing

import (
	"os"
	"strconv"
)

func GetIPList(filename string) []string {
	//return []string{}
	return []string{"137.110.91.41:27321", "137.110.90.199:27321", "137.110.92.134:27321" }
}

func GetIndexNumber(filename string) int {
	return 1
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
