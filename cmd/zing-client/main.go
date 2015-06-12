package main

import (
	"zing"
	"flag"
	"os"
	"fmt"
)

func runPrompt (client *zing.Client, cmd string){
	var v error
	switch cmd {

	case "log":
		v=client.Log()
		logError(v)
	case "status":
		v=client.Status()
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

func runCmd(client *zing.Client, cmd string, args []string){
	var v error
	switch cmd{
	case "init":
		v=client.Init(args[0])
		logError(v)
	case "add":
		v=client.Add(args[0])
		logError(v)
	case "rm":
		v=client.Rm(args[0])
		logError(v)
	case "revert":
		v=client.Revert(args[0])
		logError(v)
	case "clone":
		v=client.Clone(args[0], args[1])
		logError(v)
	case "commit":
		v=client.Commit(args[0])
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
	client:=zing.InitializeClient()
	cmdArgs := args[1:]
	if len(cmdArgs) == 0 {
		runPrompt(client,cmd)
		fmt.Println()
	} else {
		runCmd(client, cmd, cmdArgs)
	}
}
