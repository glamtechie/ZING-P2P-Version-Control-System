package zing

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
)

const (
	TEMP = "temp"
)

func zing_init(id int) error {
	out, err := exec.Command("/bin/sh", "filesystem_scripts/zing_init.sh", strconv.Itoa(id)).Output()
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Printf("%s\n", out)
	return nil
}

func zing_pull(branch string) error {
	out, err := exec.Command("/bin/sh", "filesystem_scripts/zing_pull.sh", branch).Output()
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Printf("%s\n", out)
	return nil
}

func zing_add(filename string) error {
	out, err := exec.Command("/bin/sh", "filesystem_scripts/zing_add.sh", filename).Output()
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Printf("%s\n", out)
	return nil
}

func zing_commit() error {
	out, err := exec.Command("/bin/sh", "filesystem_scripts/zing_commit.sh").Output()
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Printf("%s\n", out)
	return nil
}

func zing_make_patch_for_push(branch string, patchname string) error {
	out, err := exec.Command("/bin/sh", "filesystem_scripts/zing_make_patch_for_push.sh", branch, patchname).Output()
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Printf("%s\n", out)
	return nil
}

func zing_delete_branch(branch string) error {
	out, err := exec.Command("git branch -D", branch).Output()
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Printf("%s\n", out)
	return nil
}

func zing_abort_push() {
	//TODO: ?
}

func zing_process_push(patchname string) error {
	out, err := exec.Command("/bin/sh", "filesystem_scripts/zing_process_push.sh", patchname).Output()
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Printf("%s\n", out)
	return nil
}
