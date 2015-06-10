package main

import (
	"log"
	"zing"
    "flag"
    "os"
)

func main(){
    flag.Parse()
    args := flag.Args()
    if len(args) < 1 {
        os.Exit(1)
    }

    port:=args[0]
	server:=zing.InitializeServer(port)
	e:=zing.StartServer(server)
	if e!=nil{
		log.Fatal(e)
	}
}
