package main

import (
	"log"
	"zing"
)

func main(){
	server:=zing.InitializeServer("info.txt")
	e:=zing.StartServer(server)
	if e!=nil{
		log.Fatal(e)
	}
}
