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

func Ls(args []string) []byte {
	cmd := exec.Command("ls", args...)
	cmdOut, err := cmd.Output()
	if err != nil {
		log.Println(":ls command failed with error: ", err.Error())
		return []byte{}
	}
	return cmdOut
}

func Cd(dirname string) []byte {
	err := os.Chdir(dirname)
	if err != nil {
		log.Println(":cd command failed with error: ", err.Error())
		return []byte{}
	}

	// Display current working directory back to client.
	dir, _ := os.Getwd()
	return []byte(dir)
}

func Cwd() []byte {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(":cwd command failed with error: ", err.Error())
		return nil
	}
	return []byte(dir)
}

func Mkdir(dirname string) []byte {
	err := os.Mkdir(dirname, 0755)
	if err != nil {
		log.Println(":mkdir command failed with error: ", err.Error())
	}
	return Cwd()
}

func Rmdir(dirname string) []byte {
	err := os.RemoveAll(dirname)
	if err != nil {
		log.Panic(":rmdir command failed with error: ", err.Error())
	}
	return Cwd()
}

func Tree() []byte {
	panic(errors.New("Tree command is not implemented yet."))
}

func Touch(filename string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Println(":touch command failed with error: ", err.Error())
		return
	}
	defer f.Close()
}

func Cat(filepath string) []byte {
	contents, err := os.ReadFile(filepath)
	if err != nil {
		log.Println(":cat command failed with error: ", err.Error())
		return []byte{}
	}
	contents = append(contents, '\n')
	return contents
}

// Very small subset of what rm command actually can support.
func Rm(args []string) {
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
		return
	}
}

// Read the file on the server side into a byte array and send to the client.
func Get(filepath string) *os.File {
	file, err := os.Open(filepath)
	if err != nil {
		log.Println(err)
		return nil
	}

	// NOTE(alx): What if we return already closed file descriptor?
	// defer file.Close()
	return file
}
