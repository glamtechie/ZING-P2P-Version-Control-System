package main

import (
	"log"
	"zing"
    "flag"
    "os"
)

func main(){

	server:=zing.InitializeServer()
	e:=zing.StartServer(server)
	if e!=nil{
		log.Fatal(e)
	}
}
