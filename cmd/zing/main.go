package main

import (
	"zing"
	"flag"
	"os"
	"fmt"
)

func runPrompt (client *Client, cmd string){
	var v error
	switch cmd {
	case "commit":
		v=client.Commit()
		logError(v)
	case "pull":
		v=client.Pull()
                logError(v)
	case "push":
		v=client.Push()
                logError(v)
	
	default:
		logError(fmt.Errorf("bad command, try \"help\"."))
	}
}

func runCmd(client *Client, cmd string, args []string){
	var v error
	switch cmd{
	case "add":
		v=client.Add(args[0])
		logError(v)
	default:
		logError(fmt.Errorf("bad command, try \"help\"."))
        }

}
const help = `Usage:
   zing <command> [command <args...>]
With no command specified to enter interactive mode. 
` + cmdHelp

const cmdHelp = `Command list:
	add <filename>
	commit
	pull
	push
`

func logError(e error) {
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
	}
}
func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, help)
		os.Exit(1)
	}

	cmd:=args[0]
	client:=InitializeClient("info.txt")
	cmdArgs := args[1:]
	if len(cmdArgs) == 0 {
		runPrompt(client,cmd)
		fmt.Println()
	} else {
		runCmd(client, cmd, cmdArgs)
	}
}
