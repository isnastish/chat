package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type FileInfo struct {
	FullPath string
	Size     int64
}

func MatchCommand(command, match string) bool {
	return bool(command == match)
}

// temporary
func serverInternalError() []byte {
	return []byte("Server internal error: 500\n\n")
}

// temporary
func appendLF(src []byte) []byte {
	src = append(src, 0x0A)
	src = append(src, 0x0A)

	return src
}

// temporary
func directoryIsNotSpecified() []byte {
	return []byte("Directory is not specified.\n\n")
}

func traverseDir(rootPath string, files chan<- FileInfo, errorStream chan<- []byte) {
	entries, err := os.ReadDir(rootPath)
	if err != nil {
		errorStream <- serverInternalError()
	}

	for _, entry := range entries {
		if entry.IsDir() {
			subDir := filepath.Join(rootPath, entry.Name())
			traverseDir(subDir, files, errorStream)
		} else {
			path, _ := os.Getwd()
			fileFullPath := filepath.Join(path, entry.Name())
			info, _ := entry.Info()
			files <- FileInfo{FullPath: fileFullPath, Size: info.Size()}
		}
	}
}

func diskUsage(args ...string) []byte {
	if len(args) != 0 {
		rootDir := args[0]
		files := make(chan FileInfo)
		errorStream := make(chan []byte, 1)

		go func() {
			traverseDir(rootDir, files, errorStream)
			close(errorStream)
			close(files)
		}()

		for err := range errorStream {
			return err
		}

		var result []string
		var totalFiles, totalBytes int64

		for file := range files {
			info := fmt.Sprintf("%-64s\t[%.2f] KB\n", file.FullPath, float32(file.Size)*(1.0/1024.0))
			result = append(result, info)

			totalFiles++
			totalBytes += file.Size
		}

		totalUsage := fmt.Sprintf("\n\nTotal files: %d, total size: [%.2f] KB\n\n", totalFiles, float32(totalBytes)*(1.0/1024.0))
		result = append(result, totalUsage)

		return []byte(strings.Join(result, ""))
	}
	return directoryIsNotSpecified()
}

func mv(args ...string) []byte {
	cmd := exec.Command("mv", args...)
	out, err := cmd.Output()
	if err != nil {
		// Failed to run the command.
		// Log on the server side.
		return serverInternalError()
	}

	out = appendLF(out)
	return out
}

func pwd(args ...string) []byte {
	cwd, err := os.Getwd()
	if err != nil {
		// TODO(alx): Log on the server side,
		// capture the command, maybe not here, but in the main processing loop.
		return serverInternalError()
	}
	return []byte(cwd + "\n\n")
}

func ls(args ...string) []byte {
	if len(args) == 0 {
		dirEntries, err := os.ReadDir(".")
		if err != nil {
			// TODO(alx): Log on the server side why the command failed.
			// and capture the command.
			return serverInternalError()
		}

		const rowEntries int = 5
		var result []string
		for index, entry := range dirEntries {
			if index != 0 && (index%rowEntries) == 0 {
				result = append(result, "\n")
			}
			name := fmt.Sprintf("%-20s", entry.Name())
			result = append(result, name)
		}
		result = append(result, "\n\n")

		return []byte(strings.Join(result, ""))
	} else {
		// This is more advance.
		// Just capture the output from executing ls in a subprocess.
		// Implement your own?
		cmd := exec.Command("ls", args...)
		out, err := cmd.Output()
		if err != nil {
			return serverInternalError()
		}
		out = appendLF(out)
		return out
	}
}

func cd(args ...string) []byte {
	if len(args) != 0 {
		err := os.Chdir(args[0])
		if err != nil {
			// Log, capture the command.
			return serverInternalError()
		}
		dir, _ := os.Getwd()
		return []byte(dir + "\n\n")
	}
	return directoryIsNotSpecified()
}

func cwd(args ...string) []byte {
	dir, err := os.Getwd()
	if err != nil {
		// Log, capture the command.
		return serverInternalError()
	}
	return []byte(dir + "\n\n")
}

func mkdir(args ...string) []byte {
	if len(args) != 0 {
		err := os.Mkdir(args[0], 0755)
		if err != nil {
			return serverInternalError()
		}
		return cwd()
	}
	return directoryIsNotSpecified()
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
	return directoryIsNotSpecified()
}

func tree(args ...string) []byte {
	panic(errors.New("Tree command is not implemented yet."))
}

func touch(args ...string) []byte {
	if len(args) != 0 {
		f, err := os.Create(args[0])
		if err != nil {
			// Log, capture the command.
			return serverInternalError()
		}
		defer f.Close()
	}
	return []byte("File is not specified.\n\n")
}

func cat(args ...string) []byte {
	if len(args) != 0 {
		contents, err := os.ReadFile(args[0])
		if err != nil {
			// Log, capture the command.
			return serverInternalError()
		}
		contents = appendLF(contents)
		return contents
	}
	return []byte("File doesn't exist.\n\n")
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
