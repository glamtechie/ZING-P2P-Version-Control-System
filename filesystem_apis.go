package zing

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
)

func zing_init(id int) {
	out, err := exec.Command("/bin/sh", "filesystem_scripts/zing_init.sh", strconv.Itoa(id)).Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", out)
}

func zing_pull(branch string) {
	out, err := exec.Command("/bin/sh", "filesystem_scripts/zing_pull.sh", branch).Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", out)
}

func zing_add(filename string) {
	out, err := exec.Command("/bin/sh", "filesystem_scripts/zing_add.sh", filename).Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", out)
}

func zing_commit() {
	out, err := exec.Command("/bin/sh", "filesystem_scripts/zing_commit.sh").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", out)
}
