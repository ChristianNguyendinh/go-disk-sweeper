package main

import (
	"log"
	"os/exec"
	"strings"
	"fmt"
)

// env GOOS=darwin GOARCH=amd64 go build -o dist/sweeper-darwin-amd64  sweeper.go
// env GOOS=linux GOARCH=386 go build -o dist/sweeper-linux-386  sweeper.go
// env GOOS=linux GOARCH=amd64 go build -o dist/sweeper-linux-amd64  sweeper.go
// env GOOS=windows GOARCH=amd64 go build -o dist/sweeper-windows-amd64  sweeper.go

func main() {
	// command to get rid of any alias to ensure format
	// split by newline - for each, string.Fields to get fields - find way to deal with space(s)! in the name - recurse if directory
	// what happens if you dont have access?
	cmd := exec.Command("command", "ls", "-l")
	log.Printf("Running command...")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Command ran!")

	spl := strings.Split(string(out), "\n")
	take := spl[1 : (len(spl) - 1)]
	var contents [][9]string

	for _, line := range take {
		var info [9]string

		fields := strings.Fields(line)

		// account for spaces before and after file/folder name
		// kinda ugly, and not that efficient, but works
		name := strings.Split(line, (fields[7] + " "))[1]
		copy(info[0:8], fields[0:8])
		info[8] = name

		contents = append(contents, info)
		//fmt.Printf("%#v\n", info)
	}


	// fmt.Printf("%#v\n", contents)
	// print so i can see
	fmt.Printf("[\n")
	for _, c := range contents {
		fmt.Printf("\t%#v\n", c)
	}
	fmt.Printf("]\n")
}