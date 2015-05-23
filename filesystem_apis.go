package main

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
)

func zing_init(id int) {
	out, err := exec.Command("/bin/sh", "filesystem_scripts/testscript.sh", strconv.Itoa(id)).Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", out)
}
