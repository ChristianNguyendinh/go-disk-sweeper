package main

import (
    "strings"
    "os/exec"
    "strconv"
    "log"
)

// account for spaces before and after file/folder name
// kinda ugly, and not that efficient, but works
func linux_format_name(location string, prev string) (name string) {
    // so apparently you need the space to be ascii code... rip 90 minutes
    // get everything after prev, in that result, replace all spaces w/ ascii code
    name = strings.Replace(strings.Split(location, (prev + " "))[1], " ", "\x20", -1)
    return name
}

// Helper that execs ls -l to generate an Info struct of everything inside the given directory
// for files it just saves the info
// for directories it recurses to get the sum size of all items inside that directory
func linux_scan_dir_contents(location string) ([]*Info, uint64) {
    var total_size uint64 = 0

    //fmt.Printf("Scanning %s ...\n", location)

    // command to get rid of any alias to ensure format
    cmd := exec.Command("command", "ls", flags, location)
    out, err := cmd.Output()
    // gracefully continue - return errored info struct list
    if err != nil {
        bad_dirs = append(bad_dirs, location)

        return []*Info{
            &Info{
                directory : true,
                owner     : "ERR",
                group     : "ERR",
                size      : 0,
                name      : "ERROR - NO ACCESS TO PARENT",
                children  : []*Info{},
                parent    : nil,
            },
        }, 0
    }

    if len(out) == 0 {
        // return an empty slice if nothing in folder
        return []*Info{}, 0
    }

    // split and ignore first info line
    spl := strings.Split(string(out), "\n")
    take := spl[1 : (len(spl) - 1)]

    var contents []*Info

    for _, line := range take {
        fields := strings.Fields(line)

        // index 7 is right before the file name
        name := linux_format_name(line, fields[7])

        directory := (fields[0][0] == 'd')

        var size uint64
        var children []*Info
        if name == "." || name == ".." {
            continue
        } else if directory {
            children, size = linux_scan_dir_contents(location + "/" + name)
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

        contents = append(contents, &info)
    }

    // pprint_children(contents)
    // fmt.Printf("Size of %s: %d\n==========================\n", location, total_size)

    // fmt.Printf("CONTENTS:%#v\n\n", contents)
    return contents, total_size
}

// Scan a directory - takes path to directory to scan
// Produces an Info struct with calculated size and directory attributes
func linux_scan_dir(location string) Info {
    // replace spaces with ascii code - to make it work with exec
    name := strings.Replace(location, " ", "\x20", -1)
    children, size := linux_scan_dir_contents(name)
    
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