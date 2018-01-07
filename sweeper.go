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

// Holds directories that produce an error when accessing
var bad_dirs []string

// Print given number of tabs before given printf format string and argument
// Pretty hacky... not very robust - only takes 1 arg...
// but pretty printing will prolly go away in the future anyway
func tab_print(format string, tabs int, arg interface{}) {
    for i := 0; i < tabs; i++ {
        fmt.Printf("\t")
    }

    if (arg == nil) {
        fmt.Printf(format)
    } else {
        fmt.Printf(format, arg)
    }
}

// Pretty print, with tabs and newlines, the contents of an Info struct, recursing on children
// its actually not that pretty - esp for deep recurses
func pprint_children(contents []Info, tabs int) {
    tab_print("[\n", tabs, nil)
    for _, c := range contents {
        tab_print("\t{\n", tabs, nil)
        tab_print("\tdirectory : %t,\n", tabs, c.directory)
        tab_print("\towner : %s,\n", tabs, c.owner)
        tab_print("\tgroup : %s,\n", tabs, c.group)
        tab_print("\tsize : %d,\n", tabs, c.size)
        tab_print("\tname : %s,\n", tabs, c.name)

        if c.parent == nil {
            tab_print("\tparent (deref'ed) : %s,\n", tabs, "None")            
        } else {
            tab_print("\tparent (deref'ed) : %s,\n", tabs, (*c.parent).name)
        }

        if len(c.children) == 0 {
            tab_print("\tchildren : []\n", tabs, nil)
        } else {
            tab_print("\tchildren : \n", tabs, nil)
            pprint_children(c.children, tabs + 1)
        }
        tab_print("\t},\n", tabs, nil)
    }
    tab_print("]\n", tabs, nil)
}

// account for spaces before and after file/folder name
// kinda ugly, and not that efficient, but works
func format_name(location string, prev string) (name string) {
    // so apparently you need the space to be ascii code... rip 90 minutes
    // get everything after prev, in that result, replace all spaces w/ ascii code
    name = strings.Replace(strings.Split(location, (prev + " "))[1], " ", "\x20", -1)
    return name
}

func report_errors(bds []string) {
    if len(bds) > 0 {
        fmt.Printf("=== WARNING ===\nCould not access following directories:\n")
        fmt.Printf("[\n")
        for _, dir := range bds {
            fmt.Printf("\t%s,", dir)
        }
        fmt.Printf("\n]\n")
    }
}

// Helper that execs ls -l to generate an Info struct of everything inside the given directory
// for files it just saves the info
// for directories it recurses to get the sum size of all items inside that directory
func scan_dir_contents(location string) ([]Info, uint64) {
    var total_size uint64 = 0

    //fmt.Printf("Scanning %s ...\n", location)

    // command to get rid of any alias to ensure format
    cmd := exec.Command("command", "ls", "-l", location)
    out, err := cmd.Output()
    // gracefully continue - return errored info struct list
    if err != nil {
        bad_dirs = append(bad_dirs, location)

        return []Info{
            Info{
                directory : true,
                owner     : "ERR",
                group     : "ERR",
                size      : 0,
                name      : "ERROR - NO ACCESS TO PARENT",
                children  : []Info{},
                parent    : nil,
            },
        }, 0
    }

    spl := strings.Split(string(out), "\n")
    take := spl[1 : (len(spl) - 1)]
    var contents []Info

    for _, line := range take {
        fields := strings.Fields(line)

        // index 7 is right before the file name
        name := format_name(line, fields[7])

        directory := (fields[0][0] == 'd')

        var size uint64
        var children []Info
        if directory {
            children, size = scan_dir_contents(location + "/" + name)
        } else {
            size, err = strconv.ParseUint(fields[4], 10, 64)
            if err != nil {
                log.Fatal(err)
            }
            children = nil // explicit
        }
        // add to total size for current directory
        total_size += size

        info := Info{   
            directory : directory,
            owner     : fields[2],
            group     : fields[3],
            size      : size,
            name      : name,
            children  : children,
            parent    : nil,
        }

        // set parent for children - not sure if this is the best way
        // use index cause range copies the values, when we need to modify
        for i, _ := range info.children {
            info.children[i].parent = &info
        }

        contents = append(contents, info)
    }

    // pprint_children(contents)
    // fmt.Printf("Size of %s: %d\n==========================\n", location, total_size)

    // fmt.Printf("CONTENTS:%#v\n\n", contents)
    return contents, total_size
}

// Scan a directory - takes path to directory to scan
// Produces an Info struct with calculated size and directory attributes
func scan_dir(location string) Info {
    // replace spaces with ascii code - to make it work with exec
    name := strings.Replace(location, " ", "\x20", -1)
    children, size := scan_dir_contents(name)
    
    info := Info{   
        directory : true,
        owner     : "n/a",
        group     : "n/a",
        size      : size,
        name      : name,
        children  : children,
        parent    : nil,
    }

    for i, _ := range info.children {
        info.children[i].parent = &info
    }

    return info
}

func main() {
    // scans the directory the user is currently in
    // NOT the directory the executable is in!
    dir, err := os.Getwd()
    if err != nil {
        log.Fatal(err)
    }

    // pretty print list of size 1 consisting of returned Info struct
    //pprint_children([]Info{scan_dir(dir)}, 0)

    info := scan_dir(dir)
    report_errors(bad_dirs)

    start_prompt(info)


    // dont refresh if already done?
    // currently ignoring hidden files - add flags to optionally do?
    // refactor - some stuff (var names) are bad
}


