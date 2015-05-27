package zing

import (
	"os"
	"strconv"
)

func GetIPList(filename string) []string {
	return []string{}
}

func GetIndexNumber(filename string) int {
	return 0
}

func GetVersionNumber(filename string) int {
	file, e := os.Open(filename)
	if e != nil {
		panic("Can't open the metadata file")
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
