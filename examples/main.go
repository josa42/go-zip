package main

import (
	"fmt"
	"os"
	"path/filepath"

	zip "github.com/josa42/go-zip"
)

func main() {
	if len(os.Args) < 2 {
		runHelp()
	}

	switch os.Args[1] {
	case "compress":
		runCompress(os.Args[2:])
	case "extract":
		runExtract(os.Args[2:])
	case "list":
		runList(os.Args[2:])
	default:
		runHelp()

	}
}

func runHelp() {
	fmt.Println("Usage:")
	fmt.Println("  go run example/main.go compress <archive> <path>")
	fmt.Println("  go run example/main.go extract <archive>")
	fmt.Println("  go run example/main.go list <archive>")
	fmt.Println("")
	os.Exit(1)
}

func runCompress(args []string) {
	if len(args) != 2 {
		runHelp()
	}

	archive := args[0]
	path := args[1]

	fmt.Println("Create:  ", archive)

	a, _ := zip.CreateArchive(archive)
	defer a.Close()

	a.Add("/", path, func(path string, sourcePath string) bool {
		if path != "." && filepath.Base(path)[0] == '.' {
			fmt.Println("> Ignore:", path)
			return false
		}

		fmt.Println("> Add:   ", path)
		return true
	})
}

func runExtract(args []string) {
	if len(args) != 1 {
		runHelp()
	}
}

func runList(args []string) {
	if len(args) != 1 {
		runHelp()
	}

	archive := args[0]

	a, _ := zip.OpenArchive(archive)
	defer a.Close()

	filePaths, _ := a.List()
	for _, f := range filePaths {
		fmt.Println(f)
	}
}
