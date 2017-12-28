package main

import (
	"log"
	"os/exec"
)

// env GOOS=darwin GOARCH=amd64 go build -o dist/sweeper-darwin-amd64  sweeper.go
// env GOOS=linux GOARCH=386 go build -o dist/sweeper-linux-386  sweeper.go
// env GOOS=linux GOARCH=amd64 go build -o dist/sweeper-linux-amd64  sweeper.go
// env GOOS=windows GOARCH=amd64 go build -o dist/sweeper-windows-amd64  sweeper.go

func main() {
	cmd := exec.Command("ls", "-lh")
	log.Printf("Running command...")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Output:\n %s\n", out)
}