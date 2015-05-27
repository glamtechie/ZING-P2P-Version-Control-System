package zing

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

)

const (
	TEMP = "temp"
)

// this path need changes
var absPath string = "/Users/Vector/Workspace/Spring 2015/Distributed System/P2P-Version-Control-System/src/zing/"

func zing_init(id int) error {
	out, err := exec.Command("/bin/sh", absPath + "filesystem_scripts/zing_init.sh", strconv.Itoa(id)).Output()
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Printf("%s\n", out)
	return nil
}

func zing_pull(branch string) error {
	out, err := exec.Command("/bin/sh", absPath +  "filesystem_scripts/zing_pull.sh", branch).Output()
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Printf("%s\n", out)
	return nil
}

func zing_add(filename string) error {
	out, err := exec.Command("/bin/sh", absPath + "filesystem_scripts/zing_add.sh", filename).Output()
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Printf("%s\n", out)
	return nil
}

func zing_commit() error {
	out, err := exec.Command("/bin/sh", absPath + "filesystem_scripts/zing_commit.sh").Output()
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Printf("%s\n", out)
	return nil
}

func zing_make_patch_for_push(branch string, patchname string) (error, []byte) {
	out, err := exec.Command("/bin/sh", absPath + "filesystem_scripts/zing_make_patch_for_push.sh", branch, patchname).Output()
	b_array := make([]byte, 0)
	if err != nil {
		log.Fatal(err)
		return err, b_array
	}
	fmt.Printf("%s\n", out)
	b_array = zing_read_patch(patchname)
	return nil, b_array
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

func zing_patch_path(patchname string) string {
	return ".zing/global/" + patchname
}

func zing_abort_push() {
	//TODO: ?
}

func zing_process_push(patchname string, filecontent []byte) error {
	zing_write_patch(patchname, filecontent)

	out, err := exec.Command("/bin/sh", absPath + "filesystem_scripts/zing_process_push.sh", patchname).Output()
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Printf("%s\n", out)
	filepath := zing_patch_path(patchname)
	out, err = exec.Command("rm", filepath).Output()
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func zing_read_patch(patchname string) []byte {
	filepath := zing_patch_path(patchname)
	file, err := os.Open(filepath) // For read access.
	if err != nil {
		panic("Can't open the patch file")
	}

	result := make([]byte, 0)
	data := make([]byte, 100)
	count := 100
	for count == 100 {
		count, err = file.Read(data)
		if err != nil {
			panic("Read file error")
		} else {
			result = append(result, data...)
		}
	}

	file.Close()
	return result
}

func zing_write_patch(patchname string, filecontent []byte) {
	filepath := zing_patch_path(patchname)
	file, err := os.Create(filepath) // For read access.
	if err != nil {
		panic("Can't create the patch file")
	}

	count, e := file.Write(filecontent)
	if e != nil || count != len(filecontent) {
		panic("Write file error")
	}

	file.Close()
	return
}
