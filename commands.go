package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func MatchCommand(command, match string) bool {
	return bool(command == match)
}

func ls(args ...string) []byte {
	cmd := exec.Command("ls", args...)
	cmdOut, err := cmd.Output()
	if err != nil {
		return []byte("ls command failed\n")
	}
	return cmdOut
}

func cd(args ...string) []byte {
	if len(args) != 0 {
		err := os.Chdir(args[0])
		if err != nil {
			log.Println(":cd command failed with error: ", err.Error())
			return []byte{}
		}
		dir, _ := os.Getwd()
		return []byte(dir)
	}
	return []byte{}
}

func cwd(args ...string) []byte {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(":cwd command failed with error: ", err.Error())
		return []byte{}
	}
	return []byte(dir)
}

func mkdir(args ...string) []byte {
	if len(args) != 0 {
		err := os.Mkdir(args[0], 0755)
		if err != nil {
			log.Println(":mkdir command failed with error: ", err.Error())
			return []byte("Internal server error\n")
		}
		return cwd()
	}
	return []byte("Directory is not specified\n")
}

func rmdir(args ...string) []byte {
	if len(args) != 0 {
		err := os.RemoveAll(args[0])
		if err != nil {
			log.Panic(":rmdir command failed with error: ", err.Error())
			return []byte("Internal server error\n")
		}
		return cwd()
	}
	return []byte("Directory is not specified\n")
}

func tree(args ...string) []byte {
	panic(errors.New("Tree command is not implemented yet."))
}

func touch(args ...string) []byte {
	if len(args) != 0 {
		f, err := os.Create(args[0])
		if err != nil {
			log.Println(":touch command failed with error: ", err.Error())
			// Return internal server error?
			return []byte{}
		}
		defer f.Close()
	}
	return []byte("file is not specified\n")
}

func cat(args ...string) []byte {
	if len(args) != 0 {
		contents, err := os.ReadFile(args[0])
		if err != nil {
			log.Println(":cat command failed with error: ", err.Error())
			return []byte{}
		}
		contents = append(contents, '\n')
		return contents
	}
	return []byte("file doesn't exist\n")
}

// Very small subset of what rm command actually can support.
func rm(args ...string) []byte {
	var handleError = func(err error) {
		if err != nil {
			log.Println(":rm command failed with error: ", err.Error())
		}
	}

	// TODO(alx): Do we really need to iterate in a for-loop?
	for k := 0; k < len(args); k++ {
		if OneOfMany(args[k], "-f", "--force") {
			if k < len(args)-1 {
				err := os.Remove(args[k+1])
				handleError(err)
			} else {
				log.Println("File is not specified.")
			}
		} else if OneOfMany(args[k], "-r", "-R", "--recursive") {
			if k < len(args)-1 {
				err := os.RemoveAll(args[k+1])
				handleError(err)
			} else {
				log.Println("File or dir is not specified.")
			}
		} else if OneOfMany(args[k], "-rf", "-Rf") {
			if k < len(args)-1 {
				err := os.Remove(args[k+1])
				handleError(err)
			} else {
				log.Println("File is not specified.")
			}
		} else {
			log.Println(":rm Invalid command arguments: ", args)
		}
	}
	return []byte{}
}
