package main

import (
    "os/exec"
    "os"
    "strings"
    "fmt"
    "log"
    "strconv"
)

// env GOOS=darwin GOARCH=amd64 go build -o dist/sweeper-darwin-amd64  sweeper.go
// env GOOS=linux GOARCH=386 go build -o dist/sweeper-linux-386  sweeper.go
// env GOOS=linux GOARCH=amd64 go build -o dist/sweeper-linux-amd64  sweeper.go
// env GOOS=windows GOARCH=amd64 go build -o dist/sweeper-windows-amd64  sweeper.go

type Info struct {
    directory   bool
    owner       string
    group       string
    size        uint64
    name        string
}

func pprint(contents []Info) {
    fmt.Printf("[\n")
    for _, c := range contents {
        fmt.Printf("\t%#v\n", c)
    }
    fmt.Printf("]\n\n")
}

func scanContents(location string) ([]Info, uint64) {
    var total_size uint64 = 0

    fmt.Printf("Scanning %s ...\n", location)

    // command to get rid of any alias to ensure format
    cmd := exec.Command("command", "ls", "-l", location)
    out, err := cmd.Output()
    if err != nil {
        fmt.Printf("%s\n", cmd.Stderr)
        log.Fatal(err)
    }

    spl := strings.Split(string(out), "\n")
    take := spl[1 : (len(spl) - 1)]
    var contents []Info

    for _, line := range take {

        fields := strings.Fields(line)
        // account for spaces before and after file/folder name
        // kinda ugly, and not that efficient, but works
        name := strings.Replace(strings.Split(line, (fields[7] + " "))[1], " ", "\\ ", -1)
        fmt.Printf("NAME: %s\n", name)

        directory := (fields[0][0] == 'd')

        var size uint64
        if directory {
            _, size = scanContents(location + "/" + name)
        } else {
            size, err = strconv.ParseUint(fields[4], 10, 64)
            if err != nil {
                log.Fatal(err)
            }
        }
        // add to total size for current directory
        total_size += size

        info := Info{   
                        directory : directory,
                        owner     : fields[2],
                        group     : fields[3],
                        size      : size,
                        name      : name,
                    }

        contents = append(contents, info)
    }

    pprint(contents)

    return contents, total_size
}

func main() {
    dir, err := os.Getwd()
    if err != nil {
        log.Fatal(err)
    }

    // _, size := scanContents(dir)
    // fmt.Printf("Size of that Directory: %d\n", size)

    // so apparently you need the space to be ascii code... rip 90 minutes
    cmd := exec.Command("command", "ls", "-l", "\x20\x20\x20\x20test\x20\x20\x20\x20folder\x20\x20\x20\x20")
    fmt.Printf("%#v \n", cmd.Args)

    out, err := cmd.Output()
    if err != nil {
        fmt.Printf("%s\n", cmd.Stderr)
        log.Fatal(err)
    }
    fmt.Printf("%s\n%s\n", dir, out)

    // what happens if you dont have access?
    // json?
}


